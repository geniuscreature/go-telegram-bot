package source

import (
	"context"
	"github.com/SlyMarbo/rss"
	"github.com/geniuscreature/go-telegram-bot/internal/models"
)

type RSSSource struct {
	URL        string
	SourceID   int64
	SourceName string
}

func NewRSSSource(s models.Source) RSSSource {
	return RSSSource{
		URL:        s.Url,
		SourceID:   s.ID,
		SourceName: s.Name,
	}
}

func (s RSSSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := s.loadFeed(ctx, s.URL)
	if err != nil {
		return nil, err
	}

	items := make([]models.Item, 0, len(feed.Items))
	for _, v := range feed.Items {
		item := models.Item{
			Title:      v.Title,
			Categories: v.Categories,
			Link:       v.Link,
			Date:       v.Date,
			Summary:    v.Summary,
			SourceName: s.SourceName,
		}

		items = append(items, item)
	}

	return items, nil
}

func (s RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(url)

		if err != nil {
			errCh <- err
			return
		}

		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil
	}
}

func (s RSSSource) ID() int64 {
	return s.SourceID
}
