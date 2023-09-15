package main

import (
	"flag"
	"github.com/auwendil/crud-app/internal/repository"
	"github.com/auwendil/crud-app/internal/repository/book"
)

func prepareRepo(dbType, connString string) repository.BookRepo {
	var repo repository.BookRepo
	var err error

	switch dbType {
	case "postgresql":
		// TODO: take from config file / flags
		repo, err = book.NewPostgreSQLRepo(connString)
		if err != nil {
			panic(err)
		}
	case "mongodb":
		// TODO: take from config file / flags
		repo, err = book.NewMongoDBRepo(connString)
		if err != nil {
			panic(err)
		}
	}

	return repo
}

func main() {
	// default postgresql string: "postgresql://postgres:postgrespw@localhost:32768"
	// default mongodb string: "mongodb://localhost:27017"

	dbType := flag.String("db_type", "postgresql", "Type of database to use, available: [postgresql, mongodb]")
	connString := flag.String("conn_string", "postgresql://postgres:postgrespw@localhost:32768", "Connection string to chosen database")
	flag.Parse()

	repo := prepareRepo(*dbType, *connString)

	s := NewServer(":3000", repo)
	s.Start()
}
