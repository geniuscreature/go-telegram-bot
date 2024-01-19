package storage

import (
	"context"
	"database/sql"
	"github.com/geniuscreature/go-telegram-bot/internal/models"
	"time"
)

type SourceMysqlStorage struct {
	db *sql.DB
}

func NewSourceStorage(db *sql.DB) *SourceMysqlStorage {
	return &SourceMysqlStorage{
		db: db,
	}
}

func (s *SourceMysqlStorage) Sources(ctx context.Context) ([]models.Source, error) {
	rows, err := s.db.QueryContext(ctx, "select * from sources")
	if err != nil {
		return []models.Source{}, err
	}

	var sources []models.Source

	for rows.Next() {
		var source models.Source

		var createdAt, updatedAt []uint8

		if err = rows.Scan(
			&source.ID,
			&source.Name,
			&source.Url,
			&createdAt,
			&updatedAt,
		); err != nil {
			return []models.Source{}, nil
		}

		createdAtString, updatedAtString := string(createdAt), string(updatedAt)

		createdAtTime, err := time.Parse("2006-01-02 15:04:05", createdAtString)
		if err != nil {
			return []models.Source{}, err
		}

		updatedAtTime, err := time.Parse("2006-01-02 15:04:05", updatedAtString)
		if err != nil {
			return []models.Source{}, err
		}

		source.CreatedAt, source.UpdatedAt = createdAtTime, updatedAtTime

		sources = append(sources, source)
	}

	if rows.Err() != nil {
		return []models.Source{}, rows.Err()
	}

	return sources, nil
}

func (s *SourceMysqlStorage) SourceByID(ctx context.Context, id int64) (models.Source, error) {
	var source models.Source

	if err := s.db.QueryRowContext(
		ctx,
		"select * from sources where id = ?",
		id,
	).Scan(
		&source.Name,
		&source.Url,
		&source.CreatedAt,
		&source.UpdatedAt,
	); err != nil {
		return models.Source{}, err
	}

	return source, nil
}

func (s *SourceMysqlStorage) Add(ctx context.Context, source models.Source) (int64, error) {
	res, err := s.db.ExecContext(
		ctx,
		"insert into sources (name, url) values (?, ?)",
		source.Name,
		source.Url,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourceMysqlStorage) Delete(ctx context.Context, id int64) error {
	if _, err := s.db.ExecContext(
		ctx,
		"delete from sources where id = ?",
		id,
	); err != nil {
		return err
	}

	return nil
}
