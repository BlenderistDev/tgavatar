package check

import (
	"context"

	"github.com/gotd/td/telegram/auth"
	"github.com/pkg/errors"
	"tgavatar/internal/telegram"
)

//go:generate mockgen -source=check.go -destination=./mock_check.go -package=check

var NoAuthorizedErr = errors.New("user is not authorized")

type client interface {
	Auth() *auth.Client
	Run(ctx context.Context, f func(ctx context.Context) error) (err error)
}

// authChecker checks auth from telegram *auth.Client
type authChecker interface {
	CheckAuth(ctx context.Context, client client) (bool, error)
}

type authCheck struct {
	statusChecker statusChecker
}

// CheckAuth checks auth from telegram *auth.Client
func (s authCheck) CheckAuth(ctx context.Context, client client) (bool, error) {
	authorized, err := s.statusChecker.CheckAuth(ctx, client.Auth())
	if err != nil {
		return false, errors.Wrap(err, "failed to check auth from telegram auth")
	}

	return authorized, nil
}

func (c check) getCheckerFunc(client client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		res, err := c.authChecker.CheckAuth(ctx, client)
		if err != nil {
			return errors.Wrap(err, "failed to check auth")
		}

		if !res {
			return NoAuthorizedErr
		}

		return nil
	}
}

type tgAuthInterface interface {
	Status(ctx context.Context) (*auth.Status, error)
}

// statusChecker checks auth from telegram *auth.Status
type statusChecker interface {
	CheckAuth(ctx context.Context, auth tgAuthInterface) (bool, error)
}

type statusCheck struct {
}

// CheckAuth checks auth from telegram *auth.Status
func (s statusCheck) CheckAuth(ctx context.Context, a tgAuthInterface) (bool, error) {
	status, err := a.Status(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get auth status for check auth")
	}

	return status.Authorized, nil
}

type telegramFactory interface {
	GetClient() (*telegram.TGClient, error)
}

type tgFactory struct {
	factory telegramFactory
}

func (t tgFactory) GetClient() (client, error) {
	return t.factory.GetClient()
}

type tgFactoryInterface interface {
	GetClient() (client, error)
}

// Checker authorization check
type Checker interface {
	CheckAuth(ctx context.Context) (bool, error)
}

// check struct for authorization checking
type check struct {
	telegramFactory tgFactoryInterface
	authChecker     authChecker
}

// NewChecker Checker authorization check constructor
func NewChecker(telegramFactory telegramFactory) Checker {
	return check{
		telegramFactory: tgFactory{factory: telegramFactory},
		authChecker:     authCheck{statusChecker: statusCheck{}},
	}
}

// CheckAuth checks telegram authorization for current session
func (c check) CheckAuth(ctx context.Context) (bool, error) {
	client, err := c.telegramFactory.GetClient()
	if err != nil {
		return false, errors.Wrap(err, "failed to create TGClient for check auth")
	}

	if err := client.Run(ctx, c.getCheckerFunc(client)); err != nil {
		if errors.Is(err, NoAuthorizedErr) {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to start TGClient for check auth")
	}

	return true, nil
}
