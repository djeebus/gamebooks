package executor

import (
	"gamebooks/pkg/models"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"regexp"
	"strings"
)

func processMarkdownPage(book *models.Book, page *models.Page, pagePath string) (models.PageResult, error) {
	data, err := os.ReadFile(pagePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	markdown := string(data)
	metadata := map[string]any{}
	metadataText, content := splitMetadataFromContent(markdown)
	if len(metadataText) > 0 {
		metadata["markdown"] = content
		if err := yaml.Unmarshal([]byte(metadataText), &metadata); err != nil {
			return nil, errors.Wrap(err, "failed to parse metadata")
		}
	} else {
		metadata["markdown"] = markdown
	}

	return newMarkdownPageResults(metadata), nil
}

var lineRegexp = regexp.MustCompile("(?m)^-+$")

func splitMetadataFromContent(markdown string) (string, string) {
	results := lineRegexp.FindAllStringIndex(markdown, -1)
	if len(results) < 2 {
		return "", markdown
	}

	res1 := results[0]
	res2 := results[1]

	metadata := markdown[res1[0]:res2[1]]
	markdown = markdown[res2[1]:]

	return metadata + "\n\n", strings.TrimSpace(markdown)
}

type markdownPageResults struct {
	metadata map[string]interface{}
}

var _ models.PageResult = new(markdownPageResults)

func (m markdownPageResults) Get(key string) any {
	return m.metadata[key]
}

func (m markdownPageResults) UpdateResults(dict map[string]any) {
	for key, value := range dict {
		m.metadata[key] = value
	}
}

func (m markdownPageResults) Markdown() string {
	val, ok := m.metadata["markdown"].(string)
	if !ok {
		return ""
	}
	return val
}

func (m markdownPageResults) OnCommand(command string) (string, error) {
	if c, ok := m.metadata["on_command"].(models.Callable); ok {
		result, err := c([]any{command}, nil)
		if err != nil {
			return "", errors.Wrap(err, "failed to call OnCommand")
		}

		return result.(string), nil
	}
	return "", nil
}

func (m markdownPageResults) OnPage() (string, error) {
	return "", nil
}

func newMarkdownPageResults(metadata map[string]interface{}) *markdownPageResults {
	return &markdownPageResults{metadata: metadata}
}
