package repo

import (
	"fmt"
	"gamebooks/pkg/executor"
	"gamebooks/pkg/models"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"path/filepath"
)

func NewWithLiveReload(executor *executor.Executor) Repo {
	return &LiveReload{executor: executor}
}

var _ Repo = new(LiveReload)

type LiveReload struct {
	executor *executor.Executor
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
			ID:   bookEntry.Name(),
			Path: bookEntry.Name(),
		}

		books = append(books, &book)
	}

	return books, nil
}

func (l *LiveReload) GetBookByID(bookID string) (*models.Book, error) {
	book := models.Book{
		ID:   bookID,
		Path: filepath.Join("books", bookID),
	}
	return &book, nil
}

func (l *LiveReload) findPageFile(book *models.Book, currentPagePath, pageID string) (string, error) {
	currentDir := ""
	if currentPagePath != "" {
		currentDir = filepath.Dir(currentPagePath)
	}

	pagesPath, err := url.JoinPath(book.Path, currentDir, pageID+".*")
	if err != nil {
		return "", errors.Wrap(err, "failed to find file")
	}

	files, err := filepath.Glob(pagesPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to find page: %s", pagesPath)
	}

	switch len(files) {
	case 0:
		return "", ErrNotFound
	case 1:
		path, err := filepath.Rel(book.Path, files[0])
		if err != nil {
			return "", errors.Wrap(err, "failed to find a relative path")
		}
		return path, nil
	default:
		return "", fmt.Errorf("file not found: %s", pagesPath)
	}
}

func (l *LiveReload) FindPagePath(book *models.Book, currentPathPath string, pageID string) (string, error) {
	pagePath, err := l.findPageFile(book, currentPathPath, pageID)
	if err != nil {
		return "", errors.Wrap(err, "failed to find page file")
	}
	return pagePath, nil
}

func (l *LiveReload) GetPage(bookID, pageID string) (*models.Page, error) {
	return &models.Page{
		BookID: bookID,
		PageID: pageID,
	}, nil
}
