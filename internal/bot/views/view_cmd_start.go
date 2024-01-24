package views

import (
	"context"
	"github.com/geniuscreature/go-telegram-bot/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdStart() bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(
			tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Started",
			),
		); err != nil {
			return err
		}

		return nil
	}
}
