package repo

import (
	"fmt"
	"gamebooks/pkg/executor"
	"gamebooks/pkg/models"
	"github.com/pkg/errors"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

var pageExtensions = []string{".md", ".star"}

func (l *LiveReload) GetPages(book *models.Book) ([]*models.Page, error) {
	var pages []*models.Page

	if err := filepath.Walk(book.Path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk %q", path)
		}

		rel, err := filepath.Rel(book.Path, path)
		if err != nil {
			return errors.Wrap(err, "failed to calculate relative path")
		}

		if rel == "game.star" {
			return nil
		}

		if rel == "lib" {
			return fs.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if !slices.Contains(pageExtensions, ext) {
			return nil
		}
		base := strings.TrimSuffix(rel, ext)

		pages = append(pages, &models.Page{
			Book:     book,
			PageID:   base,
			PagePath: rel,
		})
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to walk book")
	}

	return pages, nil
}

func (l *LiveReload) GetPage(book *models.Book, currentPagePath, pageIDHint string) (*models.Page, error) {
	pagePath, err := l.findPageFile(book, currentPagePath, pageIDHint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find page file")
	}

	fullPageID := pageIDfromPath(pagePath)

	return &models.Page{
		Book:     book,
		PageID:   fullPageID,
		PagePath: pagePath,
	}, nil
}

func pageIDfromPath(path string) string {
	if pos := strings.LastIndexByte(path, '.'); pos != -1 {
		return path[:pos]
	}
	return path
}
