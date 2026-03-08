package container

import (
	"flag"

	customRPC "github.com/pnaskardev/URL-Shortner-V1/auth/infrastructure/rpc"
	authhandler "github.com/pnaskardev/URL-Shortner-V1/auth/pkg/handlers/auth"
	"github.com/pnaskardev/URL-Shortner-V1/auth/pkg/handlers/health"
	authPB "github.com/pnaskardev/URL-Shortner-V1/url-shortner-rpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var addr = flag.String("addr", "localhost:3000", "RPC server address")

type Container struct {
	Server      *customRPC.RPCServer
	AuthHandler authPB.AuthServer
}

func Init() *Container {
	flag.Parse()

	healthHandler := health.NewHealthCheckHandler()
	authHandler := authhandler.NewAuthHandler()

	srv := customRPC.NewRPCServer(*addr)

	// Register Health Server First
	srv.Register(func(s *grpc.Server) {
		grpc_health_v1.RegisterHealthServer(s, healthHandler)
	})

	srv.Register(func(s *grpc.Server) {
		authPB.RegisterAuthServer(s, authHandler)
	})

	return &Container{
		Server:      srv,
		AuthHandler: authHandler,
	}
}
