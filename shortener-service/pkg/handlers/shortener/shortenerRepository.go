package shortener

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/pnaskardev/URL-Shortner-V1/core/middlewares"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/helpers/constants"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/helpers/utils"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/helpers/views"
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

	shortenPayload := new(views.ShortenRequest)

	userID := c.RequestCtx().UserValue(middlewares.UserIDKey).(string)

	if err := c.Bind().Body(shortenPayload); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	shortenedURL := utils.EncodeToBase62([]byte(shortenPayload.Url))

	shortenedQueueBody := views.ShortenerQueueBody{
		ID:              uuid.New(),
		UserID:          userID,
		ShortenedURLKey: shortenedURL,
		LongURL:         shortenPayload.Url,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	bodyBytes, err := json.Marshal(shortenedQueueBody)
	if err != nil {
		return err
	}

	err = r.queueClient.Publish(
		ctx,
		constants.URL_CREATED_QUEUE,
		bodyBytes,
	)

	if err != nil {

		slog.Error("ERROR", "PUBLISH ERROR", err)
		panic(err)
	}

	return c.SendStatus(fiber.StatusOK)
}
