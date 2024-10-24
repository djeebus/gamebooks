package models

type Page struct {
	BookID, PageID string
	Path           string
}

type PageResult struct {
	AllowReturn  bool
	ClearHistory bool
	Markdown     string
	Title        string
}
