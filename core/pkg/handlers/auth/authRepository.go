package authrepository

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgconn"
	requesthelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/requestHelper"
	responsehelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/responseHelper"
	"github.com/pnaskardev/URL-Shortner-V1/core/helpers/utils"
	corevalidator "github.com/pnaskardev/URL-Shortner-V1/core/helpers/validator"
	"github.com/pnaskardev/URL-Shortner-V1/core/helpers/views"
	"github.com/pnaskardev/URL-Shortner-V1/core/infrastructure/database/models"
	"gorm.io/gorm"
)

type Repository interface {
	SignInHandler(c fiber.Ctx) error
	SignUpHandler(c fiber.Ctx) error
}

type repository struct {
	requestClient requesthelper.RetryableHTTPClient
	dbClient      *gorm.DB
	// If we have DB client in here
	// All of the routes will get the DB client and no need to make multiple connections
}

func New(requestClient requesthelper.RetryableHTTPClient, dbClient *gorm.DB) Repository {
	return &repository{
		requestClient: requestClient,
		dbClient:      dbClient,
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

	err = r.dbClient.Create(&user).Error
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
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

	var user models.User

	err := r.dbClient.Model(&models.User{}).Where("username = ? AND deleted = ?", authPayload.Username, false).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responsehelper.NotFound(c, "user not found")
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// real DB error (constraint, syntax, etc.)
			fmt.Println(pgErr.Code, pgErr.Message)
		}

		return responsehelper.InternalServerError(c)
	}

	saltBytes, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return responsehelper.InternalServerError(c)
	}

	hashBytes, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return responsehelper.InternalServerError(c)
	}

	verificationResult := utils.VerifyPassword(authPayload.Password, saltBytes, hashBytes)

	if !verificationResult {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Wrong Password/Username"})
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		slog.Error("ERROR", "ERROR WHILE GENERATING ACCESS TOKEN", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	authResponse := fiber.Map{
		"message":      "Authentication Succesful",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	return c.Status(fiber.StatusAccepted).JSON(authResponse)
}
