package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/pnaskardev/URL-Shortner-V1/core/httpClients"
	"github.com/pnaskardev/URL-Shortner-V1/core/middlewares"
	authrepository "github.com/pnaskardev/URL-Shortner-V1/core/pkg/handlers/auth"
	"github.com/pnaskardev/URL-Shortner-V1/core/pkg/handlers/shortner"
	"gorm.io/gorm"
)

func ApiRouter(app *fiber.App, requestClient *httpClients.RetryableHTTPClient, dbClient *gorm.DB) {

	router := app.Group("/api")

	authRouter := router.Group("/auth")

	// THESE APIs should be OPEN FOR ALL but RATE LIMITED
	authHandler := authrepository.New(*requestClient, dbClient)
	authRouter.Post("/sign-in", authHandler.SignInHandler)
	authRouter.Post("/sign-up", authHandler.SignUpHandler)

	authFiberHandler := adaptor.HTTPMiddleware(middlewares.AuthMiddleware)

	// THIS API SHOULD BE RATE LIMITED AND AUTHENTICATED
	shortenHandler := shortner.New(*requestClient, dbClient)
	router.Post("/shorten", authFiberHandler, shortenHandler.ShortenURL)
}
