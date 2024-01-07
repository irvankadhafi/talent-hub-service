package model

import (
	"context"
	"github.com/irvankadhafi/talent-hub-service/internal/helper"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

type (
	CandidateUsecase interface {
		Create(ctx context.Context, input CreateCandidateInput) (*Candidate, error)
		FindByID(ctx context.Context, id int64) (*Candidate, error)
	}

	CandidateRepository interface {
		FindByID(ctx context.Context, id int64) (*Candidate, error)
		FindPasswordByID(ctx context.Context, id int64) ([]byte, error)
		FindByEmail(ctx context.Context, email string) (*Candidate, error)
		FindUnscopedByEmail(ctx context.Context, email string) (*Candidate, error)
		FindUnscopedByPhone(ctx context.Context, phone string) (*Candidate, error)
		FindByPhone(ctx context.Context, phone string) (*Candidate, error)
		Create(ctx context.Context, candidate *Candidate) error
		Update(ctx context.Context, candidate *Candidate) error
	}

	Candidate struct {
		ID             int64          `json:"id"`
		FullName       string         `json:"full_name"`
		Email          null.String    `json:"email"`
		Phone          null.String    `json:"phone"`
		Password       string         `json:"password"`
		DateOfBirth    null.Time      `json:"date_of_birth"`
		Gender         Gender         `json:"gender"`
		CityID         int64          `json:"city_id"`
		ProvinceID     int64          `json:"province_id"`
		LastEducation  time.Time      `json:"last_education"`
		LastExperience time.Time      `json:"last_experience"`
		LoginDate      time.Time      `json:"login_date"`
		CreatedAt      time.Time      `json:"created_at" gorm:"->;<-:create"`
		UpdatedAt      time.Time      `json:"updated_at"`
		DeletedAt      gorm.DeletedAt `json:"deleted_at"`

		SessionID int64  `json:"-" gorm:"-"`
		Latitude  string `json:"latitude" gorm:"-"`
		Longitude string `json:"longitude" gorm:"-"`
	}
)

// Gender the candidate's gender
type Gender string

// Gender constants
const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

// CreateCandidateInput :nodoc:
type CreateCandidateInput struct {
	FullName             string `json:"full_name" validate:"required"`
	Email                string `json:"email" validate:"omitempty,emailEligibility"`
	Phone                string `json:"phone" validate:"omitempty,phonenumber"`
	Gender               Gender `json:"gender" validate:"required"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,eqfield=Password"`
}

// ValidateAndFormat do field validation and format the PhoneNumber
// Flow: Remove Leading 0 -> Validate -> Formating With Country Code
func (c *CreateCandidateInput) ValidateAndFormat() error {
	_ = helper.RemoveLeadingZeroPhoneNumber(&c.Phone)

	if err := validate.Struct(c); err != nil {
		return err
	}

	if c.Phone != "" {
		if err := helper.FormatPhoneNumberWithCountryCode(&c.Phone, "ID"); err != nil {
			return err
		}
	}

	return nil
}
