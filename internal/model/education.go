package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type (
	Education struct {
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

	EducationRepository interface {
		FindByID(ctx context.Context, id int64) (*Education, error)
		FindByCandidateID(ctx context.Context, candidateID int64) (*Education, error)
		Create(ctx context.Context, education *Education) error
		Update(ctx context.Context, education *Education) error
	}

	EducationUsecase interface {
		FindByID(ctx context.Context, id int64) (*Education, error)
	}

	// CreateEducationInput create education input
	CreateEducationInput struct {
		Name             string  `json:"name" validate:"required,min=3,max=60"`
		Email            string  `json:"email" validate:"required,emailEligibility"`
		PhoneNumber      *string `json:"phone_number" validate:"omitempty,phonenumber"`
		Address          string  `json:"address" validate:"required,min=1"`
		LogoImageID      int64   `json:"logo_image_id" validate:"omitempty,required"`
		IdentityColorHex string  `json:"identity_color_hex" validate:"required,hexcolor"`
	}
)
