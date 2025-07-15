package client

import (
	pb "github.com/travis-james/DBCache/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb.DBCacheServiceClient
	*grpc.ClientConn
}

func Start() (*Client, error) {
	addr := "localhost:50051"
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		pb.NewDBCacheServiceClient(conn),
		conn,
	}, nil
}

func (cc *Client) Close() {
	cc.ClientConn.Close()
}
