package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/auwendil/crud-app/internal/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

var storedBooks = []*models.Book{
	{ID: "1", Name: "Name1", Author: "Author1"},
	{ID: "2", Name: "Name2", Author: "Author2"},
	{ID: "3", Name: "Name3", Author: "Author3"},
}

func Test_Server_HandleGetBooks(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/book", nil)
	w := httptest.NewRecorder()

	storageSize := 3
	ts := &Server{
		dbRepo: prepareDbRepoMock(storageSize),
	}

	// when
	ts.handleGetAllBooks(w, req)

	// then
	httpResponse := w.Result()
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d(%s) but received: %d(%s)\n",
			http.StatusOK, http.StatusText(http.StatusOK),
			httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
	}

	jsonResponse := parseHttpResponse(t, httpResponse)
	if jsonResponse.Error {
		t.Fatal("Encountered error but there should be none:", jsonResponse)
	}

	receivedBooks := getBooksFromResponse(t, jsonResponse.Data)

	if len(receivedBooks) != storageSize {
		t.Fatalf("Returned books array has wrong size: has %d, should be: %d\n", len(receivedBooks), storageSize)
	}

	if !bookArraysEquals(t, receivedBooks, storedBooks) {
		t.Fatalf("Received books not match stored, has: %v, should be: %v\n", receivedBooks, storedBooks)
	}
}

func Test_Server_HandleGetBookById(t *testing.T) {
	// setup
	expectedBook := storedBooks[0]
	storageSize := 3
	ts := &Server{
		dbRepo: prepareDbRepoMock(storageSize),
	}

	t.Run("Returns expected book with correct id", func(t *testing.T) {
		// given
		req := httptest.NewRequest(http.MethodGet, "/book/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", expectedBook.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		// when
		ts.handleGetBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusOK, http.StatusText(http.StatusOK),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if jsonResponse.Error {
			t.Fatal("Encountered error but there should be none:", jsonResponse)
		}

		receivedBook := getBookFromResponse(t, jsonResponse.Data)

		if !bookEquals(receivedBook, expectedBook) {
			t.Fatalf("Received books not match stored, has: %v, should be: %v\n", receivedBook, storedBooks)
		}
	})

	t.Run("Returns error when book is not found", func(t *testing.T) {
		// given
		notExistingID := "-1"
		req := httptest.NewRequest(http.MethodGet, "/book/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", notExistingID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		// when
		ts.handleGetBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusNotFound, http.StatusText(http.StatusNotFound),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if !jsonResponse.Error {
			t.Fatal("Response does not contain error but should be one:", jsonResponse)
		}
	})
}

func Test_Server_HandleAddBook(t *testing.T) {
	// setup
	storageSize := 3
	ts := &Server{
		dbRepo: prepareDbRepoMock(storageSize),
	}

	t.Run("Should create new book if not exists", func(t *testing.T) {
		// given
		newBook := &models.Book{ID: "5", Name: "NewBook", Author: "NewAuthor"}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(newBook)
		if err != nil {
			t.Fatalf("Encountered error while encoding test data: %s\n", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/book", &buf)
		w := httptest.NewRecorder()

		// when
		ts.handleAddBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusCreated {
			t.Errorf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusCreated, http.StatusText(http.StatusCreated),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if jsonResponse.Error {
			t.Fatal("Encountered error but there should be none:", jsonResponse)
		}
	})

	t.Run("Should create new book if not exists", func(t *testing.T) {
		// given
		existingBook := storedBooks[0]
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(existingBook)
		if err != nil {
			t.Fatalf("Encountered error while encoding test data: %s\n", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/book", &buf)
		w := httptest.NewRecorder()

		// when
		ts.handleAddBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if !jsonResponse.Error {
			t.Fatal("Does not encountered error but there should be one:", jsonResponse)
		}
	})

	t.Run("Should fail with malformed payload", func(t *testing.T) {
		// given
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(struct {
			ID int `json:"id"`
		}{-1})
		if err != nil {
			t.Fatalf("Encountered error while encoding test data: %s\n", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/book", &buf)
		w := httptest.NewRecorder()

		// when
		ts.handleAddBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusBadRequest, http.StatusText(http.StatusBadRequest),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if !jsonResponse.Error {
			t.Fatal("Does not encountered error but there should be one:", jsonResponse)
		}
	})
}

func Test_Server_HandleUpdateBook(t *testing.T) {
	// setup
	storageSize := 3
	ts := &Server{
		dbRepo: prepareDbRepoMock(storageSize),
	}

	t.Run("Should update book", func(t *testing.T) {
		// given
		updatedName := "UpdatedBookName"
		updatedAuthor := "UpdatedBookAuthor"
		updatedBook := &models.Book{ID: storedBooks[0].ID, Name: updatedName, Author: updatedAuthor}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(updatedBook)
		if err != nil {
			t.Fatalf("Encountered error while encoding test data: %s\n", err)
		}

		req := httptest.NewRequest(http.MethodPut, "/book/{id}", &buf)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", storedBooks[0].ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		// when
		ts.handleUpdateBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusNoContent, http.StatusText(http.StatusNoContent),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if jsonResponse.Error {
			t.Fatalf("Encountered error but there should be none: %+v\n", jsonResponse)
		}
	})

	t.Run("Should not update not existing book", func(t *testing.T) {
		// given
		updatedName := "UpdatedBookName"
		updatedAuthor := "UpdatedBookAuthor"
		notExistingId := "-1"
		updatedBook := &models.Book{ID: notExistingId, Name: updatedName, Author: updatedAuthor}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(updatedBook)
		if err != nil {
			t.Fatalf("Encountered error while encoding test data: %s\n", err)
		}

		req := httptest.NewRequest(http.MethodPut, "/book", &buf)
		w := httptest.NewRecorder()

		// when
		ts.handleUpdateBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if !jsonResponse.Error {
			t.Fatalf("Does not encountered error but there should be one: %+v", jsonResponse)
		}
	})

	t.Run("Should fail with malformed payload", func(t *testing.T) {
		// given
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(struct {
			ID int `json:"id"`
		}{-1})
		if err != nil {
			t.Fatalf("Encountered error while encoding test data: %s\n", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/book", &buf)
		w := httptest.NewRecorder()

		// when
		ts.handleUpdateBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusBadRequest, http.StatusText(http.StatusBadRequest),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if !jsonResponse.Error {
			t.Fatalf("Does not encountered error but there should be one: %+v\n", jsonResponse)
		}
	})
}

func Test_Server_HandleDeleteBook(t *testing.T) {
	// setup
	storageSize := 3
	ts := &Server{
		dbRepo: prepareDbRepoMock(storageSize),
	}

	t.Run("Should delete book", func(t *testing.T) {
		// given
		deletedBook := storedBooks[0]
		req := httptest.NewRequest(http.MethodDelete, "/book/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", deletedBook.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		// when
		ts.handleDeleteBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusNoContent, http.StatusText(http.StatusNoContent),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if jsonResponse.Error {
			t.Fatalf("Encountered error but there should be none: %+v\n", jsonResponse)
		}

		if book, _ := ts.dbRepo.GetBook(deletedBook.ID); book != nil {
			t.Fatalf("Book(%s) should be deleted but are not\n", deletedBook.ID)
		}
	})

	t.Run("Should not delete book if not exists", func(t *testing.T) {
		// given
		notExistingID := "-1"
		req := httptest.NewRequest(http.MethodDelete, "/book/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", notExistingID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		// when
		ts.handleDeleteBook(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusBadRequest, http.StatusText(http.StatusBadRequest),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if !jsonResponse.Error {
			t.Fatalf("Does not encountered error but there should be one: %+v\n", jsonResponse)
		}
	})
}

func Test_Server_HandleDeleteAllBooks(t *testing.T) {
	// setup
	storageSize := 3
	ts := &Server{
		dbRepo: prepareDbRepoMock(storageSize),
	}

	t.Run("Should delete all book", func(t *testing.T) {
		// given
		req := httptest.NewRequest(http.MethodDelete, "/book", nil)
		w := httptest.NewRecorder()

		// when
		ts.handleDeleteAll(w, req)

		// then
		httpResponse := w.Result()
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status %d(%s) but received: %d(%s)\n",
				http.StatusNoContent, http.StatusText(http.StatusNoContent),
				httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode))
		}

		jsonResponse := parseHttpResponse(t, httpResponse)
		if jsonResponse.Error {
			t.Fatalf("Encountered error but there should be none: %+v\n", jsonResponse)
		}

		if books, _ := ts.dbRepo.GetAllBooks(); len(books) > 0 {
			t.Fatalf("All books should be removed from repo, but there are still %d available\n", len(books))
		}
	})
}

// utils

type dbRepoMock struct {
	m map[string]*models.Book
}

func (r *dbRepoMock) GetAllBooks() ([]*models.Book, error) {
	var books = []*models.Book{}
	for _, v := range r.m {
		books = append(books, v)
	}
	return books, nil
}

func (r *dbRepoMock) GetBook(id string) (*models.Book, error) {
	b, ok := r.m[id]
	if !ok {
		return nil, fmt.Errorf("book (id=%s) not found", id)
	}
	return b, nil
}

func (r *dbRepoMock) AddBook(b *models.Book) (*models.Book, error) {
	if _, ok := r.m[b.ID]; ok {
		return nil, fmt.Errorf("Book already exists")
	}

	r.m[b.ID] = b
	return b, nil
}

func (r *dbRepoMock) UpdateBook(id string, updatedBook *models.Book) error {
	if _, ok := r.m[id]; !ok {
		return fmt.Errorf("Book does not exist")
	}
	r.m[id] = updatedBook
	return nil
}

func (r *dbRepoMock) DeleteBook(id string) error {
	if _, ok := r.m[id]; !ok {
		return fmt.Errorf("Book does not exist")
	}
	delete(r.m, id)
	return nil
}

func (r *dbRepoMock) DeleteAllBooks() error {
	r.m = make(map[string]*models.Book)
	return nil
}

func prepareDbRepoMock(amountOfBooksLoaded int) *dbRepoMock {
	repo := &dbRepoMock{
		m: make(map[string]*models.Book),
	}

	loadedBooks := 0
	for _, book := range storedBooks {
		if loadedBooks >= amountOfBooksLoaded {
			break
		}
		repo.m[book.ID] = book
	}
	return repo
}

func parseHttpResponse(t *testing.T, res *http.Response) JSONResponse {
	var response JSONResponse
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatal("Encountered error while decoding json response:", err)
	}
	return response
}

func getBooksFromResponse(t *testing.T, data interface{}) []*models.Book {
	parsedData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Encountered error while marshalling received data (%+v): %s\n", data, err)
	}

	var books []*models.Book
	err = json.Unmarshal(parsedData, &books)
	if err != nil {
		t.Fatalf("Encountered error while unmarshalling received data (%+v): %s\n", data, err)
	}
	return books
}

func getBookFromResponse(t *testing.T, data interface{}) *models.Book {
	parsedData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Encountered error while marshalling received data (%+v): %s\n", data, err)
	}

	var book *models.Book
	err = json.Unmarshal(parsedData, &book)
	if err != nil {
		t.Fatalf("Encountered error while unmarshalling received data (%+v): %s\n", data, err)
	}
	return book
}

func bookEquals(a, b *models.Book) bool {
	return a.ID == b.ID && a.Name == b.Name && a.Author == b.Author
}

func bookArraysEquals(t *testing.T, resultArr, expectedArr []*models.Book) bool {
	if len(resultArr) != len(expectedArr) {
		t.Errorf("Different size of result (%d) and expected (%d) arrays\n", len(resultArr), len(expectedArr))
		return false
	}

	areEquals := true
	for i := 0; i < len(expectedArr); i++ {
		resultBook := expectedArr[i]
		expectedBook := findBook(resultArr, resultBook.ID)

		if expectedBook == nil {
			t.Errorf("Result books array does not contain expected book: %+v\n", expectedBook)
			areEquals = false
		}

		if !bookEquals(resultBook, expectedBook) {
			t.Errorf("Books not match: %+v vs %+v\n", resultBook, expectedBook)
			areEquals = false
		}
	}
	return areEquals
}

func findBook(arr []*models.Book, id string) *models.Book {
	for i := 0; i < len(arr); i++ {
		if arr[i].ID == id {
			return arr[i]
		}
	}
	return nil
}
