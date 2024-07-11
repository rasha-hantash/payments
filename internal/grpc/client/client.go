package client

import (
	"context"

	"google.golang.org/grpc"
	pb "myproject/pkgs/proto/api"
)

type ApiClient struct {
	client pb.ApiServiceClient
}

func NewApiClient(addr string) (*ApiClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &ApiClient{client: pb.NewApiServiceClient(conn)}, nil
}

func (c *ApiClient) CreateUser(ctx context.Context, name string, email string) error {
	_, err := c.client.CreateUser(ctx, &pb.CreateUserRequest{Name: name, Email: email})
	return err
}
