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
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"html/template"
	"net/http"
)

const (
	keyBookID = "--book-id--"
	keyPageID = "--page-id--"
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
				extension.NewTable(),
				meta.Meta,
				NewLinkTracker(),
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
	}

	page, err := v.game.GetPage(bookID, pageID)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	pageResults, err := v.executor.ExecutePage(book, page, playerStorage)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	if err = bookResults.OnPage(pageResults); err != nil {
		return errors.Wrap(err, "failed to execute on_page handler")
	}

	viewModel, links, err := v.generatePageViewModel(pageResults, bookResults, page)
	if err != nil {
		return errors.Wrapf(err, "failed to generate viewModel bookID=%s/pageID=%s", bookID, pageID)
	}

	if nextPageID := c.QueryParam("goto"); nextPageID != "" {
		for _, link := range links {
			if nextPageID == link {
				playerStorage.Set(keyPageID, nextPageID)
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}
		}
		return c.String(http.StatusBadRequest, "invalid next page ID: "+nextPageID)
	}

	return v.renderTemplate(c, "page.gohtml", viewModel)
}

type pageModel struct {
	Title string
	Book  models.BookResult
	Page  *models.Page
	HTML  template.HTML
}

func (v *views) generatePageViewModel(results *models.PageResult, book models.BookResult, page *models.Page) (pageModel, []string, error) {
	markdown := results.Markdown
	if results.Title != "" {
		markdown = fmt.Sprintf("# %s\n%s", results.Title, markdown)
	}

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

	metadata := meta.Get(context)
	if titleIface, ok := metadata["title"]; ok {
		title, ok := titleIface.(string)
		if ok {
			text = "<h1>" + title + "</h1>\n" + text
			result.Title = title
		}
	}

	links := getLinksFromContext(context)

	result.Title = book.GetName()
	result.HTML = template.HTML(text)

	return result, links, nil
}

func (v *views) gameNext(c echo.Context) error {
	panic("implement me")
}
