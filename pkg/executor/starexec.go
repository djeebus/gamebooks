package executor

import (
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/pkg/errors"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"path/filepath"
)

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

func processPageStarlarkScript(path string, book *models.Book, page *models.Page, s storage.Storage) (models.PageResult, error) {
	var t starlark.Thread

	rootDir := filepath.Dir(path)

	predeclared := starlarkPredeclared(s)

	opts := syntax.FileOptions{}

	t.Load = starlarkLoad(rootDir, &opts, predeclared)

	result, err := starlark.ExecFileOptions(&opts, &t, path, nil, predeclared)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec starlark page script")
	}

	unwrapped, err := unwrapStarlarkDict(&t, result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unwrap starlark value")
	}

	return newStarlarkPageResult(&t, unwrapped), nil
}
