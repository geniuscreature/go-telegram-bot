package storage

import (
	"context"
	"github.com/geniuscreature/go-telegram-bot/internal/models"
	"github.com/geniuscreature/go-telegram-bot/storage/mysql"
)

type SourceMysqlStorage struct {
	store mysql.Storage
}

func (s *SourceMysqlStorage) Sources(ctx context.Context) ([]models.Source, error) {
	stmt, err := s.store.DB.Prepare("select * from sources")
	if err != nil {
		return []models.Source{}, err
	}

	rows, err := stmt.QueryContext(
		ctx,
	)
	if err != nil {
		return []models.Source{}, err
	}

	var sources []models.Source

	for rows.Next() {
		var source models.Source
		if err = rows.Scan(
			&source.ID,
			&source.Name,
			&source.Url,
			&source.CreatedAt,
			&source.UpdatedAt,
		); err != nil {
			return []models.Source{}, nil
		}

		sources = append(sources, source)
	}

	if rows.Err() != nil {
		return []models.Source{}, rows.Err()
	}

	return sources, nil
}

func (s *SourceMysqlStorage) SourceByID(ctx context.Context, id int64) (models.Source, error) {
	stmt, err := s.store.DB.Prepare("select * from sources where id = ?")
	if err != nil {
		return models.Source{}, err
	}

	var source models.Source

	if err := stmt.QueryRowContext(
		ctx,
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
	stmt, err := s.store.DB.Prepare("insert into sources (name, url) values (?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, source.Name, source.Url)
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
	stmt, err := s.store.DB.Prepare("delete from sources where id = ?")
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, id); err != nil {
		return err
	}

	return nil
}
