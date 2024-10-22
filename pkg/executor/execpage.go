package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
	"path/filepath"
)

func (p *Player) ExecutePage(book *models.Book, page *models.Page, storage storage.Storage) (*models.PageResult, error) {
	switch filepath.Ext(page.Path) {
	case ".lua":
		results, err := processPageScript(book, page, storage)
		return results, errors.Wrap(err, "failed to process script")
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

	if results.Markdown, err = getLuaStringField(l, -1, -1, "markdown"); err != nil {
		return nil, errors.Wrap(err, "failed to load markdown")
	}

	if results.Title, err = getLuaStringField(l, -2, -1, "title"); err != nil {
		return nil, errors.Wrap(err, "failed to load title")
	}

	return &results, nil
}
