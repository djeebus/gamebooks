package web

import (
	"bytes"
	"embed"
	"fmt"
	"gamebooks/pkg/executor"
	"gamebooks/pkg/markdown"
	"gamebooks/pkg/models"
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"html/template"
	"net/http"
	"strings"
)

const (
	keyPagePath          = "--page-id--"
	previousPagePathsKey = "--previous-pages--"
)

//go:embed templates/*
var fs embed.FS

type views struct {
	repo     bookRepo.Repo
	storage  storage.Storage
	executor *executor.Executor
	markdown goldmark.Markdown
}

func newViews(
	repo bookRepo.Repo,
	storage storage.Storage,
	executor *executor.Executor,
	markdown goldmark.Markdown,
) (*views, error) {
	return &views{
		repo:     repo,
		storage:  storage,
		executor: executor,
		markdown: markdown,
	}, nil
}

type listBooksModel struct {
	Title string
	Books []*models.Book
}

const initKey = "--init-book--"

func (v *views) listBooks(c echo.Context) error {
	books, err := v.repo.GetBooks()
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

	book, err := v.repo.GetBookByID(bookID)
	if err != nil {
		return errors.Wrap(err, "failed to find book")
	}

	bookResults, err := v.executor.ExecuteBook(book, s)
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

	pagePath := storage.GetString(s, keyPagePath)
	if pagePath == "" {
		startPage := bookResults.GetStartPage()
		pagePath, err := v.repo.FindPagePath(book, "", startPage)
		if err != nil {
			return errors.Wrapf(err, "failed to get page path: %q", startPage)
		}
		if err = s.Set(keyPagePath, pagePath); err != nil {
			return errors.Wrap(err, "failed to set current page path")
		}
		storage.Push[string](s, previousPagePathsKey, pagePath)
		return reloadPage(c)
	}

	page, err := v.repo.GetPage(bookID, pagePath)
	if err != nil {
		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pagePath)
	}

	pageResults, err := v.executor.ExecutePage(book, page, pagePath, s)
	if err != nil {
		return errors.Wrapf(err, "failed to execute page bookID=%s/pageID=%s", bookID, pagePath)
	}

	if pageID, err := pageResults.OnPage(); err != nil {
		return errors.Wrap(err, "failed to execute page.on_page handler")
	} else if pageID != "" {
		return v.navigateThroughPageID(c, s, book, pagePath, pageID)
	}

	if pageID, err := bookResults.OnPage(page, pageResults); err != nil {
		return errors.Wrap(err, "failed to execute book.on_page handler")
	} else if pageID != "" {
		return v.navigateThroughPageID(c, s, book, pagePath, pageID)
	}

	value := pageResults.Get("clear_history")
	if clearHistory, ok := value.(bool); ok && clearHistory {
		if err = s.Remove(previousPagePathsKey); err != nil {
			return errors.Wrap(err, "failed to clear history")
		}
	}

	viewModel, links, err := v.generatePageViewModel(pageResults, bookResults, book, page, pagePath, s)
	if err != nil {
		return errors.Wrapf(err, "failed to generate viewModel bookID=%s/pageID=%s", bookID, pagePath)
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
			return v.navigateThroughPageID(c, s, book, pagePath, nextPageID)
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
	s := v.storage
	s = storage.NamespacedStorage(v.storage, userID)
	s = storage.NamespacedStorage(v.storage, bookID)
	return s
}

func reloadPage(c echo.Context) error {
	bookID := getBookID(c)
	path := fmt.Sprintf("/b/%s", bookID)
	return c.Redirect(http.StatusTemporaryRedirect, path)
}

func (v *views) navigateThroughPageID(c echo.Context, s storage.Storage, book *models.Book, currentPagePath, pageID string) error {
	filePath, err := v.repo.FindPagePath(book, currentPagePath, pageID)
	if err != nil {
		return errors.Wrapf(err, "failed to find page: %q", pageID)
	}

	return v.navigateThroughPagePath(c, s, filePath)
}

func (v *views) navigateThroughPagePath(c echo.Context, s storage.Storage, pagePath string) error {
	storage.Push[string](s, previousPagePathsKey, pagePath)
	if err := s.Set(keyPagePath, pagePath); err != nil {
		return errors.Wrap(err, "failed to set new page path")
	}

	log.Info().Str("page_id", pagePath).Msg("navigating forward")
	return reloadPage(c)
}

func (v *views) navigateToPrevious(c echo.Context, s storage.Storage) error {
	storage.Pop[string](s, previousPagePathsKey) // throw away current page id
	prevPagePath := storage.Peek[string](s, previousPagePathsKey)

	if err := s.Set(keyPagePath, prevPagePath); err != nil {
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
	pageIDs := storage.GetSlice[string](s, previousPagePathsKey)
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
	text := pageResult.Markdown()

	pages := storage.GetSlice[string](s, previousPagePathsKey)
	if len(pages) > 1 {
		text = fmt.Sprintf("%s\n\n[go back](__previous__)", text)
	}

	text = v.buildBreadcrumbs(s) + "\n\n" + text

	text += "\n\n\npage id: " + page.PageID

	context := parser.NewContext()
	markdown.SetCurrentBook(context, book)
	markdown.SetCurrentPagePath(context, pagePath)

	var buf bytes.Buffer
	if err := v.markdown.Convert([]byte(text), &buf, parser.WithContext(context)); err != nil {
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
	book, err := v.repo.GetBookByID(bookID)
	if err != nil {
		return errors.Wrap(err, "failed to get book")
	}

	pageID := c.Param("pageID")
	s := v.getBookStorage(userID, bookID)

	return v.navigateThroughPageID(c, s, book, "", pageID)
}
