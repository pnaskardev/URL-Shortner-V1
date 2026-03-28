package routes

import (
	"github.com/gofiber/fiber/v3"
	auth_handler "github.com/pnaskardev/URL-Shortner-V1/core/pkg/handlers/auth"
)

func ApiRouter(app *fiber.App) {

	router := app.Group("/api")

	authRouter := router.Group("/auth")
	authRouter.Post("/sign-in", auth_handler.New().SignInHandler)

}
