package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/infrastructure/queue"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/middlewares"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/pkg/handlers/shortener"
	"gorm.io/gorm"
)

func ApiRouter(app *fiber.App, dbClient *gorm.DB, queueClient *queue.QueueClient) {

	router := app.Group("/api")

	authFiberHandler := adaptor.HTTPMiddleware(middlewares.MicroServiceAuthMiddleware)

	// THIS API SHOULD BE RATE LIMITED AND AUTHENTICATED
	shortenHandler := shortener.New(dbClient, queueClient)
	router.Post("/shorten", authFiberHandler, shortenHandler.ShortenURL)
}
