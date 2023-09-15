package repository

import (
	"github.com/auwendil/crud-app/internal/models"
)

type BookRepo interface {
	GetAllBooks() ([]*models.Book, error)
	GetBook(id string) (*models.Book, error)
	AddBook(b *models.Book) (*models.Book, error)
	UpdateBook(id string, updatedBook *models.Book) error
	DeleteBook(id string) error
	DeleteAllBooks() error
}
