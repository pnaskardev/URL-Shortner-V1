package authrepository

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	requesthelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/requestHelper"
	responsehelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/responseHelper"
	corevalidator "github.com/pnaskardev/URL-Shortner-V1/core/helpers/validator"
	"github.com/pnaskardev/URL-Shortner-V1/core/helpers/views"
)

type Repository interface {
	SignInHandler(c fiber.Ctx) error
}

type repository struct {
	requestClient requesthelper.RetryableHTTPClient

	// If we have DB client in here
	// All of the routes will get the DB client and no need to make multiple connections
}

func New(requestClient requesthelper.RetryableHTTPClient) Repository {
	return &repository{
		requestClient: requestClient,
	}
}

func (r *repository) SignInHandler(c fiber.Ctx) error {
	authPayload := new(views.AuthSignInPayload)

	if err := c.Bind().Body(authPayload); err != nil {
		return responsehelper.BadRequest(c)
	}

	// Validate the request
	if errs, err := corevalidator.ValidateStruct(authPayload); err != nil {
		return responsehelper.ValidationError(c, errs)
	}

	slog.Info("AUTH PAYLOAD", "AUTHPAYLOAD", authPayload)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	request, err := requesthelper.CreateHTTPRequest(ctx, "POST", "http://localhost:8001/sign-in", &authPayload)
	if err != nil {
		slog.Error("SIGN IN HANDLER ERROR", "REQUEST CREATION FAILED", err)
		return responsehelper.InternalServerError(c)
	}

	response, err := r.requestClient.DoWithRetry(ctx, request)
	if err != nil {
		slog.Error("SIGN IN HANDLER ERROR", "AUTH REQUEST FAILED", err)
		return responsehelper.InternalServerError(c)
	}

	slog.Debug("SIGN IN HANDLER", "AUTH RESPONSE", response)

	return c.Status(200).JSON(response)
}
