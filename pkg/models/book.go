package models

type Book struct {
	ID   string
	Path string
}

type BookResult interface {
	GetName() string
	GetStartPage() string
	OnStart() error
	OnPage(*Page, *PageResult) error
}
