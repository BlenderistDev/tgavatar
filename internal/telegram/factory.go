package telegram

import (
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/pkg/errors"
)

// storagePath path to store telegram session json file
const storagePath = "storage/session"

// Factory telegram client factory
type Factory interface {
	GetClient() (*telegram.Client, error)
}

type factory struct {
}

// NewFactory telegram client factory constructor
func NewFactory() Factory {
	return factory{}
}

// GetClient build new telegram client
func (f factory) GetClient() (*telegram.Client, error) {
	client, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{Path: storagePath},
	})

	if err != nil {
		return nil, errors.Wrap(err, "telegram client creating error")
	}

	return client, nil
}
