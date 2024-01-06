package model

import (
	"gorm.io/gorm"
	"time"
)

type Experience struct {
	ID             int64
	CandidateID    int64
	CompanyName    string
	CompanyAddress string
	Position       string
	JobDescription string
	StartYear      time.Time
	EndYear        time.Time
	UntilNow       bool
	Flag           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}
