package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Storage struct {
	DB *sql.DB
}

func New() (*Storage, error) {
	connString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatalf("Couldn't connect to db: %s ", err)
	}

	return &Storage{
		DB: db,
	}, nil
}
