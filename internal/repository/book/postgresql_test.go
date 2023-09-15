package book

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/auwendil/crud-app/internal/models"
	"regexp"

	"testing"
)

var booksPostgresqlRows = []string{"id", "name", "author"}

func Test_Postgresql_GetAllBooks_ShouldReturnExpectedArray(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	expectedBooks := []*models.Book{
		{ID: "1", Name: "Book1", Author: "Author1"},
		{ID: "2", Name: "Book2", Author: "Author2"},
		{ID: "3", Name: "Book3", Author: "Author3"},
		{ID: "4", Name: "Book4", Author: "Author4"},
	}

	dbRows := sqlmock.NewRows(booksPostgresqlRows)
	for _, book := range expectedBooks {
		dbRows.AddRow(book.ID, book.Name, book.Author)
	}

	mock.ExpectQuery(`SELECT id, name, author FROM books;`).
		WillReturnRows(dbRows)
	mock.ExpectCommit()

	testServer := PostgreSQLRepo{DB: db}
	resBooks, err := testServer.GetAllBooks()
	if err != nil {
		t.Error(err)
	}

	if !bookArraysEquals(t, resBooks, expectedBooks) {
		t.Error("Result books are different than expected")
	}

	for i := 0; i < len(resBooks); i++ {
		if !bookEquals(resBooks[i], expectedBooks[i]) {
			t.Errorf("Books not match: %+v vs %+v\n", resBooks[i], expectedBooks[i])
		}
	}
}

func Test_Postgresql_GetAllBooks_ShouldReturnEmptyArrayWhenTableIsEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	dbRows := sqlmock.NewRows(booksPostgresqlRows)

	mock.ExpectQuery(`SELECT id, name, author FROM books;`).
		WillReturnRows(dbRows)
	mock.ExpectCommit()

	testServer := PostgreSQLRepo{DB: db}
	resBooks, err := testServer.GetAllBooks()
	if err != nil {
		t.Error(err)
	}

	if resBooks == nil {
		t.Error("Result should be an empty array and not nil value")
	}

	if len(resBooks) > 0 {
		t.Errorf("Result array should be empty but has length of %d\n", len(resBooks))
	}
}

func Test_Postgresql_GetBook_ShouldCallSelectQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testServer := PostgreSQLRepo{DB: db}

	expectedBook := &models.Book{ID: "3", Name: "Book3", Author: "Author3"}

	dbRows := sqlmock.NewRows(booksPostgresqlRows)
	dbRows.AddRow(expectedBook.ID, expectedBook.Name, expectedBook.Author)

	resultBookID := "3"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, author FROM books WHERE id=$1;`)).
		WithArgs(resultBookID).
		WillReturnRows(dbRows)
	mock.ExpectCommit()

	resBook, err := testServer.GetBook(resultBookID)
	if err != nil {
		t.Fatal(err)
	}

	if !bookEquals(resBook, expectedBook) {
		t.Errorf("Books not match: %+v vs %+v\n", resBook, expectedBook)
	}
}

func Test_Postgresql_AddBook_ShouldCallInsertQuery(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testServer := PostgreSQLRepo{DB: db}

	testBook := &models.Book{ID: "3", Name: "Book3", Author: "Author3"}

	dbRows := sqlmock.NewRows([]string{"id"})
	dbRows.AddRow(testBook.ID)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO books (name, author) VALUES ($1, $2) RETURNING id;`)).
		WithArgs(testBook.Name, testBook.Author).
		WillReturnRows(dbRows)
	mock.ExpectCommit()

	//when
	book, err := testServer.AddBook(testBook)

	//then
	if err != nil {
		t.Fatal(err)
	}

	if book.ID != testBook.ID {
		t.Fatalf("Returned id (%s) is different than expected: %s\n", book.ID, testBook.ID)
	}
}

func Test_Postgresql_UpdateBook_ShouldCallUpdateQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testServer := PostgreSQLRepo{DB: db}

	testBook := &models.Book{ID: "3", Name: "Book3", Author: "Author3"}

	res := sqlmock.NewResult(1, 1)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET name = $2, author = $3 WHERE id = $1;`)).
		WithArgs(testBook.ID, testBook.Name, testBook.Author).
		WillReturnResult(res)
	mock.ExpectCommit()

	err = testServer.UpdateBook(testBook.ID, testBook)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Postgresql_DeleteBook_ShouldCallDeleteQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testServer := PostgreSQLRepo{DB: db}

	testBook := &models.Book{ID: "3", Name: "Book3", Author: "Author3"}

	res := sqlmock.NewResult(1, 1)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1;`)).
		WithArgs(testBook.ID).
		WillReturnResult(res)
	mock.ExpectCommit()

	err = testServer.DeleteBook(testBook.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Postgresql_DeleteAllBooks_ShouldCallDeleteAllQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testServer := PostgreSQLRepo{DB: db}

	res := sqlmock.NewResult(1, 1)

	mock.ExpectExec(`DELETE FROM books;`).WillReturnResult(res)
	mock.ExpectCommit()

	err = testServer.DeleteAllBooks()
	if err != nil {
		t.Fatal(err)
	}
}

// utility functions
