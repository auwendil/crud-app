package book

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/auwendil/crud-app/internal/models"
	"regexp"

	"testing"
)

var booksPostgresqlRows = []string{"id", "name", "author"}

func Test_Postgresql_GetAllBooks_ShouldReturnExpectedArray(t *testing.T) {
	// setup
	testServer, mock := prepareTestDB(t)
	defer testServer.DB.Close()

	// given
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

	// when
	resBooks, err := testServer.GetAllBooks()

	// then
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
	// setup
	testServer, mock := prepareTestDB(t)
	defer testServer.DB.Close()

	// given
	dbRows := sqlmock.NewRows(booksPostgresqlRows)
	mock.ExpectQuery(`SELECT id, name, author FROM books;`).
		WillReturnRows(dbRows)
	mock.ExpectCommit()

	// when
	resBooks, err := testServer.GetAllBooks()

	// then
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
	// setup
	testServer, mock := prepareTestDB(t)
	defer testServer.DB.Close()

	// given
	resultBookID := "3"
	expectedBook := &models.Book{ID: resultBookID, Name: "Book3", Author: "Author3"}

	dbRows := sqlmock.NewRows(booksPostgresqlRows)
	dbRows.AddRow(expectedBook.ID, expectedBook.Name, expectedBook.Author)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, author FROM books WHERE id=$1;`)).
		WithArgs(resultBookID).
		WillReturnRows(dbRows)
	mock.ExpectCommit()

	// when
	resBook, err := testServer.GetBook(resultBookID)

	// then
	if err != nil {
		t.Fatal(err)
	}

	if !bookEquals(resBook, expectedBook) {
		t.Errorf("Books not match: %+v vs %+v\n", resBook, expectedBook)
	}
}

func Test_Postgresql_AddBook_ShouldCallInsertQuery(t *testing.T) {
	// setup
	testServer, mock := prepareTestDB(t)
	defer testServer.DB.Close()

	// given
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
	// setup
	testServer, mock := prepareTestDB(t)
	defer testServer.DB.Close()

	// given
	testBook := &models.Book{ID: "3", Name: "Book3", Author: "Author3"}

	res := sqlmock.NewResult(1, 1)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET name = $2, author = $3 WHERE id = $1;`)).
		WithArgs(testBook.ID, testBook.Name, testBook.Author).
		WillReturnResult(res)
	mock.ExpectCommit()

	// when
	err := testServer.UpdateBook(testBook.ID, testBook)

	// then
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Postgresql_DeleteBook_ShouldCallDeleteQuery(t *testing.T) {
	// setup
	testServer, mock := prepareTestDB(t)
	defer testServer.DB.Close()

	// given
	testBook := &models.Book{ID: "3", Name: "Book3", Author: "Author3"}

	res := sqlmock.NewResult(1, 1)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1;`)).
		WithArgs(testBook.ID).
		WillReturnResult(res)
	mock.ExpectCommit()

	// when
	err := testServer.DeleteBook(testBook.ID)

	// then
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Postgresql_DeleteAllBooks_ShouldCallDeleteAllQuery(t *testing.T) {
	// setup
	testServer, mock := prepareTestDB(t)
	defer testServer.DB.Close()

	// given
	res := sqlmock.NewResult(1, 1)
	mock.ExpectExec(`DELETE FROM books;`).WillReturnResult(res)
	mock.ExpectCommit()

	// when
	err := testServer.DeleteAllBooks()

	// then
	if err != nil {
		t.Fatal(err)
	}
}

func prepareTestDB(t *testing.T) (*PostgreSQLRepo, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("[SETUP] Encountered error while preparing test db mock: %s\n", err)
	}

	testServer := &PostgreSQLRepo{DB: db}
	return testServer, mock
}
