package authhandler

import (
	"context"
	"fmt"

	authPB "github.com/pnaskardev/URL-Shortner-V1/url-shortner-rpc/auth"
)

type AuthHandler struct {
	authPB.UnimplementedAuthServer
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) SignIn(ctx context.Context, req *authPB.SignInRequest) (*authPB.SignInResponse, error) {
	// validate, check password, generate token
	fmt.Println(req.Username)
	return &authPB.SignInResponse{Message: fmt.Sprintf("HELLO, %s", req.Username)}, nil
}

func (h *AuthHandler) SignUp(ctx context.Context, req *authPB.SignUpRequest) (*authPB.SignUpResponse, error) {
	// validate, create user, generate token
	fmt.Println("SIGN UP INVOKED")
	return &authPB.SignUpResponse{}, nil
}
