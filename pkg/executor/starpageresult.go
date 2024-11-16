package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

type starlarkPageResult struct {
	t      *starlark.Thread
	page   *models.Page
	result map[string]any
}

var _ models.PageResult = new(starlarkPageResult)

func newStarlarkPageResult(t *starlark.Thread, page *models.Page, result map[string]any) *starlarkPageResult {
	return &starlarkPageResult{
		t:      t,
		page:   page,
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

func (s starlarkPageResult) OnPage() (string, error) {
	onPage, ok := s.result["on_page"].(models.Callable)
	if !ok {
		return "", nil
	}

	page := map[string]any{
		"page_id":  s.page.PageID,
		"markdown": s.result["markdown"],
	}

	result, err := onPage([]any{page}, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to call function")
	}

	if pageID, ok := result.(string); ok {
		return pageID, nil
	}

	if dict, ok := result.(map[string]any); ok {
		s.UpdateResults(dict)
	}

	return "", nil
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
