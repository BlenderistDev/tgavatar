package checker

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/pkg/errors"
	"tgavatar/internal/auth/checker/auth"
)

var noAuthorizedErr = errors.New("user is not authorized")

type telegramFactory interface {
	GetClient() (*telegram.Client, error)
}

// Checker authorization checker
type Checker interface {
	CheckAuth(ctx context.Context) (bool, error)
}

// checker struct for authorization checking
type checker struct {
	telegramFactory telegramFactory
	checkerAuth     auth.CheckerAuth
}

// NewChecker Checker authorization checker constructor
func NewChecker(telegramFactory telegramFactory, checkerAuth auth.CheckerAuth) Checker {
	return checker{
		telegramFactory: telegramFactory,
		checkerAuth:     checkerAuth,
	}
}

// CheckAuth checks telegram authorization for current session
func (c checker) CheckAuth(ctx context.Context) (bool, error) {
	client, err := c.telegramFactory.GetClient()
	if err != nil {
		return false, errors.Wrap(err, "failed to create client for check auth")
	}

	if err := client.Run(ctx, func(ctx context.Context) error {
		client.Auth()
		res, err := c.checkerAuth.CheckAuth(ctx, client)
		if err != nil {
			return errors.Wrap(err, "failed to check auth")
		}

		if !res {
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
