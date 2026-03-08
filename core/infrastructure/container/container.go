package container

import (
	"context"
	"time"

	"github.com/pnaskardev/URL-Shortner-V1/core/infrastructure/rpc"
	"github.com/pnaskardev/URL-Shortner-V1/core/rpc-services/user"
	authPB "github.com/pnaskardev/URL-Shortner-V1/url-shortner-rpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type RPCCLients struct {
	AuthRPCClient *rpc.RPCClient
	// TODO ADD MORE
}

type RPCServices struct {
	AuthRPCService *user.AuthService
	// TODO ADD MORE
}

type Container struct {
	Clients  *RPCCLients
	Services *RPCServices
}

var instance *Container

func InitContainer() *Container {

	address := "localhost:3000"

	authClient, err := rpc.NewRPCClient(address)
	if err != nil {
		panic(err)
	}

	// DO a quick health check
	waitForReady(authClient.GetConn())

	clients := &RPCCLients{
		AuthRPCClient: authClient,
	}

	services := &RPCServices{
		AuthRPCService: user.New(authPB.NewAuthClient(authClient.GetConn())),
	}

	instance = &Container{
		Clients:  clients,
		Services: services,
	}
	return instance
}

func waitForReady(conn *grpc.ClientConn) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	_, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})

	if err != nil {
		panic(err)
	}

	return nil

}

func Get() *Container {
	return instance
}
