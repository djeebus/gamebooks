package cmd

import (
	"bytes"
	"fmt"
	"gamebooks/pkg"
	"gamebooks/pkg/container"
	"gamebooks/pkg/game"
	"gamebooks/pkg/markdown"
	"gamebooks/pkg/storage"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/parser"
	"os"
	"path/filepath"
	"slices"
)

type linterView struct {
	links []string
}

func (l linterView) Reload() error {
	//TODO implement me
	panic("implement me")
}

func (l linterView) Render(opts game.RenderOptions) error {
	//TODO implement me
	panic("implement me")
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

	s := storage.NewInMemory()
	_, err = ctr.Executor.ExecuteBook(book, s)
	if err != nil {
		return errors.Wrap(err, "failed to execute book")
	}

	var (
		rendered = pkg.NewSet()
		needed   = pkg.NewSet()
	)

	pages, err := ctr.Repo.GetPages(book)
	for _, page := range pages {
		fullPagePath := filepath.Join(book.Path, page.PagePath)

		var view linterView

		g := game.New(ctr, &view)
		if err := g.Execute(book.ID, s, game.ExecuteOptions{QueryParams: map[string][]string{}}); err != nil {
			return errors.Wrap(err, "failed to execute page")
		}

		for _, link := range view.links {
			needed.Add(link)
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
