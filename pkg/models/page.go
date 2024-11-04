package models

type Page struct {
	BookID string
	PageID string
}

type PageResult interface {
	Markdown() string
	OnCommand(command string) (string, error)
	OnPage() (string, error)
	UpdateResults(map[string]any)
	Get(string) interface{}
}
