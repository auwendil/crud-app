package main

import (
	"encoding/json"
	"github.com/auwendil/crud-app/internal/models"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) handleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := s.dbRepo.GetAllBooks()
	if err != nil {
		_ = handleErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = handleSuccessfulJSON(w, "", books, http.StatusOK)
}

func (s *Server) handleGetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book, err := s.dbRepo.GetBook(id)
	if err != nil {
		_ = handleErrorJSON(w, err, http.StatusNotFound)
		return
	}

	_ = handleSuccessfulJSON(w, "", book, http.StatusOK)
}

func (s *Server) handleAddBook(w http.ResponseWriter, r *http.Request) {
	var book *models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		_ = handleErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	book, err := s.dbRepo.AddBook(book)
	if err != nil {
		_ = handleErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = handleSuccessfulJSON(w, "", book, http.StatusCreated)
}

func (s *Server) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var book *models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		_ = handleErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err := s.dbRepo.UpdateBook(id, book)
	if err != nil {
		_ = handleErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = handleSuccessfulJSON(w, "", book, http.StatusNoContent)
}

func (s *Server) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := s.dbRepo.DeleteBook(id)
	if err != nil {
		_ = handleErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	_ = handleSuccessfulJSON(w, "", nil, http.StatusNoContent)
}

func (s *Server) handleDeleteAll(w http.ResponseWriter, r *http.Request) {
	err := s.dbRepo.DeleteAllBooks()
	if err != nil {
		_ = handleErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = handleSuccessfulJSON(w, "", nil, http.StatusNoContent)
}
