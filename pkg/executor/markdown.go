package executor

import (
	"gamebooks/pkg/models"
	"github.com/pkg/errors"
	"go.starlark.net/starlark"
	"os"
)

func processMarkdownPage(path string) (models.PageResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	return newMarkdownPageResults(string(data)), nil
}

type markdownPageResults struct {
	markdown string
}

func (m markdownPageResults) UpdateResults(dict *starlark.Dict) {
	val, ok, err := dict.Get(starlark.String("markdown"))
	if err == nil && ok {
		s, ok := val.(starlark.String)
		if ok {
			m.markdown = string(s)
		}
	}
}

func (m markdownPageResults) AllowReturn() bool {
	return true
}

func (m markdownPageResults) ClearHistory() bool {
	return false
}

func (m markdownPageResults) Markdown() string {
	return m.markdown
}

func (m markdownPageResults) Title() string {
	return ""
}

func (m markdownPageResults) OnCommand(command string) (string, error) {
	return "", nil
}

func newMarkdownPageResults(markdown string) *markdownPageResults {
	return &markdownPageResults{markdown: markdown}
}
