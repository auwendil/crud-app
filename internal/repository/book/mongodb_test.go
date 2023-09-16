package book

import (
	"github.com/auwendil/crud-app/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func Test_MongoDB_GetAllBooks(t *testing.T) {
	// setup
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("Should return expected array", func(mt *mtest.T) {
		// given
		ts := MongoDBRepo{
			collection: mt.Coll,
		}

		expectedBooks := []*models.Book{
			{ID: "111111111111111111111111", Name: "Book1", Author: "Author1"},
			{ID: "222222222222222222222222", Name: "Book2", Author: "Author2"},
			{ID: "333333333333333333333333", Name: "Book3", Author: "Author3"},
			{ID: "444444444444444444444444", Name: "Book4", Author: "Author4"},
		}

		resBooks := []bson.D{
			createBsonForBook(t, expectedBooks[0]),
			createBsonForBook(t, expectedBooks[1]),
			createBsonForBook(t, expectedBooks[2]),
			createBsonForBook(t, expectedBooks[3]),
		}
		cursorResponses := []bson.D{
			mtest.CreateCursorResponse(1, "db.books", mtest.FirstBatch, resBooks[0]),
			mtest.CreateCursorResponse(1, "db.books", mtest.NextBatch, resBooks[1]),
			mtest.CreateCursorResponse(1, "db.books", mtest.NextBatch, resBooks[2]),
			mtest.CreateCursorResponse(1, "db.books", mtest.NextBatch, resBooks[3]),
			mtest.CreateCursorResponse(0, "db.books", mtest.NextBatch),
		}
		mt.AddMockResponses(cursorResponses...)

		// when
		books, err := ts.GetAllBooks()

		// then
		if err != nil {
			t.Fatal("Encountered error while retrieving books from db:", err)
		}

		if !bookArraysEquals(t, books, expectedBooks) {
			t.Fatalf("Returned book array are not equal to expected book array: %v vs %v\n", books, expectedBooks)
		}
	})

	mt.Run("Should return empty array and not nil when there is no books", func(mt *mtest.T) {
		// given
		ts := MongoDBRepo{
			collection: mt.Coll,
		}

		cursorResponses := []bson.D{
			mtest.CreateCursorResponse(0, "db.books", mtest.FirstBatch),
		}

		mt.AddMockResponses(cursorResponses...)

		// when
		books, err := ts.GetAllBooks()

		// then
		if err != nil {
			t.Fatal("Encountered error while retrieving books from db:", err)
		}

		if books == nil {
			t.Fatal("Returns nil but instead should return empty array")
		}

		if len(books) > 0 {
			t.Fatal("Returned some books but there should be none:", books)
		}
	})
}

func Test_MongoDB_GetBook(t *testing.T) {
	// setup
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("Should return expected book", func(mt *mtest.T) {
		// given
		ts := MongoDBRepo{
			collection: mt.Coll,
		}

		expectedBook := &models.Book{ID: "333333333333333333333333", Name: "Book3", Author: "Author3"}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.books", mtest.FirstBatch, createBsonForBook(t, expectedBook)))

		// when
		book, err := ts.GetBook(expectedBook.ID)

		// then
		if err != nil {
			t.Fatalf("Encountered error while retrieving book (id=%s) from db: %s\n", expectedBook.ID, err)
		}

		if !bookEquals(book, expectedBook) {
			t.Fatalf("Result book are not equal to expected: %v vs %v\n", book, expectedBook)
		}
	})

	mt.Run("Should return error with invalid id", func(mt *mtest.T) {
		// given
		ts := MongoDBRepo{
			collection: mt.Coll,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.books", mtest.FirstBatch, nil))

		notExistingID := "111111111111111111111111"

		// when
		book, err := ts.GetBook(notExistingID)

		// then
		if err == nil {
			t.Fatal("Expected to return error but returned nil instead")
		}

		if book != nil {
			t.Fatal("Returned book should be nil")
		}
	})
}

func Test_MongoDB_AddBook(t *testing.T) {
	// setup
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("Should create book", func(mt *mtest.T) {
		// given
		ts := MongoDBRepo{
			collection: mt.Coll,
		}

		createdBook := &models.Book{ID: "333333333333333333333333", Name: "Book3", Author: "Author3"}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// when
		book, err := ts.AddBook(createdBook)

		// then
		if err != nil {
			t.Fatalf("Encountered error while creating book (id=%s): %s\n", createdBook.ID, err)
		}

		if book == nil {
			t.Fatalf("Returned book should not be nil\n")
		}
	})
}

func Test_MongoDB_UpdateBook(t *testing.T) {
	// setup
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("Should update book", func(mt *mtest.T) {
		// given
		ts := MongoDBRepo{
			collection: mt.Coll,
		}

		updatedBook := &models.Book{ID: "333333333333333333333333", Name: "Book3", Author: "Author3"}

		updated := bson.D{{"ok", 1}, {"value", createBsonForBook(t, updatedBook)}}
		mt.AddMockResponses(updated)

		// when
		err := ts.UpdateBook(updatedBook.ID, updatedBook)

		// then
		if err != nil {
			t.Fatalf("Encountered error while updating book (id=%s): %s\n", updatedBook.ID, err)
		}
	})
}

func Test_MongoDB_DeleteBook(t *testing.T) {
	// setup
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("Should delete book", func(mt *mtest.T) {
		// given
		ts := MongoDBRepo{
			collection: mt.Coll,
		}

		deletedBook := &models.Book{ID: "333333333333333333333333", Name: "Book3", Author: "Author3"}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// when
		err := ts.DeleteBook(deletedBook.ID)

		// then
		if err != nil {
			t.Fatalf("Encountered error while updating book (id=%s): %s\n", deletedBook.ID, err)
		}
	})

	mt.Run("Should return error while trying to remove not existing book", func(mt *mtest.T) {
		// given
		ts := MongoDBRepo{
			collection: mt.Coll,
		}

		notExistingID := "111111111111111111111111"
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    1,
			Message: "key not exists error",
		}))

		// when
		err := ts.DeleteBook(notExistingID)

		// then
		if err == nil {
			t.Fatal("Expected to return error but returned nil instead")
		}
	})
}

func Test_MongoDB_DeleteAllBooks_ShouldCallDeleteAllQuery(t *testing.T) {
	// setup
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("Should delete all books", func(mt *mtest.T) {
		// given
		ts := MongoDBRepo{
			collection: mt.Coll,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// when
		err := ts.DeleteAllBooks()

		// then
		if err != nil {
			t.Fatalf("Encountered error while deleting all books: %s\n", err)
		}
	})
}

// utility functions

func createBsonForBook(t *testing.T, b *models.Book) bson.D {
	objID, err := primitive.ObjectIDFromHex(b.ID)
	if err != nil {
		t.Fatalf("Encountered error while parsing test data: %s\n", err)
	}

	parsedBook := bson.D{
		{"_id", objID},
		{"name", b.Name},
		{"author", b.Author},
	}

	return parsedBook
}
