package executor

import (
	"gamebooks/pkg/models"
	"github.com/Shopify/go-lua"
	"github.com/pkg/errors"
	"os"
)

func New() *Player {
	return &Player{}
}

type Player struct {
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

func getLuaStringField(l *lua.State, idx1, idx2 int, fieldName string) (string, error) {
	l.Field(idx1, fieldName)

	val, ok := l.ToString(idx2)
	if !ok {
		return "", ErrNoField
	}

	return val, nil
}
