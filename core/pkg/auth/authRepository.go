package authrepository

import (
	"github.com/gofiber/fiber/v3"
)

type Repository interface {
	SignInHandler(c fiber.Ctx) error
}

type repository struct {
	// If we have DB client in here
	// All of the routes will get the DB client and no need to make multiple connections
}

func New() Repository {
	return &repository{}
}

func (r *repository) SignInHandler(c fiber.Ctx) error {

	return c.SendStatus(200)
}
