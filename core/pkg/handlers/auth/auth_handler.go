package auth_handler

import (
	"github.com/gofiber/fiber/v3"
)

type Repository interface {
	SignInHandler(c fiber.Ctx) error
}

type repository struct{}

func New() Repository {
	return &repository{}
}

func (r *repository) SignInHandler(c fiber.Ctx) error {

	return c.SendStatus(200)
}
