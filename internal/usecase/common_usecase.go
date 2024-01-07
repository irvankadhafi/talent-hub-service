package usecase

import (
	"context"
	"errors"
	"github.com/irvankadhafi/talent-hub-service/internal/config"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"github.com/irvankadhafi/talent-hub-service/utils"
	"github.com/mattheath/base62"
	"strings"
	"time"
)

// GenerateToken and check uniqueness
func GenerateToken(sr model.SessionRepository, candidateID int64) (token string, err error) {
	sleep := 10 * time.Millisecond
	ctxTimeout := 50 * time.Millisecond
	candidateIDEnc := base62.EncodeInt64(candidateID)
	err = utils.Retry(3, sleep, func() error {
		sb := strings.Builder{}
		sb.WriteString(candidateIDEnc)
		sb.WriteString("_")

		randomAlphanum := utils.GenerateRandomAlphanumeric(config.DefaultSessionTokenLength)
		sb.WriteString(randomAlphanum)
		token = sb.String()

		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		exist, err := sr.CheckToken(ctx, token)
		if err != nil {
			return err
		}
		if exist {
			return errors.New("token exists, retry")
		}

		return nil
	})

	return token, err
}
