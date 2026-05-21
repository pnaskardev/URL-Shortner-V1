package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pnaskardev/URL-Shortner-V1/core/httpClients"
	authrepository "github.com/pnaskardev/URL-Shortner-V1/core/pkg/handlers/auth"
	"gorm.io/gorm"
)

func ApiRouter(app *fiber.App, requestClient *httpClients.RetryableHTTPClient, dbClient *gorm.DB) {

	router := app.Group("/api")

	authRouter := router.Group("/auth")

	authHandler := authrepository.New(*requestClient, dbClient)
	authRouter.Post("/sign-in", authHandler.SignInHandler)
	authRouter.Post("/sign-up", authHandler.SignUpHandler)

	// shortenRouter := router.Group("/shorten")

}
