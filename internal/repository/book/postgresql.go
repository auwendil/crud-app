package book

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/auwendil/crud-app/internal/models"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type PostgreSQLRepo struct {
	DB *sql.DB
}

const postgresDBTimeout = time.Second * 3
const postgresDBDriverName = "pgx"
const postgresDBName = "books"

func NewPostgreSQLRepo(dataName string) (*PostgreSQLRepo, error) {
	db, err := sql.Open(postgresDBDriverName, fmt.Sprintf("%s/%s", dataName, postgresDBName))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	repo := &PostgreSQLRepo{
		DB: db,
	}

	return repo, nil
}

func (r *PostgreSQLRepo) GetAllBooks() ([]*models.Book, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), postgresDBTimeout)
	defer cancelFn()

	query := `
		SELECT id, name, author
		FROM books;
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*models.Book{}
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Name, &book.Author)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	return books, nil
}

func (r *PostgreSQLRepo) GetBook(id string) (*models.Book, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), postgresDBTimeout)
	defer cancelFn()

	query := `
		SELECT id, name, author
		FROM books
		WHERE id=$1;
	`

	var book models.Book
	row := r.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&book.ID, &book.Name, &book.Author)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *PostgreSQLRepo) AddBook(b *models.Book) (*models.Book, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), postgresDBTimeout)
	defer cancelFn()

	query := `
		INSERT INTO books (name, author)
		VALUES ($1, $2)
		RETURNING id;
	`

	var newId int
	err := r.DB.QueryRowContext(ctx, query, b.Name, b.Author).Scan(&newId)
	if err != nil {
		return nil, err
	}

	createdBook := &models.Book{
		ID:     fmt.Sprintf("%d", newId),
		Name:   b.Name,
		Author: b.Author,
	}
	return createdBook, nil
}

func (r *PostgreSQLRepo) UpdateBook(id string, updatedBook *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), postgresDBTimeout)
	defer cancel()

	query := `
		UPDATE books
		SET name = $2, author = $3
		WHERE id = $1;
	`

	_, err := r.DB.ExecContext(ctx, query, id, updatedBook.Name, updatedBook.Author)
	return err
}

func (r *PostgreSQLRepo) DeleteBook(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), postgresDBTimeout)
	defer cancel()

	query := `
		DELETE FROM books
		WHERE id = $1;
	`

	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

func (r *PostgreSQLRepo) DeleteAllBooks() error {
	ctx, cancel := context.WithTimeout(context.Background(), postgresDBTimeout)
	defer cancel()

	query := `
		DELETE FROM books;
	`

	_, err := r.DB.ExecContext(ctx, query)
	return err
}
