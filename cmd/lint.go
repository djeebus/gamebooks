package cmd

import (
	"bytes"
	"fmt"
	"gamebooks/pkg"
	"gamebooks/pkg/container"
	"gamebooks/pkg/markdown"
	"gamebooks/pkg/storage"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/parser"
	"os"
	"path/filepath"
	"slices"
)

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
		results, err := ctr.Executor.ExecutePage(book, page, s)
		if err != nil {
			fmt.Printf("%s: %v\n", fullPagePath, err)
			continue
		}

		context := parser.NewContext()
		markdown.SetCurrentBook(context, book)
		markdown.SetCurrentPageID(context, page.PageID)

		result, err := results.OnPage()
		if err != nil {
			fmt.Printf("%s: %v\n", fullPagePath, err)
		}
		if result != "" {
			needed.Add(result)
			continue
		}

		text, ok := results.Get("markdown").(string)
		if !ok {
			fmt.Printf("%s: failed to find markdown\n", fullPagePath)
			continue
		}

		var buf bytes.Buffer
		if err := ctr.Markdown.Convert([]byte(text), &buf, parser.WithContext(context)); err != nil {
			fmt.Printf("%s: %v\n", fullPagePath, err)
			continue
		}

		links := markdown.GetLinksFromContext(context)
		for _, link := range links {
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
