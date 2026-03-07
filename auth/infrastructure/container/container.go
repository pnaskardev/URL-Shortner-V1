package container

import (
	"flag"

	customRPC "github.com/pnaskardev/URL-Shortner-V1/auth/infrastructure/rpc"
	authhandler "github.com/pnaskardev/URL-Shortner-V1/auth/pkg/handlers/auth"
	authPB "github.com/pnaskardev/URL-Shortner-V1/url-shortner-rpc/auth"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", "localhost:3000", "RPC server address")

type Container struct {
	Server      *customRPC.Server
	AuthHandler authPB.AuthServer
}

func Init() *Container {
	flag.Parse()

	authHandler := authhandler.NewAuthHandler()

	srv := customRPC.NewServer(*addr)
	srv.Register(func(s *grpc.Server) {
		authPB.RegisterAuthServer(s, authHandler)
	})

	return &Container{
		Server:      srv,
		AuthHandler: authHandler,
	}
}
