package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TelegramBotToken     string `required:"true"`
	TelegramChannelID    int64  `required:"true"`
	DatabaseConn         string
	NotificationInterval time.Duration `default:"1m"`
	FetchInterval        time.Duration `default:"10m"`
}

func New() (*Config, error) {
	var cfg Config

	cfg.DatabaseConn = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	fetchInterval, err := parseInterval("FETCH_INTERVAL")
	if err != nil {
		log.Println(err)
	}
	cfg.FetchInterval = fetchInterval

	notificationInterval, err := parseInterval("NOTIFICATION_INTERVAL")
	if err != nil {
		log.Println(err)
	}
	cfg.NotificationInterval = notificationInterval

	cfg.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	cfg.TelegramChannelID, _ = strconv.ParseInt(os.Getenv("TELEGRAM_CHANNEL_ID"), 10, 64)

	return &cfg, nil
}

func parseInterval(interval string) (time.Duration, error) {
	fetchIntervalStr := os.Getenv(interval)
	fetchInterval, err := time.ParseDuration(fetchIntervalStr)
	if err != nil {
		return fetchInterval, errors.New(fmt.Sprintf("Could not parse %s value, setting to default", interval))
	}

	return fetchInterval, nil
}
