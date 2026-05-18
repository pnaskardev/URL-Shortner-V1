package authhandler

import (
	"github.com/gofiber/fiber/v3"
	responsehelper "github.com/pnaskardev/URL-Shortner-V1/auth/helpers/responseHelper"
	"github.com/pnaskardev/URL-Shortner-V1/auth/helpers/views"
)

type Repository interface {
	SignUpHandler(c fiber.Ctx) error
}

type repository struct {

	// If we have DB client in here
	// All of the routes will get the DB client and no need to make multiple connections
}

func New() Repository {
	return &repository{}
}

func (r *repository) SignUpHandler(c fiber.Ctx) error {

	authPayload := new(views.AuthSignInPayload)
	if err := c.Bind().Body(authPayload); err != nil {
		return responsehelper.BadRequest(c)
	}

	return nil
}
