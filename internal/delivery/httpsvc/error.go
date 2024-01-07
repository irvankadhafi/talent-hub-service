package httpsvc

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/irvankadhafi/talent-hub-service/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// http errors
var (
	ErrInvalidArgument            = echo.NewHTTPError(http.StatusBadRequest, "invalid argument")
	ErrNotFound                   = echo.NewHTTPError(http.StatusNotFound, "record not found")
	ErrInternal                   = echo.NewHTTPError(http.StatusInternalServerError, "internal system error")
	ErrEntityTooLarge             = echo.NewHTTPError(http.StatusRequestEntityTooLarge, "entity too large")
	ErrUnauthenticated            = echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	ErrUnauthorized               = echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	ErrEmailTokenNotMatch         = echo.NewHTTPError(http.StatusUnauthorized, "email or token not match")
	ErrEmailPasswordNotMatch      = echo.NewHTTPError(http.StatusUnauthorized, "email or password not match")
	ErrPermissionDenied           = echo.NewHTTPError(http.StatusForbidden, "permission denied")
	ErrSpaceNotSelected           = echo.NewHTTPError(http.StatusBadRequest, "space not selected")
	ErrLoginByEmailPasswordLocked = echo.NewHTTPError(http.StatusLocked, "user is locked from logging in using email and password")
	ErrInvitationExpired          = echo.NewHTTPError(http.StatusBadRequest, "invitation expired")
	ErrFailedPrecondition         = echo.NewHTTPError(http.StatusPreconditionFailed, "precondition failed")
)

// httpValidationOrInternalErr return valdiation or internal error
func httpValidationOrInternalErr(err error) error {
	switch t := err.(type) {
	case validator.ValidationErrors:
		_ = t
		errVal := err.(validator.ValidationErrors)

		fields := map[string]interface{}{}
		for _, ve := range errVal {
			fields[ve.Field()] = fmt.Sprintf("Failed on the '%s' tag", ve.Tag())
		}

		return echo.NewHTTPError(http.StatusBadRequest, utils.Dump(fields))
	default:
		return ErrInternal
	}
}
