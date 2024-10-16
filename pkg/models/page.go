package models

type Page struct {
	BookID, PageID string
	Path           string
	Markdown       string
	Title          string
	Duration       int
}
