package routes

import (
	"github.com/gofiber/fiber/v3"
	authrepository "github.com/pnaskardev/URL-Shortner-V1/core/pkg/auth"
)

func ApiRouter(app *fiber.App) {

	router := app.Group("/api")

	authRouter := router.Group("/auth")
	authRouter.Post("/sign-in", authrepository.New().SignInHandler)
}
