package model

import (
	"time"
)

type City struct {
	ID         int64
	ProvinceID int64
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
