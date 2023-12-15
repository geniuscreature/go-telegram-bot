package storage

import (
	"context"
	"github.com/geniuscreature/go-telegram-bot/internal/models"
	"github.com/geniuscreature/go-telegram-bot/storage/mysql"
)

type ArticleMysqlStorage struct {
	store mysql.Storage
}

func (s *ArticleMysqlStorage) Store(ctx context.Context, article models.Article) error {
	stmt, err := s.store.DB.Prepare("insert into articles (source_id, title, link, summary, published_at) values (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(
		ctx,
		article.SourceID,
		article.Title,
		article.Link,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *ArticleMysqlStorage) AllNotPosted(ctx context.Context) ([]models.Article, error) {
	stmt, err := s.store.DB.Prepare("select * from articles where posted_at is null >= current_timestamp order by published_at")

	if err != nil {
		return []models.Article{}, err
	}

	rows, err := stmt.QueryContext(
		ctx,
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
	stmt, err := s.store.DB.Prepare("update articles set posted_at = current_timestamp where id = ?")
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(
		ctx,
		id,
	); err != nil {
		return err
	}

	return nil
}
