package bookRepo

import (
	"gamebooks/pkg/executor"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func NewWithLiveReload(player *executor.Player) Game {
	return &LiveReload{player: player}
}

type LiveReload struct {
	player *executor.Player
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

		book := models.Book{
			ID:      bookEntry.Name(),
			Path:    bookEntry.Name(),
			LuaPath: filepath.Join(bookEntry.Name(), "game.lua"),
		}

		if err = l.player.ExecuteBook(&book, storage.Noop()); err != nil {
			return nil, errors.Wrap(err, "failed to execute gamebook")
		}

		books = append(books, &book)
	}

	return books, nil
}

func (l *LiveReload) GetBookByID(bookID string) (*models.Book, error) {
	book := models.Book{
		ID:      bookID,
		Path:    filepath.Join("books", bookID),
		LuaPath: filepath.Join("books", bookID, "game.lua"),
	}
	err := l.player.ExecuteBook(&book, storage.Noop())
	return &book, errors.Wrapf(err, "failed to build book %q", bookID)
}

var ErrDone = errors.New("done")

func (l *LiveReload) findPageFile(bookID, pageID string) (string, error) {
	pagesPath := filepath.Join("books", bookID)

	var pagePath string

	err := filepath.Walk(pagesPath, func(path string, info fs.FileInfo, err error) error {
		if info == nil { // but why??
			return nil
		}

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
