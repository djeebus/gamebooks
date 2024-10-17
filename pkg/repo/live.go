package bookRepo

import (
	"fmt"
	"gamebooks/pkg/models"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func NewWithLiveReload() Game {
	return new(LiveReload)
}

type LiveReload struct {
}

func (l *LiveReload) GetBooks() ([]*models.Book, error) {
	bookEntries, err := os.ReadDir("books")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list books directory")
	}

	books := make([]*models.Book, 0)
	for _, bookEntry := range bookEntries {
		if !bookEntry.IsDir() {
			continue
		}

		book := models.Book{ID: bookEntry.Name()}

		err = processBookScript(&book)

		books = append(books, &book)
	}

	return books, nil
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

var ErrNoField = errors.New("no field")

func getLuaStringField(l *lua.State, idx1, idx2 int, fieldName string) (string, error) {
	l.Field(idx1, fieldName)

	val, ok := l.ToString(idx2)
	if !ok {
		return "", ErrNoField
	}

	return val, nil
}

func (l *LiveReload) GetBookByID(bookID string) (*models.Book, error) {
	book := models.Book{ID: bookID}
	err := processBookScript(&book)
	return &book, errors.Wrapf(err, "failed to find book %q", bookID)
}

var ErrDone = errors.New("done")

func (l *LiveReload) findPageFile(bookID, pageID string) (string, error) {
	pagesPath := fmt.Sprintf("books/%s/pages", bookID)

	var pagePath string

	err := filepath.Walk(pagesPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		filename := filepath.Base(path)
		if !strings.HasPrefix(filename, pageID+".") {
			return nil
		}

		pagePath = path
		return ErrDone
	})

	if err == nil {
		return "", ErrNotFound
	}

	if errors.Is(err, ErrDone) {
		return pagePath, nil
	}

	return "", errors.Wrap(err, "failed to find page path")
}

func (l *LiveReload) GetPage(bookID, pageID string) (*models.Page, error) {
	pagePath, err := l.findPageFile(bookID, pageID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find page file")
	}

	return &models.Page{
		BookID: bookID,
		PageID: pageID,
		Path:   pagePath,
	}, nil
}
