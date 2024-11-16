package models

type Page struct {
	Book     *Book
	PageID   string
	PagePath string
}

type PageResult interface {
	OnCommand(command string) (string, error)
	OnPage() (string, error)
	UpdateResults(map[string]any)
	Get(string) interface{}
}
