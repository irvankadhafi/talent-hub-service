package delivery

import (
	"context"
	"github.com/irvankadhafi/talent-hub-service/auth"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
)

// GetAuthCandidateFromCtx ..
func GetAuthCandidateFromCtx(ctx context.Context) *model.Candidate {
	authCandidate := auth.GetCandidateFromCtx(ctx)
	if authCandidate == nil {
		return nil
	}

	user := &model.Candidate{
		ID:        authCandidate.ID,
		SessionID: authCandidate.SessionID,
	}

	return user
}
