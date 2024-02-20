package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/sawitpro/UserService/repository"
)

type Client struct {
	DB *sql.DB
}

type ClientOptions struct {
	DSN string
}

func NewClient(opts ClientOptions) repository.RepositoryInterface {
	db, err := sql.Open("postgres", opts.DSN)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	return &Client{
		DB: db,
	}
}
