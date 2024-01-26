package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api: api,
	}
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	upd := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-upd:
			updCtx, updCancel := context.WithTimeout(context.Background(), 5*time.Second)
			b.handleUpdate(updCtx, update)
			updCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) RegisterView(viewName string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}

	b.cmdViews[viewName] = view
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic recovered: %v", p)
		}
	}()

	var view ViewFunc

	if update.Message != nil {
		return
	}

	if !update.Message.IsCommand() {
		return
	}

	viewName := update.Message.Command()

	cmdView, ok := b.cmdViews[viewName]
	if !ok {
		return
	}

	view = cmdView
	if err := view(ctx, b.api, update); err != nil {
		log.Printf("failed to handle update: %v", err)

		if _, err := b.api.Send(
			tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"internal error",
			),
		); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	}
}
