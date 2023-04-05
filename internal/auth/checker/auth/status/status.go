package status

import (
	"context"

	tgAuth "github.com/gotd/td/telegram/auth"
	"github.com/pkg/errors"
)

type auth interface {
	Status(ctx context.Context) (*tgAuth.Status, error)
}

// CheckerAuthStatus checks auth from telegram *auth.Status
type CheckerAuthStatus interface {
	CheckAuth(ctx context.Context, auth auth) (bool, error)
}

type checkerAuthStatus struct {
}

// NewCheckerStatusAuth constructor for CheckerAuthStatus
func NewCheckerStatusAuth() CheckerAuthStatus {
	return checkerAuthStatus{}
}

// CheckAuth checks auth from telegram *auth.Status
func (s checkerAuthStatus) CheckAuth(ctx context.Context, auth auth) (bool, error) {
	status, err := auth.Status(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get auth status for check auth")
	}

	return status.Authorized, nil
}
