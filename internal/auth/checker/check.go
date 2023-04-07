package checker

import (
	"context"

	"github.com/gotd/td/telegram/auth"
	"github.com/pkg/errors"
	"tgavatar/internal/telegram"
)

//go:generate mockgen -source=check.go -destination=./mock_check/check.go -package=mock_check

type Client interface {
	Auth() *auth.Client
	Run(ctx context.Context, f func(ctx context.Context) error) (err error)
}

// CheckerAuth checks auth from telegram *auth.Client
type CheckerAuth interface {
	CheckAuth(ctx context.Context, client Client) (bool, error)
	GetCheckerFunc(client Client) func(ctx context.Context) error
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
func (s checkerAuth) CheckAuth(ctx context.Context, client Client) (bool, error) {
	authorized, err := s.checkerAuthStatus.CheckAuth(ctx, client.Auth())
	if err != nil {
		return false, errors.Wrap(err, "failed to check auth from telegram auth")
	}

	return authorized, nil
}

func (s checkerAuth) GetCheckerFunc(client Client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		res, err := s.CheckAuth(ctx, client)
		if err != nil {
			return errors.Wrap(err, "failed to check auth")
		}

		if !res {
			return NoAuthorizedErr
		}

		return nil
	}
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

var NoAuthorizedErr = errors.New("user is not authorized")

type telegramFactory interface {
	GetClient() (*telegram.TGClient, error)
}

type TgFactory struct {
	f telegramFactory
}

func (t TgFactory) GetClient() (Client, error) {
	return t.f.GetClient()
}

type TgFactoryInterface interface {
	GetClient() (Client, error)
}

func NewTgFactory(f telegramFactory) TgFactoryInterface {
	return TgFactory{f: f}
}

// Checker authorization checker
type Checker interface {
	CheckAuth(ctx context.Context) (bool, error)
}

// checker struct for authorization checking
type checker struct {
	telegramFactory TgFactoryInterface
	checkerAuth     CheckerAuth
}

// NewChecker Checker authorization checker constructor
func NewChecker(telegramFactory TgFactoryInterface, checkerAuth CheckerAuth) Checker {
	return checker{
		telegramFactory: telegramFactory,
		checkerAuth:     checkerAuth,
	}
}

// CheckAuth checks telegram authorization for current session
func (c checker) CheckAuth(ctx context.Context) (bool, error) {
	client, err := c.telegramFactory.GetClient()
	if err != nil {
		return false, errors.Wrap(err, "failed to create TGClient for check auth")
	}

	if err := client.Run(ctx, c.checkerAuth.GetCheckerFunc(client)); err != nil {
		if errors.Is(err, NoAuthorizedErr) {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to start TGClient for check auth")
	}

	return true, nil
}
