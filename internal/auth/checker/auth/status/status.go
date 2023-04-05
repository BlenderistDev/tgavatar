package status

import (
	"context"

	"github.com/gotd/td/telegram/auth"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=status.go -mock_names=tgAuth=tgAuth -destination=./mock_status/status.go -package=mock_status

type tgAuth interface {
	Status(ctx context.Context) (*auth.Status, error)
}

// CheckerAuthStatus checks auth from telegram *auth.Status
type CheckerAuthStatus interface {
	CheckAuth(ctx context.Context, auth tgAuth) (bool, error)
}

type checkerAuthStatus struct {
}

// NewCheckerStatusAuth constructor for CheckerAuthStatus
func NewCheckerStatusAuth() CheckerAuthStatus {
	return checkerAuthStatus{}
}

// CheckAuth checks auth from telegram *auth.Status
func (s checkerAuthStatus) CheckAuth(ctx context.Context, a tgAuth) (bool, error) {
	status, err := a.Status(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get auth status for check auth")
	}

	return status.Authorized, nil
}
