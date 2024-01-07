package httpsvc

import (
	"github.com/irvankadhafi/talent-hub-service/auth"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"github.com/labstack/echo/v4"
)

// Service http service
type Service struct {
	group          *echo.Group
	authUsecase    model.AuthUsecase
	authMiddleware *auth.AuthenticationMiddleware
}

// RouteService add dependencies and use group for routing
func RouteService(
	group *echo.Group,
	authUSecase model.AuthUsecase,
	authMiddleware *auth.AuthenticationMiddleware,
) {
	srv := &Service{
		group:          group,
		authUsecase:    authUSecase,
		authMiddleware: authMiddleware,
	}
	srv.initRoutes()
}

func (s *Service) initRoutes() {
	s.group.POST("/auth/login/", s.handleLoginByIdentifierPassword())
	s.group.POST("/auth/tokens/refresh/", s.handleRefreshToken())
	s.group.POST("/auth/logout/", s.handleLogout(), s.authMiddleware.MustAuthenticateAccessToken())
}
