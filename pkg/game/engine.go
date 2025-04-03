package game

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark/parser"

	"gamebooks/pkg/container"
	"gamebooks/pkg/markdown"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
)

const (
	initKey              = "--init-book--"
	keyPageID            = "--page-id--"
	previousPageIDsKey   = "--previous-pages--"
	oncePageIDStorageKey = "--once--"
)

type Engine struct {
	ctr  container.Container
	view View
}

type ExecuteOptions struct {
	QueryParams map[string][]string
}

type RenderOptions struct {
	Commands []string
	HTML     string
	Links    []string
	Logs     []string
	Title    string
}

type View interface {
	Reload() error
	Render(opts RenderOptions) error
}

func New(ctr container.Container, view View) Engine {
	return Engine{ctr: ctr, view: view}
}

func (e Engine) Execute(bookID string, s storage.Storage, opts ExecuteOptions) error {
	book, err := e.ctr.Repo.GetBookByID(bookID)
	if err != nil {
		return errors.Wrap(err, "failed to find book")
	}

	bookResults, err := e.ctr.Executor.ExecuteBook(book, s)
	if err != nil {
		return errors.Wrap(err, "failed to execute book")
	}

	if !storage.GetBool(s, initKey) {
		if err := bookResults.OnStart(); err != nil {
			return errors.Wrap(err, "failed to init story")
		}
		if err = s.Set(initKey, true); err != nil {
			return errors.Wrap(err, "failed to set init key")
		}
	}

	var pageID string
	for _, debugPageID := range opts.QueryParams["debug.go"] {
		pageID = debugPageID
	}

	if pageID == "" {
		pageID = storage.GetString(s, keyPageID)
	}

	if pageID == "" {
		pageID := bookResults.GetStartPage()
		if err = s.Set(keyPageID, pageID); err != nil {
			return errors.Wrap(err, "failed to set current page path")
		}

		if err = storage.Push[string](s, previousPageIDsKey, pageID); err != nil {
			return errors.Wrap(err, "failed to set previous page paths")
		}
		return e.view.Reload()
	}

	if err := e.processDebugCommands(opts, s, &pageID); err != nil {
		return errors.Wrap(err, "failed to process debug commands")
	}

	page, err := e.ctr.Repo.GetPage(book, "", pageID)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	pageResults, err := e.ctr.Executor.ExecutePage(book, page, s)
	if err != nil {
		return errors.Wrapf(err, "failed to execute page bookID=%s/pageID=%s", bookID, pageID)
	}

	if nextPageID, err := pageResults.OnPage(); err != nil {
		return errors.Wrap(err, "failed to execute page.on_page handler")
	} else if nextPageID != "" {
		return e.navigateThroughPageID(s, book, page.PagePath, nextPageID)
	}

	onceFlag, err := s.Get(oncePageIDStorageKey)
	if err != nil {
		return errors.Wrap(err, "failed to get once flag")
	}
	if onceFlag != pageID {
		if err = pageResults.Once(); err != nil {
			return errors.Wrap(err, "failed to execute page.once handler")
		}
		if err = s.Set(oncePageIDStorageKey, pageID); err != nil {
			return errors.Wrap(err, "failed to set once flag")
		}
	}

	if pageID, err := bookResults.OnPage(page, pageResults); err != nil {
		return errors.Wrap(err, "failed to execute book.on_page handler")
	} else if pageID != "" {
		return e.navigateThroughPageID(s, book, page.PagePath, pageID)
	}

	value := pageResults.Get("clear_history")
	if clearHistory, ok := value.(bool); ok && clearHistory {
		if err = s.Remove(previousPageIDsKey); err != nil {
			return errors.Wrap(err, "failed to clear history")
		}
	}

	for _, command := range opts.QueryParams["cmd"] {
		command, args := splitCommand(command)
		for {
			nextPageID, err := pageResults.OnCommand(command, args)
			if err != nil {
				return errors.Wrap(err, "failed to execute command")
			}
			if nextPageID == "" {
				return e.view.Reload()
			}

			if strings.HasPrefix(nextPageID, "!") {
				command = nextPageID
				continue
			}

			return e.navigateThroughPageID(s, book, page.PagePath, nextPageID)
		}
	}

	html, links, commands, err := e.buildHTML(s, book, page, pageResults)
	if err != nil {
		return errors.Wrap(err, "failed to build html")
	}

	for _, nextPageID := range opts.QueryParams["goto"] {
		if nextPageID == "__previous__" {
			return e.navigateToPrevious(s)
		}

		for _, link := range links {
			if nextPageID == link {
				return e.navigateThroughPagePath(s, nextPageID)
			}
		}
		return errors.Wrap(ErrCannotGoto, "invalid next page ID: "+nextPageID)
	}

	logs := storage.GetLog(s)

	return e.view.Render(RenderOptions{
		Commands: commands,
		HTML:     html,
		Links:    links,
		Logs:     logs,
		Title:    bookResults.GetName(),
	})
}

var ErrCannotGoto = errors.New("cannot goto")

func (e Engine) buildBreadcrumbs(pageIDs []string) string {
	var result string
	for _, pageID := range pageIDs {
		result += fmt.Sprintf("%s > ", pageID)
	}

	return result
}

func (e Engine) buildHTML(s storage.Storage, book *models.Book, page *models.Page, pageResult models.PageResult) (string, []string, []string, error) {
	text, ok := pageResult.Get("markdown").(string)
	if !ok {
		return "", nil, nil, errors.New("failed to find markdown field")
	}

	pageIDs := storage.GetSlice[string](s, previousPageIDsKey)
	if len(pageIDs) > 1 {
		text = fmt.Sprintf("%s\n\n[go back](__previous__)", text)
	}

	text = e.buildBreadcrumbs(pageIDs) + "\n\n" + text

	text += "\n\n\npage id: " + page.PageID

	context := parser.NewContext()
	markdown.SetCurrentBook(context, book)
	markdown.SetCurrentPageID(context, page.PageID)

	var buf bytes.Buffer
	if err := e.ctr.Markdown.Convert([]byte(text), &buf, parser.WithContext(context)); err != nil {
		return "", nil, nil, errors.Wrap(err, "failed to render markdown")
	}

	text = buf.String()

	links := markdown.GetLinksFromContext(context)
	commands := markdown.GetCommandsFromContext(context)

	return text, links, commands, nil
}

func (e Engine) navigateThroughPageID(s storage.Storage, book *models.Book, currentPagePath, pageID string) error {
	page, err := e.ctr.Repo.GetPage(book, currentPagePath, pageID)
	if err != nil {
		return errors.Wrapf(err, "failed to find page: %q", pageID)
	}

	return e.navigateThroughPagePath(s, page.PageID)
}

func (e Engine) navigateThroughPagePath(s storage.Storage, pageID string) error {
	if err := storage.Push[string](s, previousPageIDsKey, pageID); err != nil {
		return errors.Wrap(err, "failed to push previous page path")
	}

	if err := s.Set(keyPageID, pageID); err != nil {
		return errors.Wrap(err, "failed to set new page path")
	}

	log.Info().Str("page_id", pageID).Msg("navigating forward")
	return e.view.Reload()
}

func (e Engine) processDebugCommands(opts ExecuteOptions, s storage.Storage, pageID *string) error {
	var err error

	for key, values := range opts.QueryParams {
		if !strings.HasPrefix(key, "debug.") {
			continue
		}

		key = strings.TrimPrefix(key, "debug.")
		if key == "go" {
			for _, value := range values {
				*pageID = value
				if err = s.Set(keyPageID, pageID); err != nil {
					return errors.Wrap(err, "failed to set debug page id")
				}
			}
			continue
		}

		if strings.HasPrefix(key, "set:") {
			key = strings.TrimPrefix(key, "set:")
			for _, value := range values {
				var data any = value

				if strings.HasPrefix(value, "int!") {
					value = strings.TrimPrefix(value, "int!")
					data, err = strconv.ParseInt(value, 10, 64)
					if err != nil {
						return errors.Wrap(err, "failed to parse int")
					}
				}

				if err = s.Set(key, data); err != nil {
					return errors.Wrapf(err, "failed to set %s = %s", key, value)
				}
			}
		}
	}

	return nil
}

func splitCommand(command string) (string, []string) {
	command = strings.TrimPrefix(command, "!")
	parts := strings.Split(command, "!")
	switch len(parts) {
	case 0:
		return "", nil
	case 1:
		return parts[0], nil
	default:
		return parts[0], parts[1:]
	}
}

func (e Engine) navigateToPrevious(s storage.Storage) error {
	if _, err := storage.Pop[string](s, previousPageIDsKey); err != nil { // throw away current page id
		return errors.Wrap(err, "failed to pop previous page id")
	}
	prevPagePath := storage.Peek[string](s, previousPageIDsKey)

	if err := s.Set(keyPageID, prevPagePath); err != nil {
		return errors.Wrap(err, "failed to set the next page path")
	}

	log.Info().Str("page_id", prevPagePath).Msg("navigating previous")
	return e.view.Reload()
}
