package storage

import (
	"context"
	"database/sql"
	"github.com/geniuscreature/go-telegram-bot/internal/models"
	"time"
)

type ArticleMysqlStorage struct {
	db *sql.DB
}

func NewArticleStorage(db *sql.DB) *ArticleMysqlStorage {
	return &ArticleMysqlStorage{
		db: db,
	}
}

func (s *ArticleMysqlStorage) Store(ctx context.Context, article models.Article) error {
	if err := s.db.QueryRowContext(
		ctx,
		"insert into articles (source_id, title, link, summary, published_at) values (?, ?, ?, ?, ?)",
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (s *ArticleMysqlStorage) AllNotPosted(ctx context.Context, timestamp time.Time, limit int64) ([]models.Article, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
		select * 
		from articles 
		where 
		    posted_at is null and 
		    published_at >= ? 
		order by created_at
		limit ?`,
		timestamp,
		limit,
	)
	if err != nil {
		return []models.Article{}, err
	}

	var articles []models.Article

	for rows.Next() {
		var article models.Article
		if err = rows.Scan(
			&article.ID,
			&article.SourceID,
			&article.Title,
			&article.Link,
			&article.Summary,
			&article.PublishedAt,
			&article.CreatedAt,
			&article.PostedAt,
		); err != nil {
			return []models.Article{}, nil
		}

		articles = append(articles, article)
	}

	if rows.Err() != nil {
		return []models.Article{}, rows.Err()
	}

	return articles, nil
}

func (s *ArticleMysqlStorage) MarkPosted(ctx context.Context, id int64) error {
	if err := s.db.QueryRowContext(
		ctx,
		"update articles set posted_at = current_timestamp where id = ?",
		id,
	); err != nil {
		return err.Err()
	}

	return nil
}
