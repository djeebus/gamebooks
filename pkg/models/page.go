package models

import "go.starlark.net/starlark"

type Page struct {
	BookID, PageID string
	Path           string
}

type PageResult interface {
	AllowReturn() bool
	ClearHistory() bool
	Markdown() string
	Title() string
	OnCommand(command string) (string, error)
	UpdateResults(*starlark.Dict)
}
