package authrepository

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgconn"
	requesthelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/requestHelper"
	responsehelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/responseHelper"
	"github.com/pnaskardev/URL-Shortner-V1/core/helpers/utils"
	corevalidator "github.com/pnaskardev/URL-Shortner-V1/core/helpers/validator"
	"github.com/pnaskardev/URL-Shortner-V1/core/helpers/views"
	"github.com/pnaskardev/URL-Shortner-V1/core/infrastructure/database"
	"github.com/pnaskardev/URL-Shortner-V1/core/infrastructure/database/models"
)

type Repository interface {
	SignInHandler(c fiber.Ctx) error
	SignUpHandler(c fiber.Ctx) error
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

func (r *repository) SignUpHandler(c fiber.Ctx) error {
	authPayload := new(views.AuthSignInPayload)

	if err := c.Bind().Body(authPayload); err != nil {
		return responsehelper.BadRequest(c)
	}

	if errs, err := corevalidator.ValidateStruct(authPayload); err != nil {
		return responsehelper.ValidationError(c, errs)
	}

	// Hash password
	password, salt, err := utils.HashPassword(authPayload.Password)
	if err != nil {
		return responsehelper.InternalServerError(c)
	}

	hashStr := base64.StdEncoding.EncodeToString(password)
	saltStr := base64.StdEncoding.EncodeToString(salt)

	user := models.User{
		Username: authPayload.Username,
		Password: hashStr,
		Salt:     saltStr,
	}

	dbClient := database.ConnectToPostgres()
	err = dbClient.Create(&user).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {

			if pgErr.Code == "23505" {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"message": "user already exists please use a different username",
				})
			}

		}
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "user has been created please try logging in",
	})

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

	request, err := requesthelper.CreateHTTPRequest(ctx, "POST", "http://localhost:8001/api/auth/sign-in", &authPayload)
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
