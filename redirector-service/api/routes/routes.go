package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pnaskardev/URL-Shortner-V1/redirector-service/infrastructure/cache"
	"github.com/pnaskardev/URL-Shortner-V1/redirector-service/pkg/handlers/redirector"
	"gorm.io/gorm"
)

func ApiRouter(app *fiber.App, dbClient *gorm.DB) {

	// router := app.Group("/api")

	cacheClient := cache.GetRedisClient()

	// THIS API SHOULD BE RATE LIMITED AND AUTHENTICATED
	redirectHandler := redirector.New(dbClient, cacheClient)
	app.Get("/:key", redirectHandler.RedirectURL)
}
