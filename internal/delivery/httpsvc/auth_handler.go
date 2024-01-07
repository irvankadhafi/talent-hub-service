package httpsvc

import (
	"github.com/irvankadhafi/talent-hub-service/internal/delivery"
	"github.com/irvankadhafi/talent-hub-service/internal/delivery/httpsvc/dto"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"github.com/irvankadhafi/talent-hub-service/internal/usecase"
	"github.com/irvankadhafi/talent-hub-service/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Service) handleLoginByIdentifierPassword() echo.HandlerFunc {
	type loginRequest struct {
		Identifier string `json:"identifier"` // Can be either email or phone
		Password   string `json:"password"`
	}

	return func(c echo.Context) error {
		req := loginRequest{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		loginReq := model.LoginRequest{
			Identifier:    req.Identifier,
			PlainPassword: req.Password,
			IPAddress:     c.RealIP(),
			UserAgent:     c.Request().UserAgent(),
			// TODO: implement longitude and latitude
		}

		session, err := s.authUsecase.LoginByIdentifierPassword(c.Request().Context(), loginReq)
		switch err {
		case nil:
			break
		case usecase.ErrNotFound, usecase.ErrUnauthorized:
			return ErrEmailPasswordNotMatch
		case usecase.ErrLoginByEmailPasswordLocked:
			return ErrLoginByEmailPasswordLocked
		case usecase.ErrPermissionDenied:
			return ErrPermissionDenied
		default:
			logrus.Error(err)
			return ErrInternal
		}

		res := dto.LoginResponse{
			AccessToken:           session.AccessToken,
			AccessTokenExpiresAt:  utils.FormatTimeRFC3339(&session.AccessTokenExpiredAt),
			RefreshToken:          session.RefreshToken,
			RefreshTokenExpiresAt: utils.FormatTimeRFC3339(&session.RefreshTokenExpiredAt),
			TokenType:             "Bearer",
		}

		return c.JSON(http.StatusOK, dto.NewSuccessResponse(res, "Success Login"))
	}
}

func (s *Service) handleRefreshToken() echo.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	return func(c echo.Context) error {
		req := request{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		session, err := s.authUsecase.RefreshToken(c.Request().Context(), model.RefreshTokenRequest{
			RefreshToken: req.RefreshToken,
			IPAddress:    c.RealIP(),
			UserAgent:    c.Request().UserAgent(),
			// TODO: implement longitude and latitude

		})
		switch err {
		case nil:
		case usecase.ErrRefreshTokenExpired, usecase.ErrNotFound:
			return ErrUnauthenticated
		default:
			logrus.Error(err)
			return ErrInternal
		}

		res := dto.LoginResponse{
			AccessToken:           session.AccessToken,
			AccessTokenExpiresAt:  utils.FormatTimeRFC3339(&session.AccessTokenExpiredAt),
			RefreshToken:          session.RefreshToken,
			RefreshTokenExpiresAt: utils.FormatTimeRFC3339(&session.RefreshTokenExpiredAt),
			TokenType:             "Bearer",
		}

		return c.JSON(http.StatusOK, dto.NewSuccessResponse(res, "Success Refresh Token"))
	}
}

func (s *Service) handleLogout() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		requester := delivery.GetAuthCandidateFromCtx(ctx)

		err := s.authUsecase.DeleteSessionByID(c.Request().Context(), requester.SessionID)
		switch err {
		case nil:
			break
		case usecase.ErrNotFound:
			return ErrNotFound
		default:
			logrus.Error(err)
			return httpValidationOrInternalErr(err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
