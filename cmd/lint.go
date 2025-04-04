package cmd

import (
	"fmt"
	"gamebooks/pkg"
	"gamebooks/pkg/container"
	"gamebooks/pkg/game"
	"gamebooks/pkg/storage/memory"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"slices"
)

type linterView struct {
	opts game.RenderOptions
}

func (l *linterView) Reload() error {
	return nil
}

func (l *linterView) Render(opts game.RenderOptions) error {
	l.opts = opts
	return nil
}

func lintBook(ctr container.Container, bookName string) error {
	bookDir := filepath.Join("books", bookName)
	info, err := os.Stat(bookDir)
	if err != nil {
		return errors.Wrap(err, "dir does not exist")
	}

	if !info.IsDir() {
		return errors.Wrap(err, "path is not a directory")
	}

	book, err := ctr.Repo.GetBookByID(bookName)
	if err != nil {
		return errors.Wrap(err, "failed to get book")
	}

	var (
		rendered = pkg.NewSet()
		needed   = pkg.NewSet()
	)

	pages, err := ctr.Repo.GetPages(book)
	if err != nil {
		return errors.Wrap(err, "failed to get pages")
	}

	for _, page := range pages {
		s := memory.New()
		if err := s.Set("lint-mode", true); err != nil {
			return errors.Wrap(err, "failed to set lint mode")
		}

		_, err = ctr.Executor.ExecuteBook(book, s)
		if err != nil {
			return errors.Wrap(err, "failed to execute book")
		}

		var view linterView

		query := map[string][]string{"debug.go": {page.PageID}}

		g := game.New(ctr, &view)
		if err := g.Execute(book.ID, s, game.ExecuteOptions{QueryParams: query}); err != nil {
			fmt.Printf("- %s: failed to execute page: %v\n", page.PageID, err)

			continue
		}

		for _, link := range view.opts.Links {
			needed.Add(link)
		}

		for _, command := range view.opts.Commands {
			if err := g.Execute(book.ID, s, game.ExecuteOptions{QueryParams: map[string][]string{
				"cmd": {command},
			}}); err != nil {
				fmt.Printf("- %s: failed to execute %q command: %v\n", page.PageID, command, err)

				continue
			}
		}

		rendered.Add(page.PageID)
	}

	missing := needed.Without(rendered)
	if len(missing) > 0 {
		var keys []string
		for key := range missing {
			keys = append(keys, key)
		}
		slices.Sort(keys)

		fmt.Println(len(missing), "missing pages")
		for key := range missing {
			fmt.Println("- ", key)
		}
	}

	return nil
}
