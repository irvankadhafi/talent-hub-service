package model

import (
	"gorm.io/gorm"
	"time"
)

type Education struct {
	ID              int64
	CandidateID     int64
	InstitutionName string
	Major           string
	StartYear       time.Time
	EndYear         time.Time
	UntilNow        bool
	GPA             float64
	Flag            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}
