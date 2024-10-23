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

func (p *Executor) ExecuteBook(book *models.Book, storage storage.Storage) (models.BookResult, error) {
	luaFile := filepath.Join(book.Path, "game.lua")
	if _, err := os.Stat(luaFile); err == nil {
		results, err := processBookScript(luaFile, book, storage)
		return results, errors.Wrap(err, "failed to process book script")

	}

	return nil, fmt.Errorf("not implemented: %q", book.Path)
}

func processBookScript(path string, book *models.Book, storage storage.Storage) (models.BookResult, error) {
	l := lua.NewState()

	luaOpenLibraries(l, book)

	open(l, "gamebooks/storage", storageFunctions(storage))
	open(l, "gamebooks/dice", diceLibrary())

	if err := lua.DoFile(l, path); err != nil {
		return nil, errors.Wrap(err, "failed to load game lua script")
	}

	return newBookResult(l)
}
