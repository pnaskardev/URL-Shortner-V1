package rpc

import (
	"google.golang.org/grpc"
)

type Server struct {
	grpc    *grpc.Server
	address string
}

func NewServer(address string) *Server {
	return &Server{
		grpc:    grpc.NewServer(),
		address: address,
	}
}

// Register allows each service to mount itself onto the server
// Usage: server.Register(func(s *grpc.Server) { authPB.RegisterAuthServer(s, &myImpl{}) })
func (s *Server) Register(fn func(*grpc.Server)) {
	fn(s.grpc)
}

func (s *Server) GetServer() *grpc.Server {
	return s.grpc
}

func (s *Server) GetAddr() string {
	return s.address
}

func (s *Server) Stop() {
	if s.grpc != nil {
		s.grpc.GracefulStop()
	}
}
