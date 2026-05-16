package authhandler

import (
	"context"

	authPB "github.com/pnaskardev/URL-Shortner-V1/url-shortner-rpc/auth"
)

type AuthHandler struct {
	authPB.UnimplementedAuthServer
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) SignIn(ctx context.Context) error {
	return nil
}

func (h *AuthHandler) SignUp(ctx context.Context) error {
	return nil
}
