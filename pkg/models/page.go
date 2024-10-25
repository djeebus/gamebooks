package models

type Page struct {
	BookID, PageID string
	Path           string
}

type PageResult interface {
	Markdown() string
	Title() string
	OnCommand(command string) (string, error)
	UpdateResults(map[string]any)
	Get(string) interface{}
}
