package main

import (
	"github.com/geniuscreature/go-telegram-bot/storage/mysql"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("err loading env: %v", err)
	}

	_, err := mysql.New()
	if err != nil {
		log.Fatalf("couldn't connect to db: %v", err)
	}
}
