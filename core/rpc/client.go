package rpc_service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RPCClient struct {
	Conn *grpc.ClientConn
}

func RPCNewClientConnection(url string) (*RPCClient, error) {

	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &RPCClient{
		Conn: conn,
	}, nil
}
