package auth

import (
	"context"

	"github.com/gotd/td/telegram/auth"
	"github.com/pkg/errors"
	"tgavatar/internal/auth/checker/auth/status"
)

type client interface {
	Auth() *auth.Client
}

// CheckerAuth checks auth from telegram *auth.Client
type CheckerAuth interface {
	CheckAuth(ctx context.Context, client client) (bool, error)
}

type checkerAuth struct {
	checkerAuthStatus status.CheckerAuthStatus
}

// NewCheckerAuth constructor for CheckerAuth
func NewCheckerAuth(checkerAuthStatus status.CheckerAuthStatus) CheckerAuth {
	return checkerAuth{
		checkerAuthStatus: checkerAuthStatus,
	}
}

// CheckAuth checks auth from telegram *auth.Client
func (s checkerAuth) CheckAuth(ctx context.Context, client client) (bool, error) {
	authorized, err := s.checkerAuthStatus.CheckAuth(ctx, client.Auth())
	if err != nil {
		return false, errors.Wrap(err, "failed to check auth from telegram auth")
	}

	return authorized, nil
}
