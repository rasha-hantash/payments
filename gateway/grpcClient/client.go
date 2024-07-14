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
	slog.Info("connecting to api client", "addr", serverAddr)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		slog.Error("error connection to api service", err.Error())
		return nil, nil
	}
	
	return &ApiClient{client: pb.NewApiServiceClient(conn), Conn: conn}, nil
}




func (c *ApiClient) CreateUser(req *pb.CreateUserRequest) (*pb.User, error) {
    slog.Info("gateway - creating user")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
	// err := c.Conn.Invoke(ctx, "/api.ApiService/CreateUser", req, nil)
	user, err := c.client.CreateUser(ctx, req)
	if err != nil {
		slog.Error("api - error creating user", "error", err.Error())
		return nil, err
	}
	return user, nil
}
