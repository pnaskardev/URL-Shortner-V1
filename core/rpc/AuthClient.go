package rpc_service

import (
	"github.com/pnaskardev/URL-Shortner-V1/core/proto/auth"
	"google.golang.org/grpc"
)

type AuthRPCClient struct {
	Client auth.AuthServiceClient
	Conn   *grpc.ClientConn
}

func NewAuthRPCClientConnection(url string) (*AuthRPCClient, error) {
	rpcConnectionClient, err := RPCNewClientConnection(url)
	if err != nil {
		return nil, err
	}

	return &AuthRPCClient{
		Client: auth.NewAuthServiceClient(rpcConnectionClient.Conn),
		Conn:   rpcConnectionClient.Conn,
	}, nil

}
