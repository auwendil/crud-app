package book

import (
	"github.com/auwendil/crud-app/internal/models"
	"testing"
)

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
