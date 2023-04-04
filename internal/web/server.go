package web

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"tgavatar/internal/auth"
)

type log interface {
	Error(args ...interface{})
	Info(args ...interface{})
}

// Handler struct for web routes handlers
type Handler struct {
	codeChan    chan string
	authChecker auth.Checker
	auth        auth.Auth
	log         log
}

// LaunchAuthServer start web server for telegram auth
func LaunchAuthServer(authChecker auth.Checker, auth auth.Auth, log log) error {
	app := fiber.New()

	codeChan := make(chan string)
	h := Handler{
		codeChan:    codeChan,
		authChecker: authChecker,
		auth:        auth,
		log:         log,
	}

	app.Get("/", h.Auth)
	app.Post("/phone", h.Phone)
	app.Post("/code", h.Code)

	err := app.Listen(os.Getenv("HOST"))
	if err != nil {
		return errors.Wrap(err, "failed to get web server host")
	}

	return nil
}

// Auth handler for GET / route
func (h Handler) Auth(c *fiber.Ctx) error {
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

// Phone handler for POST /phone route
func (h Handler) Phone(c *fiber.Ctx) error {
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

	err = h.auth.Auth(req.Phone, h.codeChan)
	if err != nil {
		h.log.Error(errors.Wrap(err, "auth error for /phone"))
		return fiber.ErrInternalServerError
	}

	return c.SendFile("html/code.html")
}

// Code handler for POST /code route
func (h Handler) Code(c *fiber.Ctx) error {
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
