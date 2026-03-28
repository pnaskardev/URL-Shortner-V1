package user

import (
	"context"

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

func (s *AuthService) SignIn(ctx context.Context) (*authPB.SignInResponse, error) {
	return s.authClient.SignIn(ctx, &authPB.SignInRequest{Username: "priyanshu", Password: "testpass"})
}
