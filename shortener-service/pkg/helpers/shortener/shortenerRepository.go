package shortener

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/infrastructure/queue"
	"gorm.io/gorm"
)

type Repository interface {
	ShortenURL(c fiber.Ctx) error
}

type repository struct {
	dbClient    *gorm.DB
	queueClient *queue.QueueClient
	// If we have DB client in here
	// All of the routes will get the DB client and no need to make multiple connections
}

func New(dbClient *gorm.DB, queueClient *queue.QueueClient) Repository {
	return &repository{
		dbClient:    dbClient,
		queueClient: queueClient,
	}
}

func (r *repository) ShortenURL(c fiber.Ctx) error {

	slog.Debug("DEBUG", "SHORTEN URL BODY", string(c.Body()))

	return c.SendStatus(fiber.StatusOK)
}
