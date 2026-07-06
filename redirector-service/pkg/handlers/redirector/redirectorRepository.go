package redirector

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/pnaskardev/URL-Shortner-V1/redirector-service/pkg/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const cacheTTL = 1 * time.Hour

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

	shortenedKey := c.Params("key")

	if shortenedKey == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	longURL, err := r.cacheClient.Get(ctx, shortenedKey).Result()
	switch {
	case err == nil:
		// Cache hit — redirect straight away.
		return c.Redirect().Status(fiber.StatusFound).To(longURL)

	case errors.Is(err, redis.Nil):
		// Cache miss — fall through to Postgres.

	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var fetchedURL models.URLs
	dbErr := r.dbClient.WithContext(ctx).Where("short_url_key = ?", shortenedKey).Take(&fetchedURL).Error

	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if dbErr != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if !fetchedURL.IsActive || fetchedURL.ExpiresAt < time.Now().UnixMilli() {
		return c.SendStatus(fiber.StatusGone)
	}

	// Warm the cache so the next read is a hit.
	r.cacheClient.Set(ctx, shortenedKey, fetchedURL.LongURL, cacheTTL)

	return c.Redirect().Status(fiber.StatusFound).To(fetchedURL.LongURL)
}
