package user

import (
	"context"
	"fmt"

	authPB "github.com/pnaskardev/URL-Shortner-V1/url-shortner-rpc/auth"
)

type AuthService struct {
	authClient authPB.AuthClient
}

func New(authClient authPB.AuthClient) *AuthService {
	return &AuthService{
		authClient: authClient,
	}
}

func (s *AuthService) SignIn(ctx context.Context) (*authPB.HelloReply, error) {
	fmt.Println(&authPB.HelloRequest{Name: "PRIYANSHU"})
	return s.authClient.SignIn(ctx, &authPB.HelloRequest{Name: "PRIYANSHU"})
}
