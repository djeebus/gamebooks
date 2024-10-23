package executor

import (
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
)

func (p *Player) executeBook(book *models.Book, storage storage.Storage) (*lua.State, error) {
	l := lua.NewState()

	luaOpenLibraries(l, book)

	open(l, "gamebooks/storage", storageFunctions(storage))
	open(l, "gamebooks/dice", diceLibrary())

	if err := lua.DoFile(l, book.LuaPath); err != nil {
		return nil, errors.Wrap(err, "failed to load game lua script")
	}

	return l, nil
}

func (p *Player) ExecuteBook(book *models.Book, storage storage.Storage) (*models.BookResult, error) {
	l, err := p.executeBook(book, storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute book")
	}

	var result models.BookResult

	result.Name, err = getLuaStringField(l, -1, -1, "name")
	if err != nil {
		return nil, errors.Wrap(err, "failed to load name")
	}

	result.StartPage, err = getLuaStringField(l, -2, -1, "start_page")
	if err != nil {
		return nil, errors.Wrap(err, "failed to load start_page")
	}

	return &result, nil
}

func (p *Player) BeginBook(book *models.Book, storage storage.Storage) error {
	l, err := p.executeBook(book, storage)
	if err != nil {
		return errors.Wrap(err, "failed to execute book")
	}

	l.Field(-1, "on_start")

	if !l.IsFunction(-1) {
		return errors.Wrap(err, "on_start is not a function")
	}

	l.Call(0, 0)

	return nil
}
