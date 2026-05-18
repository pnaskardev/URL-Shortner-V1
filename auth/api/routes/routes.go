package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pnaskardev/URL-Shortner-V1/auth/pkg/handlers/health"
)

func ApiRouter(app *fiber.App) {

	router := app.Group("/api")

	healthCheckRepoInstance := health.New()
	router.Get("/health", healthCheckRepoInstance.HealthCheck)

}
