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
	var err error
	var results models.PageResult

	l := lua.NewState()

	if err = lua.LoadFile(l, path, ""); err != nil {
		return nil, errors.Wrap(err, "failed to load file")
	}

	lua.OpenLibraries(l)
	open(l, "gamebooks/storage", []lua.RegistryFunction{
		{Name: "get", Function: getStorage(storage)},
		{Name: "set", Function: setStorage(storage)},
	})
	open(l, "gamebooks/dice", []lua.RegistryFunction{
		{Name: "roll", Function: rollDie},
	})

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

func getLuaStringField(l *lua.State, idx1, idx2 int, fieldName string) (string, error) {
	l.Field(idx1, fieldName)

	val, ok := l.ToString(idx2)
	if !ok {
		return "", ErrNoField
	}

	return val, nil
}
