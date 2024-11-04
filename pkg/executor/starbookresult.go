package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.starlark.net/starlark"
)

type starlarkBookResult struct {
	t      *starlark.Thread
	result starlark.StringDict
}

var _ models.BookResult = new(starlarkBookResult)

func newStarlarkBookResult(t *starlark.Thread, result starlark.StringDict) models.BookResult {
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
	fn := s.result["on_start"]
	if fn == nil {
		return nil
	}

	_, err := starlark.Call(s.t, fn, nil, nil)
	if err != nil {
		return errors.Wrap(err, "failed to call on_start")
	}
	return nil
}

func (s starlarkBookResult) OnPage(page *models.Page, result models.PageResult) (string, error) {
	onPage := s.result["on_page"]
	if onPage == nil {
		return "", nil
	}

	var err error

	input := starlark.NewDict(2)

	if value := result.Get("on_command"); value != nil {
		if err = input.SetKey(starlark.String("on_command"), starlark.Bool(true)); err != nil {
			return "", errors.Wrap(err, "failed to add on_command flag")
		}
	}

	if err = input.SetKey(starlark.String("page_id"), starlark.String(page.PageID)); err != nil {
		return "", errors.Wrap(err, "failed to set page_id")
	}

	if err = input.SetKey(starlark.String("markdown"), starlark.String(result.Markdown())); err != nil {
		return "", errors.Wrap(err, "failed to set markdown")
	}

	if _, err = starlark.Call(s.t, onPage, []starlark.Value{input}, nil); err != nil {
		return "", errors.Wrap(err, "failed to call on_page")
	}

	unwrapped, err := unwrapStarlarkValue(s.t, input)
	if err != nil {
		return "", errors.Wrap(err, "failed to unwrap input")
	}

	unwrappedMap, ok := unwrapped.(map[string]any)
	if !ok {
		return "", fmt.Errorf("could not cast %T to map[string]any", unwrapped)
	}

	if value, ok := unwrappedMap["on_command"].(bool); ok && value {
		delete(unwrappedMap, "on_command")
	}

	result.UpdateResults(unwrappedMap)

	return "", nil
}
