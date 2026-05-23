package shortner

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/pnaskardev/URL-Shortner-V1/core/httpClients"
	"gorm.io/gorm"
)

type Repository interface {
	ShortenURL(c fiber.Ctx) error
}

type repository struct {
	requestClient httpClients.RetryableHTTPClient
	dbClient      *gorm.DB
	// If we have DB client in here
	// All of the routes will get the DB client and no need to make multiple connections
}

func New(requestClient httpClients.RetryableHTTPClient, dbClient *gorm.DB) Repository {
	return &repository{
		requestClient: requestClient,
		dbClient:      dbClient,
	}
}

func (r *repository) ShortenURL(c fiber.Ctx) error {

	log.Debug("SHORTEN URL BODY", string(c.Body()))
	return nil
}
