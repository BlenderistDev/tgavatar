package auth

import (
	"context"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/pkg/errors"
)

// Checker struct for authorization checking
type Checker struct {
}

// CheckAuth checks telegram authorization for current session
func (c Checker) CheckAuth(ctx context.Context) (bool, error) {
	client, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{Path: "storage/session"},
	})

	if err != nil {
		return false, errors.Wrap(err, "failed to create client for check auth")
	}

	var authorized bool

	if err := client.Run(ctx, func(ctx context.Context) error {
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get auth status for check auth")
		}

		authorized = status.Authorized

		return nil
	}); err != nil {
		return false, errors.Wrap(err, "failed to start client for check auth")
	}

	return authorized, nil
}
