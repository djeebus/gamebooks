package player

import (
	"fmt"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func New() *Player {
	return &Player{}
}

type Player struct {
}

func (p *Player) ExecutePage(page *models.Page, storage storage.Storage) (*models.PageResult, error) {
	switch filepath.Ext(page.Path) {
	case ".lua":
		results, err := processPageScript(page.Path, storage)
		return results, errors.Wrap(err, "failed to process script")
	case ".md":
		results, err := processMarkdownPage(page.Path)
		return results, errors.Wrapf(err, "failed to build markdown page")
	default:
		return nil, fmt.Errorf("not implemented: %q", page.Path)
	}
}

var ErrNoField = errors.New("no field")

func processMarkdownPage(path string) (*models.PageResult, error) {
	var results models.PageResult

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	results.Markdown = string(data)
	return &results, nil
}

func processPageScript(path string, storage storage.Storage) (*models.PageResult, error) {
	var results models.PageResult

	l := lua.NewState()
	lua.BaseOpen(l)
	if !lua.NewMetaTable(l, "pageMetaTable") {
		return nil, errors.New("failed to create new metatable")
	}

	lua.SetFunctions(l, []lua.RegistryFunction{
		{Name: "get_storage", Function: getStorage(storage)},
		{Name: "set_storage", Function: setStorage(storage)},
		{Name: "roll_die", Function: rollDie},
	}, 0)

	if err := lua.DoFile(l, path); err != nil {
		return nil, errors.Wrap(err, "failed to load page lua script")
	}

	var err error

	if results.Markdown, err = getLuaStringField(l, -1, -1, "markdown"); err != nil {
		return nil, errors.Wrap(err, "failed to load markdown")
	}

	if results.Title, err = getLuaStringField(l, -2, -1, "title"); err != nil {
		return nil, errors.Wrap(err, "failed to load title")
	}

	return &results, nil
}

func getStorage(storage.Storage) lua.Function {
	return func(state *lua.State) int {
		panic("getStorage")
	}
}

func setStorage(storage.Storage) lua.Function {
	return func(state *lua.State) int {
		panic("setStorage")
	}
}

func rollDie(*lua.State) int {
	panic("rollDie")
}

func getLuaStringField(l *lua.State, idx1, idx2 int, fieldName string) (string, error) {
	l.Field(idx1, fieldName)

	val, ok := l.ToString(idx2)
	if !ok {
		return "", ErrNoField
	}

	return val, nil
}
