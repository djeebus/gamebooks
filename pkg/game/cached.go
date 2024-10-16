package game

import (
	"fmt"
	"gamebooks/pkg/models"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func New(cache bool) (*Game, error) {
	if !cache {
		return new(Game), nil
	}

	books := make([]*models.Book, 0)
	booksByID := make(map[string]*models.Book)
	pagesByBookIDPageID := make(map[string]map[string]*models.Page)

	bookEntries, err := os.ReadDir("books")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list books directory")
	}

	for _, bookEntry := range bookEntries {
		if !bookEntry.IsDir() {
			continue
		}

		book := models.Book{ID: bookEntry.Name()}

		err = processBookScript(&book)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to process %q", book.ID)
		}

		books = append(books, &book)
		booksByID[book.ID] = &book

		pagesByBookIDPageID[book.ID] = make(map[string]*models.Page)
		if err = filepath.Walk("books/"+book.ID+"/pages/", func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			var page *models.Page

			switch filepath.Ext(path) {
			case ".lua":
				page, err = processPageScript(path)
				if err != nil {
					return errors.Wrapf(err, "failed to process %q page script", path)
				}
			case ".md":
				page, err = buildMarkdownPage(path)
				if err != nil {
					return errors.Wrapf(err, "failed to build markdown page")
				}
			default:
				panic(fmt.Sprintf("not implemented: %q", path))
			}

			page.PageID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			page.BookID = book.ID

			pagesByBookIDPageID[book.ID][page.PageID] = page
			return nil
		}); err != nil {
			return nil, errors.Wrapf(err, "failed to list pages for %q", book.ID)
		}
	}

	return &Game{books: books, booksByID: booksByID, pagesByBookIDPageID: pagesByBookIDPageID}, nil
}

type Game struct {
	books               []*models.Book
	booksByID           map[string]*models.Book
	pagesByBookIDPageID map[string]map[string]*models.Page
}

func (g *Game) GetBooks() ([]*models.Book, error) {
	return g.books, nil
}

var ErrNotFound = errors.New("not found")

func (g *Game) GetBookByID(bookID string) (*models.Book, error) {
	if g.booksByID == nil {
		var book models.Book
		err := processBookScript(&book)
		return &book, errors.Wrapf(err, "failed to find book %q", bookID)
	}

	b, ok := g.booksByID[bookID]
	if !ok {
		return nil, ErrNotFound
	}

	return b, nil
}

func (g *Game) GetPage(bookID, pageID string) (*models.Page, error) {
	if g.pagesByBookIDPageID == nil {
		basePageName := fmt.Sprintf("books/%s/page/%s", bookID, pageID)

		var page *models.Page
		if _, err := os.Stat(basePageName + ".md"); err == nil {
			if page, err = buildMarkdownPage(basePageName + ".md"); err != nil {
				return nil, errors.Wrap(err, "failed to build markdown page")
			}
			return page, nil
		} else if _, err := os.Stat(basePageName + ".lua"); err == nil {
			if page, err = processPageScript(basePageName + ".lua"); err != nil {
				return nil, errors.Wrap(err, "failed to build lua page")
			}
			return page, nil
		}
		return nil, ErrNotFound
	}

	page, ok := g.pagesByBookIDPageID[bookID][pageID]
	if !ok {
		return nil, ErrNotFound
	}

	return page, nil
}

var ErrNoField = errors.New("missing field")

func buildMarkdownPage(path string) (*models.Page, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %q", path)
	}

	return &models.Page{
		Markdown: string(data),
	}, nil
}

func processBookScript(book *models.Book) error {
	l := lua.NewState()
	lua.OpenLibraries(l)

	path := filepath.Join("books", book.ID, "game.lua")
	if err := lua.DoFile(l, path); err != nil {
		return errors.Wrap(err, "failed to load game lua script")
	}

	var err error

	book.Name, err = getLuaStringField(l, -1, -1, "name")
	if err != nil {
		return errors.Wrap(err, "failed to load name")
	}

	book.StartPage, err = getLuaStringField(l, -2, -1, "start_page")
	if err != nil {
		return errors.Wrap(err, "failed to load start_page")
	}

	return nil
}

func processPageScript(path string) (*models.Page, error) {
	var page models.Page

	l := lua.NewState()
	lua.OpenLibraries(l)

	if err := lua.DoFile(l, path); err != nil {
		return nil, errors.Wrap(err, "failed to load page lua script")
	}

	var err error

	if page.Markdown, err = getLuaStringField(l, -1, -1, "markdown"); err != nil {
		return nil, errors.Wrap(err, "failed to load markdown")
	}

	if page.Title, err = getLuaStringField(l, -2, -1, "title"); err != nil {
		return nil, errors.Wrap(err, "failed to load title")
	}

	if page.Duration, err = getLuaIntField(l, -3, -1, "duration"); err != nil {
		return nil, errors.Wrap(err, "failed to load duration")
	}

	return &page, nil
}

func getLuaStringField(l *lua.State, idx1, idx2 int, fieldName string) (string, error) {
	l.Field(idx1, fieldName)

	val, ok := l.ToString(idx2)
	if !ok {
		return "", ErrNoField
	}

	return val, nil
}

func getLuaIntField(l *lua.State, idx1, idx2 int, fieldName string) (int, error) {
	l.Field(idx1, fieldName)

	val, ok := l.ToInteger(idx2)
	if !ok {
		return 0, ErrNoField
	}

	return val, nil
}
