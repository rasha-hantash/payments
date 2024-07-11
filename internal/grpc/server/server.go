package server

import (
	"context"

	pb "myproject/pkgs/proto/api"
)

type GrpcServer struct {
	pb.UnimplementedUserServiceServer
}

func (g *GrpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Implement user creation logic here
	return &pb.CreateUserResponse{}, nil
}
