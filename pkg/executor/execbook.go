package executor

import (
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
)

func (p *Player) ExecuteBook(book *models.Book, storage storage.Storage) error {
	l := lua.NewState()

	luaOpenLibraries(l, book)

	open(l, "gamebooks/storage", storageFunctions(storage))
	open(l, "gamebooks/dice", diceLibrary())

	if err := lua.DoFile(l, book.LuaPath); err != nil {
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
