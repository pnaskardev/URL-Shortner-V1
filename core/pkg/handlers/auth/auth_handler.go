package auth_handler

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/pnaskardev/URL-Shortner-V1/core/infrastructure/container"
)

type Repository interface {
	SignInHandler(c fiber.Ctx) error
}

type repository struct{}

func New() Repository {
	return &repository{}
}

func (r *repository) SignInHandler(c fiber.Ctx) error {
	rpcContainer := container.Get()
	reply, err := rpcContainer.Services.AuthRPCService.SignIn(c.Context())
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
	return c.SendStatus(200)
}
