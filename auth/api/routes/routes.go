package routes

import (
	"github.com/gofiber/fiber/v3"
	authhandler "github.com/pnaskardev/URL-Shortner-V1/auth/pkg/handlers/auth"
	"github.com/pnaskardev/URL-Shortner-V1/auth/pkg/handlers/health"
)

func ApiRouter(app *fiber.App) {

	router := app.Group("/api")

	healthCheckRepoInstance := health.New()
	router.Get("/health", healthCheckRepoInstance.HealthCheck)

	authRepositoryInstance := authhandler.New()
	authRouter := router.Group("/auth")
	authRouter.Post("/sign-up", authRepositoryInstance.SignUpHandler)
}
