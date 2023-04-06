package check

import (
	"context"

	"github.com/gotd/td/telegram/auth"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=check.go -mock_names=client=client -destination=./mock_check/check.go -package=mock_check

type client interface {
	Auth() *auth.Client
}

// CheckerAuth checks auth from telegram *auth.Client
type CheckerAuth interface {
	CheckAuth(ctx context.Context, client client) (bool, error)
}

type checkerAuth struct {
	checkerAuthStatus CheckerAuthStatusInterface
}

// NewCheckerAuth constructor for CheckerAuth
func NewCheckerAuth(checkerAuthStatus CheckerAuthStatusInterface) CheckerAuth {
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

type TgAuthInterface interface {
	Status(ctx context.Context) (*auth.Status, error)
}

// CheckerAuthStatusInterface checks auth from telegram *auth.Status
type CheckerAuthStatusInterface interface {
	CheckAuth(ctx context.Context, auth TgAuthInterface) (bool, error)
}

type checkerAuthStatus struct {
}

// NewCheckerStatusAuth constructor for CheckerAuthStatusInterface
func NewCheckerStatusAuth() CheckerAuthStatusInterface {
	return checkerAuthStatus{}
}

// CheckAuth checks auth from telegram *auth.Status
func (s checkerAuthStatus) CheckAuth(ctx context.Context, a TgAuthInterface) (bool, error) {
	status, err := a.Status(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get auth status for check auth")
	}

	return status.Authorized, nil
}
