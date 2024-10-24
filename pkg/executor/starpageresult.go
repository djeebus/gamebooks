package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

type starlarkPageResult struct {
	result starlark.StringDict
}

func (s starlarkPageResult) UpdateResults(dict *starlark.Dict) {
	for _, key := range dict.Keys() {
		val, ok, err := dict.Get(key)
		if err != nil && ok {
			keystr, ok := val.(starlark.String)
			if ok {
				s.result[string(keystr)] = val
			}
		}
	}
}

func newStarlarkPageResult(result starlark.StringDict) *starlarkPageResult {
	return &starlarkPageResult{
		result: result,
	}
}

func (s starlarkPageResult) AllowReturn() bool {
	val, ok, err := getFromDict[starlark.Bool](s.result, "allow_return")
	if err == nil && ok {
		return val
	}
	return false
}

func (s starlarkPageResult) ClearHistory() bool {
	//TODO implement me
	panic("implement me")
}

func (s starlarkPageResult) Markdown() string {
	return asString(s.result, "markdown")
}

func (s starlarkPageResult) Title() string {
	return asString(s.result, "title")
}

func (s starlarkPageResult) OnCommand(command string) (string, error) {
	//TODO implement me
	panic("implement me")
}

var _ models.PageResult = &starlarkPageResult{}

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
