package executor

import (
	"gamebooks/pkg/models"
	"github.com/pkg/errors"
	"os"
)

func processMarkdownPage(path string) (*models.PageResult, error) {
	var results models.PageResult

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	results.Markdown = string(data)
	return &results, nil
}
