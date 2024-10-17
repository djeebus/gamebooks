package web

import (
	"bytes"
	"embed"
	"fmt"
	"gamebooks/pkg/models"
	"gamebooks/pkg/player"
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"html/template"
	"net/http"
)

//go:embed templates/*
var fs embed.FS

type views struct {
	game      bookRepo.Game
	storage   storage.Storage
	player    *player.Player
	templates map[string]*template.Template
	markdown  goldmark.Markdown
}

func newViews(game bookRepo.Game, storage storage.Storage) (*views, error) {
	files, err := fs.ReadDir("templates")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read templates")
	}

	templates := make(map[string]*template.Template)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := fmt.Sprintf("templates/%s", file.Name())
		tmpl, err := template.ParseFS(fs, path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse template: %s", file.Name())
		}

		templates[file.Name()] = tmpl
	}

	return &views{
		game:      game,
		storage:   storage,
		templates: templates,
		markdown:  goldmark.New(),
	}, nil
}

type indexModel struct {
	Books []*models.Book
}

func (v *views) index(c echo.Context) error {
	books, err := v.game.GetBooks()
	if err != nil {
		return errors.Wrap(err, "failed to get books")
	}

	tmpl, ok := v.templates["home.gohtml"]
	if !ok {
		return errors.New("no home template")
	}
	return tmpl.ExecuteTemplate(
		c.Response().Writer,
		"home.gohtml",
		indexModel{books},
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

type pageModel struct {
	Page *models.Page
	HTML template.HTML
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

	results, err := v.player.ExecutePage(page, v.storage)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	var buf bytes.Buffer
	if err = v.markdown.Convert([]byte(results.Markdown), &buf); err != nil {
		return errors.Wrap(err, "failed to render markdown")
	}

	tmpl, ok := v.templates["page.gohtml"]
	if !ok {
		return errors.New("no home template")
	}

	return tmpl.ExecuteTemplate(
		c.Response().Writer,
		"page.gohtml",
		pageModel{Page: page, HTML: template.HTML(buf.String())},
	)
}
