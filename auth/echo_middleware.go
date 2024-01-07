package auth

import (
	"context"
	"encoding/json"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"github.com/irvankadhafi/talent-hub-service/pkg/cacher"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

const (
	_authScheme          = "Bearer"
	_headerAuthorization = "Authorization"
)

// CandidateAuthenticator to perform candidate authentication
type CandidateAuthenticator interface {
	AuthenticateToken(ctx context.Context, accessToken string) (*Candidate, error)
}

// AuthenticationMiddleware middleware for authentication
type AuthenticationMiddleware struct {
	cacheManager    cacher.CacheManager
	candidateAuther CandidateAuthenticator
}

// NewAuthenticationMiddleware AuthMiddleware constructor
func NewAuthenticationMiddleware(
	candidateAuther CandidateAuthenticator,
	cacheManager cacher.CacheManager,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		candidateAuther: candidateAuther,
		cacheManager:    cacheManager,
	}
}

// AuthenticateAccessToken authenticate access token from http `Authorization` header and load a Candidate to context
func (a *AuthenticationMiddleware) AuthenticateAccessToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getAccessToken(c.Request())
			return a.authenticateAccessToken(c, next, token)
		}
	}
}

// MustAuthenticateAccessToken must authenticate access token from http `Authorization` header and load a Candidate to context
// Differ from AuthenticateAccessToken, if no token provided then return Unauthenticated
func (a *AuthenticationMiddleware) MustAuthenticateAccessToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getAccessToken(c.Request())
			if token == "" {
				return errorResp(http.StatusUnauthorized, "user is unauthenticated")
			}

			return a.authenticateAccessToken(c, next, token)
		}
	}
}

func (a *AuthenticationMiddleware) authenticateAccessToken(c echo.Context, next echo.HandlerFunc, token string) error {
	// only load user to context when token presented
	if token == "" {
		return next(c)
	}
	ctx := c.Request().Context()

	session, err := a.findSessionFromCache(token)
	switch err {
	default:
		// cache error will fallback to rpc
		logrus.WithField("sessionCacheError", "find session from cache got error").Error(err)
	case nil:
		if session == nil {
			break // fallback
		}
		if session.IsAccessTokenExpired() {
			return errorResp(http.StatusUnauthorized, "token expired")
		}

		ctx := SetUserToCtx(c.Request().Context(), NewCandidateFromSession(*session))
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}

	userSession, err := a.candidateAuther.AuthenticateToken(ctx, token)
	// fallback to rpc
	switch status.Code(err) {
	case codes.OK:
		if userSession == nil { // safety check
			return next(c)
		}

		ctx := SetUserToCtx(c.Request().Context(), *userSession)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	case codes.NotFound:
		return errorResp(http.StatusBadRequest, "token is invalid")
	case codes.Unauthenticated:
		return errorResp(http.StatusUnauthorized, "token is expired")
	default:
		logrus.Error(err)
		return errorResp(http.StatusInternalServerError, "system error")
	}
}

func (a *AuthenticationMiddleware) findSessionFromCache(token string) (*model.Session, error) {
	reply, err := a.cacheManager.Get(model.NewSessionTokenCacheKey(token))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if reply == nil {
		return nil, nil
	}

	bt, _ := reply.([]byte)
	if bt == nil {
		return nil, nil
	}

	sess := &model.Session{}
	err = json.Unmarshal(bt, sess)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return sess, nil
}

func getAccessToken(req *http.Request) (accessToken string) {
	authHeaders := strings.Split(req.Header.Get(_headerAuthorization), " ")

	if (len(authHeaders) != 2) || (authHeaders[0] != _authScheme) {
		return ""
	}

	return strings.TrimSpace(authHeaders[1])
}

func errorResp(code int, message string) error {
	return echo.NewHTTPError(code, echo.Map{"message": message})
}
