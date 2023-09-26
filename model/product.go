package model

import "time"

type Product struct {
	ID        int64
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
