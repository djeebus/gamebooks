package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func (p *Executor) ExecuteBook(book *models.Book, storage storage.Storage) (models.BookResult, error) {
	starlarkFile := filepath.Join(book.Path, "game.star")
	if _, err := os.Stat(starlarkFile); err == nil {
		results, err := processBookStarlarkScript(starlarkFile, book, storage)
		return results, errors.Wrap(err, "failed to process book starlark script")
	}

	return nil, fmt.Errorf("not implemented: %q", book.Path)
}

func (p *Executor) ExecutePage(book *models.Book, page *models.Page, storage storage.Storage) (models.PageResult, error) {
	switch filepath.Ext(page.Path) {
	case ".star":
		results, err := processPageStarlarkScript(page.Path, book, page, storage)
		return results, errors.Wrap(err, "failed to page process script")
	case ".md":
		results, err := processMarkdownPage(page.Path)
		return results, errors.Wrapf(err, "failed to build markdown page")
	default:
		return nil, fmt.Errorf("not implemented: %q", page.Path)
	}
}
