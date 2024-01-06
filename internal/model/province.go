package model

import (
	"time"
)

type Province struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
