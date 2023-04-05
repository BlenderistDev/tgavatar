package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

type log interface {
	Error(args ...interface{})
	Info(args ...interface{})
}

// noSignUp can be embedded to prevent signing up.
type noSignUp struct{}

// SignUp to implement noSignUp
func (c noSignUp) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

// AcceptTermsOfService to implement noSignUp
func (c noSignUp) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// termAuth implements authentication via channel.
type termAuth struct {
	noSignUp

	phone string

	dataChan chan string
}

// Phone return phone for auth
func (a termAuth) Phone(_ context.Context) (string, error) {
	return a.phone, nil
}

// Password return password for auth. Fake implemented
// @todo add 2FA password support
func (a termAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

// Code get authorization code from channel
func (a termAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {

	code := <-a.dataChan

	return strings.TrimSpace(code), nil
}

type telegramFactory interface {
	GetClient() (*telegram.Client, error)
}

type Authorizer interface {
	Auth(phone string, codeChan chan string) error
}

// Auth struct for telegram authorization
type authorizer struct {
	successAuthChan chan struct{}
	ctx             context.Context
	log             log
	telegramFactory telegramFactory
}

// NewAuth constructor for Auth
func NewAuth(ctx context.Context, log log, telegramFactory telegramFactory, successAuthChan chan struct{}) Authorizer {
	return authorizer{
		successAuthChan: successAuthChan,
		ctx:             ctx,
		log:             log,
		telegramFactory: telegramFactory,
	}
}

// Auth init telegram authorization
func (a authorizer) Auth(phone string, codeChan chan string) error {
	flow := auth.NewFlow(
		termAuth{phone: strings.Clone(phone), dataChan: codeChan},
		auth.SendCodeOptions{},
	)

	client, err := a.telegramFactory.GetClient()
	if err != nil {
		return errors.Wrap(err, "failed to create telegram client for auth flow init")
	}

	go func() {
		err = client.Run(a.ctx, func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
				return errors.Wrap(err, "failed to start telegram auth flow")
			}
			a.successAuthChan <- struct{}{}

			return nil
		})
		if err != nil {
			// @todo add error handling
			a.log.Error(errors.Wrap(err, "auth error"))
		}
	}()

	return nil
}
