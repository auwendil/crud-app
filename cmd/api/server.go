package main

import (
	"github.com/auwendil/crud-app/internal/repository"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	addr   string
	dbRepo repository.BookRepo
}

func NewServer(listenAddr string, repo repository.BookRepo) *Server {
	return &Server{
		addr:   listenAddr,
		dbRepo: repo,
	}
}

func (s *Server) Start() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/book", s.handleGetAllBooks)
	r.Get("/book/{id}", s.handleGetBook)
	r.Post("/book", s.handleAddBook)
	r.Put("/book/{id}", s.handleUpdateBook)
	r.Delete("/book/{id}", s.handleDeleteBook)
	r.Delete("/book", s.handleDeleteAll)

	log.Println("Starting server on", s.addr)
	if err := http.ListenAndServe(s.addr, r); err != nil {
		panic(err)
	}
}
