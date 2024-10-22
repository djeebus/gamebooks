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
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"html/template"
	"net/http"
	"path/filepath"
)

//go:embed templates/*
var fs embed.FS

type views struct {
	game     bookRepo.Game
	storage  storage.Storage
	player   *executor.Player
	markdown goldmark.Markdown
}

func newViews(game bookRepo.Game, storage storage.Storage, player *executor.Player) (*views, error) {
	return &views{
		game:    game,
		storage: storage,
		player:  player,
		markdown: goldmark.New(
			goldmark.WithExtensions(
				extension.NewTable(),
				meta.Meta,
			),
		),
	}, nil
}

type indexModel struct {
	Title string
	Books []*models.Book
}

func (v *views) index(c echo.Context) error {
	books, err := v.game.GetBooks()
	if err != nil {
		return errors.Wrap(err, "failed to get books")
	}

	viewModel := indexModel{"books", books}
	return v.renderTemplate(c, "home.gohtml", viewModel)
}

const baseTemplateName = "_layout.gohtml"

func (v *views) renderTemplate(c echo.Context, templateName string, viewModel interface{}) error {
	var err error

	tmpl, err := template.ParseFS(fs,
		filepath.Join("templates", templateName),
		filepath.Join("templates", baseTemplateName),
	)
	if err != nil {
		return errors.Wrap(err, "failed to parse template")
	}

	return tmpl.ExecuteTemplate(
		c.Response().Writer,
		baseTemplateName,
		viewModel,
	)
}

func (v *views) getBook(c echo.Context) error {
	bookID := c.Param("book")

	book, err := v.game.GetBookByID(bookID)
	if err != nil {
		return errors.Wrap(err, "failed to get book")
	}

	path := fmt.Sprintf("/p/%s/%s", book.ID, book.StartPage)
	return c.Redirect(http.StatusFound, path)
}

func (v *views) getPage(c echo.Context) error {
	bookID := c.Param("book")
	pageID := c.Param("page")

	book, err := v.game.GetBookByID(bookID)
	if err != nil {
		return errors.Wrap(err, "failed to get book")
	}

	if pageID == "start" {
		path := fmt.Sprintf("/p/%s/%s", book, book.StartPage)
		return c.Redirect(http.StatusFound, path)
	}

	page, err := v.game.GetPage(bookID, pageID)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	results, err := v.player.ExecutePage(book, page, v.storage)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	viewModel, err := v.generatePageViewModel(results, book, page)
	if err != nil {
		return errors.Wrapf(err, "failed to generate viewModel bookID=%s/pageID=%s", bookID, pageID)
	}

	return v.renderTemplate(c, "page.gohtml", viewModel)
}

type pageModel struct {
	Title string
	Book  *models.Book
	Page  *models.Page

	HTML template.HTML
}

func (v *views) generatePageViewModel(results *models.PageResult, book *models.Book, page *models.Page) (pageModel, error) {
	markdown := results.Markdown
	if results.Title != "" {
		markdown = fmt.Sprintf("# %s\n%s", results.Title, markdown)
	}

	context := parser.NewContext()

	var buf bytes.Buffer
	if err := v.markdown.Convert([]byte(markdown), &buf, parser.WithContext(context)); err != nil {
		return pageModel{}, errors.Wrap(err, "failed to render markdown")
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

	result.Title = book.Name
	result.HTML = template.HTML(text)

	return result, nil
}
