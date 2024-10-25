package web

import (
	"bytes"
	"embed"
	"fmt"
	"gamebooks/pkg/executor"
	"gamebooks/pkg/models"
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"html/template"
	"net/http"
	"regexp"
)

const (
	keyBookID        = "--book-id--"
	keyPageID        = "--page-id--"
	previousPagesKey = "--previous-pages--"
)

//go:embed templates/*
var fs embed.FS

type views struct {
	game     bookRepo.Game
	storage  storage.Storage
	executor *executor.Executor
	markdown goldmark.Markdown
}

func newViews(game bookRepo.Game, storage storage.Storage, executor *executor.Executor) (*views, error) {
	return &views{
		game:     game,
		storage:  storage,
		executor: executor,
		markdown: goldmark.New(
			goldmark.WithExtensions(
				meta.New(),
				extension.NewTable(),
				NewLinkTracker(),
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		),
	}, nil
}

type listBooksModel struct {
	Title string
	Books []*models.Book
}

const initKey = "--init-book--"

func (v *views) gameView(c echo.Context) error {
	userID := getUserID(c)
	playerStorage := storage.NamespacedStorage(v.storage, userID)

	if newBookID := c.QueryParam("bookID"); newBookID != "" {
		playerStorage.Set(keyBookID, newBookID)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	bookID := storage.GetString(playerStorage, keyBookID)
	if bookID == "" {
		books, err := v.game.GetBooks()
		if err != nil {
			return errors.Wrap(err, "failed to get books")
		}

		viewModel := listBooksModel{"books", books}
		return v.renderTemplate(c, "books.gohtml", viewModel)
	}

	book, err := v.game.GetBookByID(bookID)
	if err != nil {
		return errors.Wrap(err, "failed to find book")
	}

	bookResults, err := v.executor.ExecuteBook(book, playerStorage)
	if err != nil {
		return errors.Wrap(err, "failed to execute book")
	}

	if !storage.GetBool(playerStorage, initKey) {
		if err := bookResults.OnStart(); err != nil {
			return errors.Wrap(err, "failed to init story")
		}
		playerStorage.Set(initKey, true)
	}

	pageID := storage.GetString(playerStorage, keyPageID)
	if pageID == "" {
		pageID = bookResults.GetStartPage()
		playerStorage.Set(keyPageID, pageID)
		storage.Push[string](playerStorage, previousPagesKey, pageID)
	}

	page, err := v.game.GetPage(bookID, pageID)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	pageResults, err := v.executor.ExecutePage(book, page, playerStorage)
	if err != nil {
		return errors.Wrapf(err, "failed to execute page bookID=%s/pageID=%s", bookID, pageID)
	}

	if err = bookResults.OnPage(page, pageResults); err != nil {
		return errors.Wrap(err, "failed to execute on_page handler")
	}

	value := pageResults.Get("clear_history")
	if clearHistory, ok := value.(bool); ok && clearHistory {
		playerStorage.Remove(previousPagesKey)
	}

	viewModel, links, err := v.generatePageViewModel(pageResults, bookResults, page, playerStorage)
	if err != nil {
		return errors.Wrapf(err, "failed to generate viewModel bookID=%s/pageID=%s", bookID, pageID)
	}

	if command := c.QueryParam("cmd"); command != "" {
		nextPageID, err := pageResults.OnCommand(command)
		if err != nil {
			return errors.Wrap(err, "failed to execute command")
		}
		if nextPageID == "" {
			nextPageID = pageID
		}

		return v.navigateThroughLink(c, playerStorage, nextPageID)
	}

	if nextPageID := c.QueryParam("goto"); nextPageID != "" {
		if nextPageID == "__previous__" {
			return v.navigateToPrevious(c, playerStorage)
		}

		for _, link := range links {
			if nextPageID == link {
				return v.navigateThroughLink(c, playerStorage, nextPageID)
			}
		}
		return c.String(http.StatusBadRequest, "invalid next page ID: "+nextPageID)
	}

	return v.renderTemplate(c, "page.gohtml", viewModel)
}

func (v *views) navigateThroughLink(c echo.Context, s storage.Storage, pageID string) error {
	storage.Push[string](s, previousPagesKey, pageID)
	s.Set(keyPageID, pageID)
	log.Info().Str("page_id", pageID).Msg("navigating forward")
	return c.Redirect(http.StatusTemporaryRedirect, "/")

}

func (v *views) navigateToPrevious(c echo.Context, s storage.Storage) error {
	storage.Pop[string](s, previousPagesKey) // throw away current page id
	prevPageID := storage.Peek[string](s, previousPagesKey)

	s.Set(keyPageID, prevPageID)
	log.Info().Str("page_id", prevPageID).Msg("navigating previous")
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

type pageModel struct {
	Title string
	Book  models.BookResult
	Page  *models.Page
	HTML  template.HTML
}

func (v *views) buildBreadcrumbs(s storage.Storage) string {
	pageIDs := storage.GetSlice[string](s, previousPagesKey)
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
	results models.PageResult, book models.BookResult, page *models.Page, s storage.Storage,
) (pageModel, []string, error) {
	markdown := results.Markdown()

	if t := results.Title(); t != "" {
		markdown = fmt.Sprintf("# %s (%s)\n%s", t, page.PageID, markdown)
	} else {
		markdown = addPageIDtoFirstHeader(markdown, page.PageID)
	}

	pages := storage.GetSlice[string](s, previousPagesKey)
	if len(pages) > 1 {
		markdown = fmt.Sprintf("%s\n\n[go back](__previous__)", markdown)
	}

	markdown = v.buildBreadcrumbs(s) + "\n\n" + markdown

	markdown += "\n\n\npage id: " + page.PageID

	context := parser.NewContext()

	var buf bytes.Buffer
	if err := v.markdown.Convert([]byte(markdown), &buf, parser.WithContext(context)); err != nil {
		return pageModel{}, nil, errors.Wrap(err, "failed to render markdown")
	}

	text := buf.String()

	result := pageModel{
		Book: book,
		Page: page,
	}

	links := getLinksFromContext(context)

	result.Title = book.GetName()
	result.HTML = template.HTML(text)

	return result, links, nil
}

var headerRegexp = regexp.MustCompile(`(?m)^# .*$`)

func addPageIDtoFirstHeader(markdown string, pageID string) string {
	return headerRegexp.ReplaceAllStringFunc(markdown, func(s string) string {
		return fmt.Sprintf("%s (%s)", s, pageID)
	})
}
