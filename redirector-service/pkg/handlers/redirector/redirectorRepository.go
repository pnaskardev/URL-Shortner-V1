package redirector

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	urlshortenerv1 "github.com/pnaskardev/URL-Shortner-V1/contracts/gen/go/urlshortener/v1"
	constants "github.com/pnaskardev/URL-Shortner-V1/redirector-service/helpers"
	"github.com/pnaskardev/URL-Shortner-V1/redirector-service/helpers/utils"
	"github.com/pnaskardev/URL-Shortner-V1/redirector-service/infrastructure/queue"
	"github.com/pnaskardev/URL-Shortner-V1/redirector-service/pkg/models"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

const cacheTTL = 1 * time.Hour

type Repository interface {
	RedirectURL(c fiber.Ctx) error
}

type repository struct {
	dbClient    *gorm.DB
	cacheClient *redis.Client
	queueClient *queue.QueueClient
	// If we have DB client in here
	// All of the routes will get the DB client and no need to make multiple connections
}

func New(dbClient *gorm.DB, queueClient *queue.QueueClient, cacheClient *redis.Client) Repository {
	return &repository{
		dbClient:    dbClient,
		cacheClient: cacheClient,
		queueClient: queueClient,
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

	// Capture click context BEFORE the goroutine — fiber.Ctx is pooled and
	// reused for the next request once this handler returns.
	occurredAtMs := time.Now().UnixMilli()
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")
	referrer := c.Get("Referer")

	// Push the click to the analytics pipeline. Fire-and-forget: lossy by design,
	// must not add latency or failure to the redirect.
	go func() {

		publishCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		snowFlakeID, err := utils.NewSnowflakeID()
		if err != nil {
			slog.Error("ID GENERATION ERROR ANALYTICS EVENT", "ERROR", err)
			return
		}
		analyticEvent := urlshortenerv1.UrlRedirectedEvent{
			Id:              snowFlakeID,
			ShortenedUrlKey: shortenedKey,
			OccurredAtMs:    occurredAtMs,
			IpAddress:       ipAddress,
			UserAgent:       userAgent,
			Referrer:        referrer,
		}

		body, err := proto.Marshal(&analyticEvent)
		if err != nil {
			slog.Error("MARSHAL ERROR ANALYTICS EVENT", "ERROR", err)
			return
		}

		if err := r.queueClient.Publish(publishCtx, constants.URL_ANALYTICS_EVENT_QUEUE, body); err != nil {
			slog.Error("PUBLISH ANALYTICS EVENT", "KEY", shortenedKey, "ERROR", err)
		}

	}()

	return c.Redirect().Status(fiber.StatusFound).To(fetchedURL.LongURL)
}
