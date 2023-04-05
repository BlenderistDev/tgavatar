package telegram

import (
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/pkg/errors"
)

// storagePath path to store telegram session json file
const storagePath = "storage/session"

type Factory interface {
	GetClient() (*telegram.Client, error)
}

type factory struct {
}

func NewFactory() Factory {
	return factory{}
}

func (f factory) GetClient() (*telegram.Client, error) {
	client, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{Path: storagePath},
	})

	if err != nil {
		return nil, errors.Wrap(err, "telegram client creating error")
	}

	return client, nil
}
