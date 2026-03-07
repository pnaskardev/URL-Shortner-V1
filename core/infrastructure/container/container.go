package container

import (
	"flag"

	"github.com/pnaskardev/URL-Shortner-V1/core/infrastructure/rpc"
	"github.com/pnaskardev/URL-Shortner-V1/core/rpc-services/user"

	authPB "github.com/pnaskardev/URL-Shortner-V1/url-shortner-rpc/auth"
)

var addr = flag.String("addr", "localhost:3000", "RPC server address")

type Container struct {
	RPC         *rpc.RPCClient
	AuthService *user.AuthService
}

var instance *Container

func Init() *Container {
	flag.Parse()

	rpcClient := &rpc.RPCClient{}
	rpcClient.Connect(*addr)

	instance = &Container{
		RPC:         rpcClient,
		AuthService: user.New(authPB.NewAuthClient(rpcClient.Conn())),
	}
	return instance
}

func Get() *Container {
	return instance
}
