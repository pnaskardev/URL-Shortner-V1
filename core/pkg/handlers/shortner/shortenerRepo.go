package shortner

import (
	"io"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/protobuf/encoding/protojson"

	urlshortenerv1 "github.com/pnaskardev/URL-Shortner-V1/contracts/gen/go/urlshortener/v1"
	responsehelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/responseHelper"
	"github.com/pnaskardev/URL-Shortner-V1/core/helpers/utils"
	corevalidator "github.com/pnaskardev/URL-Shortner-V1/core/helpers/validator"
	"github.com/pnaskardev/URL-Shortner-V1/core/helpers/views"
	"github.com/pnaskardev/URL-Shortner-V1/core/httpClients"
	"github.com/pnaskardev/URL-Shortner-V1/core/middlewares"
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

	shortenPayload := new(views.ShortenRequest)

	userID := c.RequestCtx().UserValue(middlewares.UserIDKey).(string)

	if err := c.Bind().Body(shortenPayload); err != nil {
		return responsehelper.BadRequest(c)
	}

	if errs, err := corevalidator.ValidateStruct(shortenPayload); err != nil {
		return responsehelper.ValidationError(c, errs)
	}

	internalRequest := &urlshortenerv1.InternalShortenRequest{
		Url:    shortenPayload.Url,
		UserId: userID,
	}

	internalRequestBody, err := protojson.Marshal(internalRequest)
	if err != nil {
		return responsehelper.InternalServerError(c)
	}

	doWithRetryPayload := httpClients.DoWithRetryRequest{
		Ctx:     c,
		Method:  "POST",
		Service: utils.SHORTENER_SERVICE,
		Payload: internalRequestBody,
	}

	response, err := r.requestClient.DoWithRetry(doWithRetryPayload)
	if err != nil {
		return responsehelper.InternalServerError(c)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return responsehelper.InternalServerError(c)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(response.StatusCode).Send(body)

}
