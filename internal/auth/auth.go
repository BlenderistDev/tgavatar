package auth

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

// noSignUp can be embedded to prevent signing up.
type noSignUp struct{}

func (c noSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c noSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// termAuth implements authentication via terminal.
type termAuth struct {
	noSignUp

	phone string

	dataChan chan string
}

func (a termAuth) Phone(_ context.Context) (string, error) {
	return a.phone, nil
}

func (a termAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

func (a termAuth) Code(_ context.Context, s *tg.AuthSentCode) (string, error) {

	code := <-a.dataChan

	return strings.TrimSpace(code), nil
}

type Auth struct {
	successAuthChan chan struct{}
}

func NewAuth(successAuthChan chan struct{}) Auth {
	return Auth{
		successAuthChan: successAuthChan,
	}
}

func (a Auth) Auth(phone string, codeChan chan string) error {
	flow := auth.NewFlow(
		termAuth{phone: strings.Clone(phone), dataChan: codeChan},
		auth.SendCodeOptions{},
	)

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		SessionStorage: &session.FileStorage{Path: "storage/session"},
	})

	if err != nil {
		return err
	}

	go func() {
		err = client.Run(context.Background(), func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
				return err
			}
			a.successAuthChan <- struct{}{}

			return nil
		})
		if err != nil {
			// @todo add error handling
			log.Println(errors.Wrap(err, "auth error"))
		}
	}()

	return nil
}
