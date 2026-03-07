package rpc

import (
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RPCClient struct {
	conn *grpc.ClientConn
	once sync.Once
}

func NewRPCClient(addr string) *RPCClient {
	return &RPCClient{}
}

var defaultAddr = "localhost:3000"

func (c *RPCClient) Connect(addr string) {
	c.once.Do(func() {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to connect to RPC server: %v", err)
		}

		c.conn = conn
	})

}

func (c *RPCClient) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *RPCClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
