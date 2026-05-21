package responsehelper

import (
	"github.com/gofiber/fiber/v3"
	corevalidator "github.com/pnaskardev/URL-Shortner-V1/core/helpers/validator"
)

func BadRequest(c fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "INVALID REQUEST",
	})
}

func ValidationError(c fiber.Ctx, errors []corevalidator.FieldError) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"errors": errors,
	})
}

func InternalServerError(c fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Something went wrong",
	})
}

func NotFound(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": message,
	})
}
