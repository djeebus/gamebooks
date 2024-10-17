package bookRepo

import (
	"gamebooks/pkg/models"
)

type Game interface {
	GetBooks() ([]*models.Book, error)
	GetBookByID(bookID string) (*models.Book, error)
	GetPage(bookID, pageID string) (*models.Page, error)
}
