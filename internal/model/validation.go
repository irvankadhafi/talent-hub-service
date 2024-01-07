package model

import (
	"github.com/go-playground/validator/v10"
	"github.com/irvankadhafi/talent-hub-service/internal/config"
	"regexp"
	"strings"
	"sync"
)

// validate singleton, it's thread safe and cached the struct validation rules
var validate *validator.Validate

// singleton regex
var phoneNumberRgx *regexp.Regexp

var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		validate = validator.New()

		_ = validate.RegisterValidation("phonenumber", isPhoneValid)

		_ = validate.RegisterValidation("emailEligibility", isEmailValid)

		_ = validate.RegisterValidation("identifier", validateIdentifier)

		phoneNumberRgx = regexp.MustCompile(`(^(\+)|^[0-9]+$)`)
	})
}

// isPhoneValid implements validator.Func for check phone number
func isPhoneValid(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	if len(phone) < 3 {
		return false
	}

	return phoneNumberRgx.MatchString(phone)
}

// isEmailValid implements validator.Func for check email
func isEmailValid(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	// General validator for email
	if err := validate.Var(email, "email"); err != nil {
		return false
	}

	splitedEmail := strings.Split(email, "@")
	specialChars := "!%#^-"

	// the use of e-mail with the + symbol is not allowed on production
	if config.Env() == config.EnvProduction {
		specialChars += "+"
	}

	return !strings.ContainsAny(splitedEmail[0], specialChars)
}

func validateIdentifier(fl validator.FieldLevel) bool {
	return isEmailValid(fl) || isPhoneValid(fl)
}
