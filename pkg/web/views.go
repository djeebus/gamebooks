package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"slices"

	"gamebooks/pkg/container"
	"gamebooks/pkg/game"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type echoView struct {
	c echo.Context

	renderTemplate func(c echo.Context, template string, model any) error
}

func (e echoView) Render(opts game.RenderOptions) error {
	slices.Reverse(opts.Logs)

	result := pageModel{
		HTML:  template.HTML(opts.HTML),
		Log:   opts.Logs,
		Title: opts.Title,
	}

	return e.renderTemplate(e.c, "page.gohtml", result)
}

func (e echoView) Reload() error {
	return reloadPage(e.c)
}

var _ game.View = echoView{}

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

	g := game.New(v.ctr, echoView{c: c, renderTemplate: v.renderTemplate})

	return g.Execute(bookID, s, game.ExecuteOptions{
		QueryParams: c.QueryParams(),
	})
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

type pageModel struct {
	Title string
	Book  models.BookResult
	Page  *models.Page
	HTML  template.HTML
	Log   []string
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
