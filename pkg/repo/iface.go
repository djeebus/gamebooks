package repo

import (
	"gamebooks/pkg/models"
)

type Repo interface {
	FindPagePath(book *models.Book, currentPagePath string, pageID string) (string, error)
	GetBooks() ([]*models.Book, error)
	GetBookByID(bookID string) (*models.Book, error)
	GetPage(bookID, pageID string) (*models.Page, error)
}
