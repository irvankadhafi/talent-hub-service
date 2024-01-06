package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type (
	CandidateUsecase interface {
	}

	CandidateRepository interface {
		FindByID(ctx context.Context, id int64) (*Candidate, error)
		FindPasswordByID(ctx context.Context, id int64) ([]byte, error)
		FindByEmail(ctx context.Context, email string) (*Candidate, error)
		FindByPhone(ctx context.Context, phone string) (*Candidate, error)
		Create(ctx context.Context, candidate *Candidate) error
		Update(ctx context.Context, candidate *Candidate) error
	}

	Candidate struct {
		ID             int64
		FullName       string
		Email          string
		Phone          string
		Password       string
		DateOfBirth    time.Time
		Latitude       string
		Longitude      string
		Gender         Gender
		CityID         int64
		ProvinceID     int64
		LastEducation  time.Time
		LastExperience time.Time
		LoginDate      time.Time
		CreatedAt      time.Time
		UpdatedAt      time.Time
		DeletedAt      gorm.DeletedAt
	}
)

// Gender the candidate's gender
type Gender string

// Gender constants
const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)
