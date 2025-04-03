package repo

import (
	"gamebooks/pkg/models"
)

type Repo interface {
	GetBooks() ([]*models.Book, error)
	GetBookByID(bookID string) (*models.Book, error)
	GetPage(book *models.Book, currentPagePath, pageID string) (*models.Page, error)
	GetPages(book *models.Book) ([]*models.Page, error)
}
