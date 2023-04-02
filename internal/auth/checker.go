package auth

import (
	"context"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
)

type Checker struct {
}

func (c Checker) CheckAuth(ctx context.Context) (bool, error) {
	client, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{Path: "storage/session"},
	})

	if err != nil {
		return false, err
	}

	var authorized bool

	if err := client.Run(ctx, func(ctx context.Context) error {
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return err
		}

		authorized = status.Authorized

		return nil
	}); err != nil {
		return false, err
	}

	return authorized, nil
}
