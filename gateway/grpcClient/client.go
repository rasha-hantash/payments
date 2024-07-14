package grpcClient

import (
	"log/slog"
	"context"
	"time"

	pb "github.com/rasha-hantash/chariot-takehome/api/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ApiClient struct {
	client pb.ApiServiceClient
	Conn   *grpc.ClientConn
}

func NewApiClient(serverAddr string) (*ApiClient, error) {
	slog.Info("connecting to claims client", "addr", serverAddr)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		slog.Error("error connection to claims service", err)
		return nil, nil
	}
	return &ApiClient{client: pb.NewApiServiceClient(conn), Conn: conn}, nil
}




func (c *ApiClient) CreateUser(req *pb.CreateUserRequest) (*pb.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    return c.client.CreateUser(ctx, req)
}
