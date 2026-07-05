package redirector

import (
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	RedirectURL(c fiber.Ctx) error
}

type repository struct {
	dbClient    *gorm.DB
	cacheClient *redis.Client
	// If we have DB client in here
	// All of the routes will get the DB client and no need to make multiple connections
}

func New(dbClient *gorm.DB, cacheClient *redis.Client) Repository {
	return &repository{
		dbClient:    dbClient,
		cacheClient: cacheClient,
	}
}

func (r repository) RedirectURL(c fiber.Ctx) error {
	return nil
}
