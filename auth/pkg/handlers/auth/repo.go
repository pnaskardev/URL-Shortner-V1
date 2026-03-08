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

func (h *AuthHandler) SignIn(ctx context.Context, req *authPB.HelloRequest) (*authPB.HelloReply, error) {
	// validate, check password, generate token
	fmt.Println(req.Name)
	return &authPB.HelloReply{Message: fmt.Sprintf("HELLO, %s", req.Name)}, nil
}

func (h *AuthHandler) SignUp(ctx context.Context, req *authPB.HelloRequest) (*authPB.HelloReply, error) {
	// validate, create user, generate token
	fmt.Println("SIGN UP INVOKED")
	return &authPB.HelloReply{}, nil
}
