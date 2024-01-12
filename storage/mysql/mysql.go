package mysql

import (
	"database/sql"
	"github.com/geniuscreature/go-telegram-bot/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Storage struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DatabaseConn)
	if err != nil {
		log.Fatalf("Couldn't connect to db: %s ", err)
	}

	return db, nil
}
