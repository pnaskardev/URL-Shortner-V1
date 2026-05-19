package routes

import (
	"github.com/gofiber/fiber/v3"
	requesthelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/requestHelper"
	authrepository "github.com/pnaskardev/URL-Shortner-V1/core/pkg/handlers/auth"
	"gorm.io/gorm"
)

func ApiRouter(app *fiber.App, requestClient *requesthelper.RetryableHTTPClient, dbClient *gorm.DB) {

	router := app.Group("/api")

	authRouter := router.Group("/auth")

	authHandler := authrepository.New(*requestClient, dbClient)
	authRouter.Post("/sign-in", authHandler.SignInHandler)
	authRouter.Post("/sign-up", authHandler.SignUpHandler)
}
