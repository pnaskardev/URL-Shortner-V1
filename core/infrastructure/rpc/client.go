package rpc

import (
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RPCClient struct {
	conn *grpc.ClientConn
	once sync.Once
}

func NewRPCClient(addr string) (*RPCClient, error) {

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 3 * time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	return &RPCClient{conn: conn}, nil
}

func (c *RPCClient) GetConn() *grpc.ClientConn {
	return c.conn
}

func (c *RPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
