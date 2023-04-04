package auth

import (
	"context"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/pkg/errors"
)

var noAuthorizedErr = errors.New("user is not authorized")

// Checker authorization checker
type Checker interface {
	CheckAuth(ctx context.Context) (bool, error)
}

// checker struct for authorization checking
type checker struct {
	storagePath string
}

// NewChecker Checker authorization checker constructor
func NewChecker(storagePath string) Checker {
	return checker{
		storagePath: storagePath,
	}
}

// CheckAuth checks telegram authorization for current session
func (c checker) CheckAuth(ctx context.Context) (bool, error) {
	client, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{Path: c.storagePath},
	})

	if err != nil {
		return false, errors.Wrap(err, "failed to create client for check auth")
	}

	if err := client.Run(ctx, func(ctx context.Context) error {
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get auth status for check auth")
		}

		if !status.Authorized {
			return noAuthorizedErr
		}

		return nil
	}); err != nil {
		if errors.Is(err, noAuthorizedErr) {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to start client for check auth")
	}

	return true, nil
}
