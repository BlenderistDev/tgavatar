package telegram

import (
	"context"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
)

// storagePath path to store telegram session json file
const storagePath = "storage/session"

type TGClient struct {
	client *telegram.Client
}

func (c TGClient) Auth() *auth.Client {
	return c.client.Auth()
}

func (c TGClient) Run(ctx context.Context, f func(ctx context.Context) error) (err error) {
	return c.client.Run(ctx, f)
}

func (c TGClient) API() *tg.Client {
	return c.client.API()
}

//// Factory telegram TGClient factory
//type Factory interface {
//	GetClient() (*TGClient, error)
//}

type Factory struct {
}

// NewFactory telegram TGClient factory constructor
func NewFactory() Factory {
	return Factory{}
}

// GetClient build new telegram TGClient
func (f Factory) GetClient() (*TGClient, error) {
	c, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{Path: storagePath},
	})

	if err != nil {
		return nil, errors.Wrap(err, "telegram TGClient creating error")
	}

	return &TGClient{client: c}, nil
}
