package routes

import (
	"github.com/gofiber/fiber/v3"
	requesthelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/requestHelper"
	authrepository "github.com/pnaskardev/URL-Shortner-V1/core/pkg/auth"
)

func ApiRouter(app *fiber.App, requestClient *requesthelper.RetryableHTTPClient) {

	router := app.Group("/api")

	authRouter := router.Group("/auth")

	authHandler := authrepository.New(*requestClient)
	authRouter.Post("/sign-in", authHandler.SignInHandler)
}
