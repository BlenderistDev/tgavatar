package web

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"tgavatar/internal/auth"
)

type Handler struct {
	codeChan    chan string
	authChecker auth.Checker
	auth        auth.Auth
}

func LaunchAuthServer(authChecker auth.Checker, auth auth.Auth) error {
	app := fiber.New()

	codeChan := make(chan string)
	h := Handler{
		codeChan:    codeChan,
		authChecker: authChecker,
		auth:        auth,
	}

	app.Get("/", h.Auth)
	app.Post("/phone", h.Phone)
	app.Post("/code", h.Code)

	err := app.Listen(os.Getenv("HOST"))
	if err != nil {
		return err
	}

	return nil
}

func (h Handler) Auth(c *fiber.Ctx) error {
	authorized, err := h.authChecker.CheckAuth(c.Context())
	if err != nil {
		log.Println(errors.Wrap(err, "auth check failed"))
		return fiber.ErrInternalServerError
	}

	if authorized {
		return c.SendFile("html/authorized.html")
	}

	return c.SendFile("html/auth.html")
}

func (h Handler) Phone(c *fiber.Ctx) error {
	authorized, err := h.authChecker.CheckAuth(c.Context())
	if err != nil {
		log.Println(errors.Wrap(err, "auth check failed"))
		return fiber.ErrInternalServerError
	}

	if authorized {
		return fiber.NewError(fiber.StatusPreconditionFailed, "already authorized")
	}

	req := struct {
		Phone string `form:"phone"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	err = h.auth.Auth(req.Phone, h.codeChan)
	if err != nil {
		log.Println(errors.Wrap(err, "auth error"))
		return fiber.ErrInternalServerError
	}

	return c.SendFile("html/code.html")
}

func (h Handler) Code(c *fiber.Ctx) error {
	authorized, err := h.authChecker.CheckAuth(c.Context())
	if err != nil {
		log.Println(errors.Wrap(err, "auth check failed"))
		return fiber.ErrInternalServerError
	}

	if authorized {
		return fiber.NewError(fiber.StatusPreconditionFailed, "already authorized")
	}

	req := struct {
		Code string `form:"code"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	h.codeChan <- req.Code

	return c.SendFile("html/success.html")
}
