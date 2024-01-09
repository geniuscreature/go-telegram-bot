package fetcher

import (
	"context"
	"github.com/geniuscreature/go-telegram-bot/internal/models"
	rssSource "github.com/geniuscreature/go-telegram-bot/internal/source"
	"log"
	"sync"
	"time"
)

type ArticleStorage interface {
	Store(ctx context.Context, article models.Article) error
}

type SourceProvider interface {
	Sources(ctx context.Context) ([]models.Source, error)
}

type Source interface {
	Fetch(ctx context.Context) ([]models.Item, error)
	ID() int64
}

type Fetcher struct {
	articles ArticleStorage
	sources  SourceProvider

	fetchInterval time.Duration
}

func New(
	articleStorage ArticleStorage,
	sourceProvider SourceProvider,
	fetchInterval time.Duration,
) *Fetcher {
	return &Fetcher{
		articles:      articleStorage,
		sources:       sourceProvider,
		fetchInterval: fetchInterval,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sources.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, source := range sources {
		wg.Add(1)

		go func(source Source) {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				log.Fatal("Items fetch failed")
				return
			}

			if err := f.saveItems(ctx, source, items); err != nil {
				log.Fatal("Items save failed")
				return
			}

		}(rssSource.NewRSSSource(source))
	}

	wg.Wait()

	return nil
}

func (f *Fetcher) saveItems(ctx context.Context, source Source, items []models.Item) error {
	var article models.Article

	for _, item := range items {
		article = models.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Link:        item.Link,
			Summary:     item.Summary,
			PublishedAt: item.Date,
		}

		if err := f.articles.Store(ctx, article); err != nil {
			return err
		}
	}

	return nil
}
