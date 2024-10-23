package executor

import (
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
)

func processBookLuaScript(path string, book *models.Book, storage storage.Storage) (models.BookResult, error) {
	l := lua.NewState()

	luaOpenLibraries(l, book)

	open(l, "gamebooks/storage", storageFunctions(storage))
	open(l, "gamebooks/dice", diceLibrary())

	if err := lua.DoFile(l, path); err != nil {
		return nil, errors.Wrap(err, "failed to load game lua script")
	}

	return newBookResult(l)
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
