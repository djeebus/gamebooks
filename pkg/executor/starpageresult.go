package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

type starlarkPageResult struct {
	t      *starlark.Thread
	result map[string]any
}

func newStarlarkPageResult(t *starlark.Thread, result map[string]any) *starlarkPageResult {
	return &starlarkPageResult{
		t:      t,
		result: result,
	}
}

func (s starlarkPageResult) UpdateResults(dict map[string]any) {
	for key, val := range dict {
		s.result[key] = val
	}
}

func (s starlarkPageResult) Markdown() string {
	val, ok := s.result["markdown"].(string)
	if !ok {
		panic("missing required field: 'markdown'")
	}

	return val
}

func (s starlarkPageResult) Title() string {
	val, ok := s.result["title"].(string)
	if !ok {
		return ""
	}

	return val
}

func (s starlarkPageResult) OnCommand(command string) (string, error) {
	fn, ok := s.result["on_command"].(models.Callable)
	if !ok {
		return "", nil
	}

	result, err := fn([]any{command}, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to call on_command")
	}

	if result == nil {
		return "", nil
	}

	sval, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("expected a string result, got %T", result)
	}

	return sval, nil
}

func (s starlarkPageResult) Get(key string) any {
	return s.result[key]
}

var _ models.PageResult = &starlarkPageResult{}
