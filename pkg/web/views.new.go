package web

import (
	"bytes"
	"fmt"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"html/template"
	"net/http"
)

const (
	keyBookID = "--book-id--"
	keyPageID = "--page-id--"
)

type indexModel struct {
	Title string
	Books []*models.Book
}

func (v *views) listBooks(c echo.Context) error {
	books, err := v.game.GetBooks()
	if err != nil {
		return errors.Wrap(err, "failed to get books")
	}

	viewModel := indexModel{"books", books}
	return v.renderTemplate(c, "books.gohtml", viewModel)
}

func (v *views) selectBook(c echo.Context) error {
	userID := getUserID(c)
	bookID := c.Param("bookID")

	s := storage.NamespacedStorage(v.storage, userID)
	s.Set(keyBookID, bookID)

	return c.Redirect(http.StatusTemporaryRedirect, "/game")
}

func (v *views) gameView(c echo.Context) error {
	userID := getUserID(c)
	playerStorage := storage.NamespacedStorage(v.storage, userID)

	bookID := storage.GetString(playerStorage, keyBookID)
	if bookID == "" {
		return c.Redirect(http.StatusFound, "/")
	}

	book, err := v.game.GetBookByID(bookID)
	if err != nil {
		return errors.Wrap(err, "failed to find book")
	}

	bookResults, err := v.player.ExecuteBook(book, playerStorage)
	if err != nil {
		return errors.Wrap(err, "failed to execute book")
	}

	pageID := storage.GetString(playerStorage, keyPageID)
	if pageID == "" {
		pageID = bookResults.StartPage
	}

	page, err := v.game.GetPage(bookID, pageID)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	pageResults, err := v.player.ExecutePage(book, page, playerStorage)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
	}

	viewModel, err := v.generatePageViewModel(pageResults, bookResults, page)
	if err != nil {
		return errors.Wrapf(err, "failed to generate viewModel bookID=%s/pageID=%s", bookID, pageID)
	}

	return v.renderTemplate(c, "page.gohtml", viewModel)
}

type pageModel struct {
	Title string
	Book  *models.BookResult
	Page  *models.Page
	HTML  template.HTML
}

type linkTracker struct{}

func newLinkTracker() *linkTracker {
	return new(linkTracker)
}

func (v *views) generatePageViewModel(results *models.PageResult, book *models.BookResult, page *models.Page) (pageModel, error) {
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

	links := getLinksFromContext(context)
	println(links)

	result.Title = book.Name
	result.HTML = template.HTML(text)

	return result, nil
}

func (v *views) gameNext(c echo.Context) error {
	panic("implement me")
}
