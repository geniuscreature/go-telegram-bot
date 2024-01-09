package mysql

import (
	"database/sql"
	"github.com/geniuscreature/go-telegram-bot/internal/config"
	"log"
)

type Storage struct {
	DB *sql.DB
}

func New() (*sql.DB, error) {
	var cfg config.Config

	db, err := sql.Open("mysql", cfg.DatabaseConn)
	if err != nil {
		log.Fatalf("Couldn't connect to db: %s ", err)
	}

	return db, nil
}
