package web

import (
	"bytes"
	"embed"
	"fmt"
	"gamebooks/pkg/container"
	"gamebooks/pkg/markdown"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark/parser"
	"html/template"
	"net/http"
	"strings"
)

const (
	keyPageID          = "--page-id--"
	previousPageIDsKey = "--previous-pages--"
)

//go:embed templates/*
var fs embed.FS

type views struct {
	ctr container.Container
}

func newViews(ctr container.Container) (*views, error) {
	return &views{ctr: ctr}, nil
}

type listBooksModel struct {
	Title string
	Books []*models.Book
}

const initKey = "--init-book--"

func (v *views) listBooks(c echo.Context) error {
	books, err := v.ctr.Repo.GetBooks()
	if err != nil {
		return errors.Wrap(err, "failed to get books")
	}

	viewModel := listBooksModel{"books", books}
	return v.renderTemplate(c, "books.gohtml", viewModel)
}

func getBookID(c echo.Context) string {
	return c.Param("bookID")
}

func (v *views) gameView(c echo.Context) error {
	userID := getUserID(c)
	bookID := getBookID(c)

	s := v.getBookStorage(userID, bookID)

	book, err := v.ctr.Repo.GetBookByID(bookID)
	if err != nil {
		return errors.Wrap(err, "failed to find book")
	}

	bookResults, err := v.ctr.Executor.ExecuteBook(book, s)
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

	pageID := storage.GetString(s, keyPageID)
	if pageID == "" {
		pageID := bookResults.GetStartPage()
		if err = s.Set(keyPageID, pageID); err != nil {
			return errors.Wrap(err, "failed to set current page path")
		}
		if err = storage.Push[string](s, previousPageIDsKey, pageID); err != nil {
			return errors.Wrap(err, "failed to set previous page paths")
		}
		return reloadPage(c)
	}

	if debugPageID := c.QueryParam("debug.go"); debugPageID != "" {
		pageID = debugPageID
	}

	page, err := v.ctr.Repo.GetPage(book, "", pageID)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	pageResults, err := v.ctr.Executor.ExecutePage(book, page, s)
	if err != nil {
		return errors.Wrapf(err, "failed to execute page bookID=%s/pageID=%s", bookID, pageID)
	}

	if nextPageID, err := pageResults.OnPage(); err != nil {
		return errors.Wrap(err, "failed to execute page.on_page handler")
	} else if nextPageID != "" {
		return v.navigateThroughPageID(c, s, book, pageID, nextPageID)
	}

	if pageID, err := bookResults.OnPage(page, pageResults); err != nil {
		return errors.Wrap(err, "failed to execute book.on_page handler")
	} else if pageID != "" {
		return v.navigateThroughPageID(c, s, book, pageID, pageID)
	}

	value := pageResults.Get("clear_history")
	if clearHistory, ok := value.(bool); ok && clearHistory {
		if err = s.Remove(previousPageIDsKey); err != nil {
			return errors.Wrap(err, "failed to clear history")
		}
	}

	viewModel, links, err := v.generatePageViewModel(pageResults, bookResults, book, page, pageID, s)
	if err != nil {
		return errors.Wrapf(err, "failed to generate viewModel bookID=%s/pageID=%s", bookID, pageID)
	}

	if command := c.QueryParam("cmd"); command != "" {
		for {
			nextPageID, err := pageResults.OnCommand(command)
			if err != nil {
				return errors.Wrap(err, "failed to execute command")
			}
			if nextPageID == "" {
				return reloadPage(c)
			}
			if strings.HasPrefix(nextPageID, "!") {
				command = nextPageID
				continue
			}
			return v.navigateThroughPageID(c, s, book, pageID, nextPageID)
		}
	}

	if nextPageID := c.QueryParam("goto"); nextPageID != "" {
		if nextPageID == "__previous__" {
			return v.navigateToPrevious(c, s)
		}

		for _, link := range links {
			if nextPageID == link {
				return v.navigateThroughPagePath(c, s, nextPageID)
			}
		}
		return c.String(http.StatusBadRequest, "invalid next page ID: "+nextPageID)
	}

	return v.renderTemplate(c, "page.gohtml", viewModel)
}

func (v *views) getBookStorage(userID string, bookID string) storage.Storage {
	s := v.ctr.Storage
	s = storage.NamespacedStorage(s, userID)
	s = storage.NamespacedStorage(s, bookID)
	return s
}

func reloadPage(c echo.Context) error {
	bookID := getBookID(c)
	path := fmt.Sprintf("/b/%s", bookID)
	return c.Redirect(http.StatusTemporaryRedirect, path)
}

func (v *views) navigateThroughPageID(c echo.Context, s storage.Storage, book *models.Book, currentPageID, pageID string) error {
	page, err := v.ctr.Repo.GetPage(book, currentPageID, pageID)
	if err != nil {
		return errors.Wrapf(err, "failed to find page: %q", pageID)
	}

	return v.navigateThroughPagePath(c, s, page.PageID)
}

func (v *views) navigateThroughPagePath(c echo.Context, s storage.Storage, pageID string) error {
	storage.Push[string](s, previousPageIDsKey, pageID)
	if err := s.Set(keyPageID, pageID); err != nil {
		return errors.Wrap(err, "failed to set new page path")
	}

	log.Info().Str("page_id", pageID).Msg("navigating forward")
	return reloadPage(c)
}

func (v *views) navigateToPrevious(c echo.Context, s storage.Storage) error {
	storage.Pop[string](s, previousPageIDsKey) // throw away current page id
	prevPagePath := storage.Peek[string](s, previousPageIDsKey)

	if err := s.Set(keyPageID, prevPagePath); err != nil {
		return errors.Wrap(err, "failed to set the next page path")
	}
	log.Info().Str("page_id", prevPagePath).Msg("navigating previous")
	return reloadPage(c)
}

type pageModel struct {
	Title string
	Book  models.BookResult
	Page  *models.Page
	HTML  template.HTML
}

func (v *views) buildBreadcrumbs(s storage.Storage) string {
	pageIDs := storage.GetSlice[string](s, previousPageIDsKey)
	if len(pageIDs) == 0 {
		return ""
	}

	var result string
	for _, pageID := range pageIDs {
		result += fmt.Sprintf("%s > ", pageID)
	}

	return result
}

func (v *views) generatePageViewModel(
	pageResult models.PageResult,
	bookResult models.BookResult,
	book *models.Book,
	page *models.Page,
	pagePath string,
	s storage.Storage,
) (pageModel, []string, error) {
	text, ok := pageResult.Get("markdown").(string)
	if !ok {
		return pageModel{}, nil, errors.New("failed to find markdown field")
	}

	pages := storage.GetSlice[string](s, previousPageIDsKey)
	if len(pages) > 1 {
		text = fmt.Sprintf("%s\n\n[go back](__previous__)", text)
	}

	text = v.buildBreadcrumbs(s) + "\n\n" + text

	text += "\n\n\npage id: " + page.PageID

	context := parser.NewContext()
	markdown.SetCurrentBook(context, book)
	markdown.SetCurrentPageID(context, page.PageID)

	var buf bytes.Buffer
	if err := v.ctr.Markdown.Convert([]byte(text), &buf, parser.WithContext(context)); err != nil {
		return pageModel{}, nil, errors.Wrap(err, "failed to render markdown")
	}

	text = buf.String()

	result := pageModel{
		Book: bookResult,
		Page: page,
	}

	links := markdown.GetLinksFromContext(context)

	result.Title = bookResult.GetName()
	result.HTML = template.HTML(text)

	return result, links, nil
}

func (v *views) gameClear(c echo.Context) error {
	userID := getUserID(c)
	bookID := getBookID(c)
	s := v.getBookStorage(userID, bookID)
	if err := s.Clear(""); err != nil {
		return errors.Wrap(err, "failed to clear")
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/b/%s", bookID))
}

func (v *views) setPageID(c echo.Context) error {
	userID := getUserID(c)
	bookID := getBookID(c)
	book, err := v.ctr.Repo.GetBookByID(bookID)
	if err != nil {
		return errors.Wrap(err, "failed to get book")
	}

	pageID := c.Param("pageID")
	s := v.getBookStorage(userID, bookID)

	return v.navigateThroughPageID(c, s, book, "", pageID)
}
