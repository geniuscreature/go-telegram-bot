package models

import "time"

type Source struct {
	ID        int64
	Name      string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
