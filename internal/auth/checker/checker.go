package checker

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/pkg/errors"
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
	checkerAuth     CheckerAuth
}

// NewChecker Checker authorization checker constructor
func NewChecker(telegramFactory telegramFactory, checkerAuth CheckerAuth) Checker {
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

	if err := client.Run(ctx, c.checkerAuth.GetCheckerFunc(client)); err != nil {
		if errors.Is(err, noAuthorizedErr) {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to start client for check auth")
	}

	return true, nil
}
