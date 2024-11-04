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

	predeclared := starlarkPredeclared(storage, nil)

	opts := syntax.FileOptions{}

	stack := newStack[string](rootDir)

	t.Load = starlarkLoad(&opts, predeclared, stack)

	result, err := starlark.ExecFileOptions(&opts, &t, path, nil, predeclared)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec starlark book script")
	}

	return newStarlarkBookResult(&t, result), nil
}

func starlarkLoad(opts *syntax.FileOptions, predeclared starlark.StringDict, stack *stack[string]) func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	return func(t *starlark.Thread, module string) (starlark.StringDict, error) {
		rootDir := stack.peek()
		modulePath := filepath.Join(rootDir, module)
		newRootDir := filepath.Dir(modulePath)
		stack.push(newRootDir)
		result, err := starlark.ExecFileOptions(opts, t, modulePath, nil, predeclared)
		stack.pop()
		return result, err
	}
}

func processPageStarlarkScript(book *models.Book, page *models.Page, path string, s storage.Storage) (models.PageResult, error) {
	var t starlark.Thread

	rootDir := filepath.Dir(path)

	predeclared := starlarkPredeclared(s, page)

	opts := syntax.FileOptions{}

	stack := newStack(rootDir)

	t.Load = starlarkLoad(&opts, predeclared, stack)

	result, err := starlark.ExecFileOptions(&opts, &t, path, nil, predeclared)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec starlark page script")
	}

	unwrapped, err := unwrapStarlarkDict(&t, result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unwrap starlark value")
	}

	return newStarlarkPageResult(&t, page, unwrapped), nil
}
