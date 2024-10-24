package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"path/filepath"
)

type starlarkBookResult struct {
	t      *starlark.Thread
	result starlark.StringDict
}

func newStarlarkBookResult(t *starlark.Thread, result starlark.StringDict) *starlarkBookResult {
	return &starlarkBookResult{t, result}
}

func (s starlarkBookResult) GetName() string {
	return s.asString("name")
}

func (s starlarkBookResult) GetStartPage() string {
	return s.asString("start_page")
}

func (s starlarkBookResult) asString(key string) string {
	return asString(s.result, key)
}

func asString(result starlark.StringDict, key string) string {
	val, ok := result[key]
	if !ok {
		log.Error().Str("key", key).Msg("key not found")
		return ""
	}

	sval, ok := val.(starlark.String)
	if !ok {
		log.Warn().Type("value", val).Msg("unexpected data type")
		return ""
	}

	return sval.GoString()
}

func (s starlarkBookResult) OnStart() error {
	_, err := starlark.Call(s.t, s.result["on_start"], nil, nil)
	if err != nil {
		return errors.Wrap(err, "failed to call on_start")
	}
	return nil
}

func (s starlarkBookResult) OnPage(page *models.Page, result *models.PageResult) error {
	var err error

	input := starlark.NewDict(2)

	if err = input.SetKey(starlark.String("page_id"), starlark.String(page.PageID)); err != nil {
		return errors.Wrap(err, "failed to set page_id")
	}

	if err = input.SetKey(starlark.String("title"), starlark.String(result.Title)); err != nil {
		return errors.Wrap(err, "failed to set title")
	}

	if err = input.SetKey(starlark.String("markdown"), starlark.String(result.Markdown)); err != nil {
		return errors.Wrap(err, "failed to set markdown")
	}

	if _, err = starlark.Call(s.t, s.result["on_page"], []starlark.Value{input}, nil); err != nil {
		return errors.Wrap(err, "failed to call on_page")
	}

	if val, ok, err := getFromDict[starlark.String](input, "title"); err != nil {
		return errors.Wrap(err, "failed to get title")
	} else if !ok {
		return errors.Wrap(err, "missing required key 'title'")
	} else {
		result.Title = string(val)
	}

	if val, ok, err := getFromDict[starlark.String](input, "markdown"); err != nil {
		return errors.Wrap(err, "failed to get markdown")
	} else if !ok {
		return errors.Wrap(err, "missing required key 'markdown'")
	} else {
		result.Markdown = string(val)
	}

	if val, ok, err := getFromDict[starlark.Bool](input, "allow_return"); err != nil {
		return errors.Wrap(err, "failed to get allow_return")
	} else if ok {
		result.AllowReturn = bool(val)
	}

	return nil
}

func getFromDict[T starlark.Value](input *starlark.Dict, key string) (T, bool, error) {
	var t T

	value, ok, err := input.Get(starlark.String(key))
	if err != nil {
		return t, false, errors.Wrap(err, "failed to get key")
	}
	if !ok {
		return t, false, nil
	}
	t, ok = value.(T)
	if !ok {
		return t, false, fmt.Errorf("expected starlark.String, got %T", value)
	}

	return t, true, nil
}

func setStringFromDictKey(input *starlark.Dict, key string, destination *string) error {
	value, ok, err := input.Get(starlark.String(key))
	if err != nil {
		return errors.Wrap(err, "failed to get key")
	}
	if !ok {
		return nil
	}

	s, ok := value.(starlark.String)
	if !ok {
		return fmt.Errorf("expected starlark.String, got %T", value)
	}

	*destination = string(s)
	return nil
}

func processBookStarlarkScript(path string, book *models.Book, storage storage.Storage) (models.BookResult, error) {
	var t starlark.Thread

	rootDir := filepath.Dir(path)

	predeclared := starlarkPredeclared(storage)

	opts := syntax.FileOptions{}

	t.Load = starlarkLoad(rootDir, &opts, predeclared)

	result, err := starlark.ExecFileOptions(&opts, &t, path, nil, predeclared)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec starlark book script")
	}

	return newStarlarkBookResult(&t, result), nil
}

func starlarkLoad(rootDir string, opts *syntax.FileOptions, predeclared starlark.StringDict) func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	return func(t *starlark.Thread, module string) (starlark.StringDict, error) {
		modulePath := filepath.Join(rootDir, module)
		return starlark.ExecFileOptions(opts, t, modulePath, nil, predeclared)
	}
}

func processPageStarlarkScript(path string, book *models.Book, page *models.Page, s storage.Storage) (*models.PageResult, error) {
	var t starlark.Thread

	rootDir := filepath.Dir(path)

	predeclared := starlarkPredeclared(s)

	opts := syntax.FileOptions{}

	t.Load = starlarkLoad(rootDir, &opts, predeclared)

	result, err := starlark.ExecFileOptions(&opts, &t, path, nil, predeclared)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec starlark page script")
	}

	return &models.PageResult{
		Markdown: asString(result, "markdown"),
		Title:    asString(result, "title"),
	}, nil
}
