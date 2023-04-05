package web

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=server.go -destination=./mock/server.go -package=mock_server

type log interface {
	Error(args ...interface{})
	Info(args ...interface{})
}

type authChecker interface {
	CheckAuth(ctx context.Context) (bool, error)
}

type authorizer interface {
	Auth(phone string, codeChan chan string) error
}

// handler struct for web routes handlers
type handler struct {
	codeChan    chan string
	authChecker authChecker
	authorizer  authorizer
	log         log
}

// LaunchAuthServer start web server for telegram auth
func LaunchAuthServer(authChecker authChecker, authorizer authorizer, log log) (*fiber.App, error) {
	app := fiber.New()

	codeChan := make(chan string)
	h := handler{
		codeChan:    codeChan,
		authChecker: authChecker,
		authorizer:  authorizer,
		log:         log,
	}

	app.Get("/", h.auth)
	app.Post("/phone", h.phone)
	app.Post("/code", h.code)

	err := app.Listen(os.Getenv("HOST"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get web server host")
	}

	return app, nil
}

// auth handler for GET / route
func (h handler) auth(c *fiber.Ctx) error {
	authorized, err := h.authChecker.CheckAuth(c.Context())
	if err != nil {
		h.log.Error(errors.Wrap(err, "auth check failed for /auth"))
		return fiber.ErrInternalServerError
	}

	if authorized {
		return c.SendFile("html/authorized.html")
	}

	return c.SendFile("html/auth.html")
}

// phone handler for POST /phone route
func (h handler) phone(c *fiber.Ctx) error {
	authorized, err := h.authChecker.CheckAuth(c.Context())
	if err != nil {
		h.log.Error(errors.Wrap(err, "auth check failed for /phone"))
		return fiber.ErrInternalServerError
	}

	if authorized {
		return fiber.NewError(fiber.StatusPreconditionFailed, "already authorized")
	}

	req := struct {
		Phone string `form:"phone"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		h.log.Error(errors.Wrap(err, "failed to decode body for /phone"))
		return fiber.ErrBadRequest
	}

	err = h.authorizer.Auth(req.Phone, h.codeChan)
	if err != nil {
		h.log.Error(errors.Wrap(err, "auth error for /phone"))
		return fiber.ErrInternalServerError
	}

	return c.SendFile("html/code.html")
}

// code handler for POST /code route
func (h handler) code(c *fiber.Ctx) error {
	authorized, err := h.authChecker.CheckAuth(c.Context())
	if err != nil {
		h.log.Error(errors.Wrap(err, "auth check failed for /code"))
		return fiber.ErrInternalServerError
	}

	if authorized {
		return fiber.NewError(fiber.StatusPreconditionFailed, "already authorized")
	}

	req := struct {
		Code string `form:"code"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		h.log.Error(errors.Wrap(err, "failed to decode body for /code"))
		return fiber.ErrBadRequest
	}

	h.codeChan <- req.Code

	return c.SendFile("html/success.html")
}
