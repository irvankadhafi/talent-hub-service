package auth

import (
	"context"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
)

type contextKey string

// use module path to make it unique
const candidateCtxKey contextKey = "github.com/irvankadhafi/talent-hub-service/auth.Candidate"

// SetUserToCtx set user to context
func SetUserToCtx(ctx context.Context, candidate Candidate) context.Context {
	return context.WithValue(ctx, candidateCtxKey, candidate)
}

// GetCandidateFromCtx get user from context
func GetCandidateFromCtx(ctx context.Context) *Candidate {
	user, ok := ctx.Value(candidateCtxKey).(Candidate)
	if !ok {
		return nil
	}
	return &user
}

// Candidate represent an authenticated candidate
type Candidate struct {
	ID        int64 `json:"id"`
	SessionID int64 `json:"session_id"`
}

// NewCandidateFromSession return new candidate from session
func NewCandidateFromSession(sess model.Session) Candidate {
	return Candidate{
		ID:        sess.CandidateID,
		SessionID: sess.ID,
	}
}
