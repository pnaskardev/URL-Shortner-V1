package shortener

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	urlshortenerv1 "github.com/pnaskardev/URL-Shortner-V1/contracts/gen/go/urlshortener/v1"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/helpers/constants"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/helpers/utils"
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

	shortenPayload := new(urlshortenerv1.InternalShortenRequest)

	if err := protojson.Unmarshal(c.Body(), shortenPayload); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// UserId is carried in the InternalShortenRequest contract, set by core after
	// it authenticates the caller. Reading it from the body keeps this service from
	// depending on core's internal packages.
	userID := shortenPayload.GetUserId()
	if userID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	id, err := utils.NewSnowflakeID()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	shortenedURLKey := utils.EncodeToBase62(new(big.Int).SetInt64(id).Bytes())

	event := &urlshortenerv1.UrlCreatedEvent{
		Id:              uuid.New().String(),
		UserId:          userID,
		ShortenedUrlKey: shortenedURLKey,
		LongUrl:         shortenPayload.GetUrl(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	bodyBytes, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	err = r.queueClient.Publish(
		ctx,
		constants.URL_CREATED_QUEUE,
		bodyBytes,
	)

	if err != nil {
		slog.Error("failed to publish url.created event", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response := &urlshortenerv1.ShortenResponse{
		Id:              event.GetId(),
		ShortenedUrlKey: shortenedURLKey,
		LongUrl:         shortenPayload.GetUrl(),
	}

	responseBody, err := protojson.Marshal(response)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusAccepted).Send(responseBody)
}
