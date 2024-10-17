package models

type Page struct {
	BookID, PageID string
	Path           string
}

type PageResult struct {
	Markdown string
	Title    string
}
