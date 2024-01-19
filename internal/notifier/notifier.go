package notifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/geniuscreature/go-telegram-bot/internal/models"
	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
	"time"
)

type ArticleProvider interface {
	AllNotPosted(ctx context.Context, timestamp time.Time, limit int64) ([]models.Article, error)
	MarkPosted(ctx context.Context, id int64) error
}

type Notifier struct {
	articles   ArticleProvider
	bot        *tgbotapi.BotAPI
	interval   time.Duration
	timePeriod time.Duration
	channelID  int64
}

func New(
	articles ArticleProvider,
	bot *tgbotapi.BotAPI,
	interval time.Duration,
	timePeriod time.Duration,
	channelID int64,
) *Notifier {
	return &Notifier{
		articles:   articles,
		bot:        bot,
		interval:   interval,
		timePeriod: timePeriod,
		channelID:  channelID,
	}
}

func (n *Notifier) Start(ctx context.Context) error {
	ticker := time.NewTicker(n.interval)
	defer ticker.Stop()

	if err := n.PostArticle(ctx); err != nil {
		return errors.New(fmt.Sprintf("post article error %v", err))
	}

	for {
		select {
		case <-ticker.C:
			if err := n.PostArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) PostArticle(ctx context.Context) error {
	topArticle, err := n.articles.AllNotPosted(ctx, time.Now().Add(-n.timePeriod), 1)
	if err != nil {
		return err
	}

	if len(topArticle) == 0 {
		return nil
	}

	article := topArticle[0]

	summary, err := n.extractSummary(ctx, article)
	if err != nil {
		return err
	}

	if err := n.sendArticle(article, summary); err != nil {
		return err
	}

	return nil
}

func (n *Notifier) extractSummary(ctx context.Context, article models.Article) (string, error) {
	var r io.Reader

	if article.Summary == "" {
		r = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Link)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		r = resp.Body
	}

	summary, err := readability.FromReader(r, nil)
	if err != nil {
		return "", err
	}

	return "\n\n" + cleanText(summary.TextContent), nil
}

func (n *Notifier) sendArticle(article models.Article, summary string) error {
	msg := tgbotapi.NewMessage(
		n.channelID,
		fmt.Sprintf(
			article.Title,
			summary,
			article.Link,
		),
	)

	_, err := n.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func cleanText(text string) string {
	return strings.TrimSpace(text)
}
