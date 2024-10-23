package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func (p *Executor) ExecutePage(book *models.Book, page *models.Page, storage storage.Storage) (*models.PageResult, error) {
	switch filepath.Ext(page.Path) {
	case ".lua":
		results, err := processPageScript(book, page, storage)
		return results, errors.Wrap(err, "failed to page process script")
	case ".md":
		results, err := processMarkdownPage(page.Path)
		return results, errors.Wrapf(err, "failed to build markdown page")
	default:
		return nil, fmt.Errorf("not implemented: %q", page.Path)
	}
}

func processPageScript(book *models.Book, page *models.Page, storage storage.Storage) (*models.PageResult, error) {
	var err error
	var results models.PageResult

	l := lua.NewState()

	luaOpenLibraries(l, book)

	open(l, "gamebooks/storage", pageStorageFunctions(storage, page.PageID))
	open(l, "gamebooks/dice", diceLibrary())

	if err = lua.LoadFile(l, page.Path, "text"); err != nil {
		return nil, errors.Wrap(err, "failed to load file")
	}

	if err = l.ProtectedCall(0, 1, 0); err != nil {
		return nil, errors.Wrap(err, "failed to execute lua script")
	}

	var ok bool

	l.Field(-1, "markdown")
	results.Markdown, ok = l.ToString(-1)
	if !ok {
		return nil, errors.Wrap(ErrNoField, "failed to load markdown")
	}
	l.Pop(1)

	l.Field(-1, "title")
	results.Title, ok = l.ToString(-1)
	if !ok {
		return nil, errors.Wrap(ErrNoField, "failed to load title")
	}
	l.Pop(1)

	return &results, nil
}

func processMarkdownPage(path string) (*models.PageResult, error) {
	var results models.PageResult

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	results.Markdown = string(data)
	return &results, nil
}
