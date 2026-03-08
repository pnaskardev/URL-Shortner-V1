package rpc

import (
	"fmt"

	"google.golang.org/grpc"
)

type RPCServer struct {
	grpc    *grpc.Server
	address string
}

func NewRPCServer(address string) *RPCServer {
	return &RPCServer{
		grpc:    grpc.NewServer(),
		address: address,
	}
}

func (s *RPCServer) Register(fn func(*grpc.Server)) {
	fn(s.grpc)
}

func (s *RPCServer) GetServer() *grpc.Server {
	return s.grpc
}

func (s *RPCServer) GetAddr() string {
	return s.address
}

func (s *RPCServer) Stop() {
	fmt.Println("GRACEFULLY STOPPING THE FRPC SERVER")
	if s.grpc != nil {
		s.grpc.GracefulStop()
	}
}
