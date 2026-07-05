package routes

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func ApiRouter(app *fiber.App, dbClient *gorm.DB) {

	// router := app.Group("/api")

	// THIS API SHOULD BE RATE LIMITED AND AUTHENTICATED
	// shortenHandler := shortener.New(dbClient, queueClient)
	// router.Post("/shorten", authFiberHandler, shortenHandler.ShortenURL)
}
