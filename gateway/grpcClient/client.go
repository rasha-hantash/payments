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
		slog.Error("error connection to api service", "error", err.Error())
		return nil, nil
	}
	
	return &ApiClient{client: pb.NewApiServiceClient(conn), Conn: conn}, nil
}


func (c *ApiClient) CreateUser(req *pb.CreateUserRequest) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
	user, err := c.client.CreateUser(ctx, req)
	if err != nil {
		slog.Error("error creating user", "error", err.Error())
		return nil, err
	}
	return user, nil
}

func (c *ApiClient) CreateAccount(req *pb.CreateAccountRequest) (*pb.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
	account, err := c.client.CreateAccount(ctx, req)
	if err != nil {
		slog.Error("error creating user", "error", err.Error())
		return nil, err
	}
	return account, nil
}

func (c *ApiClient) DepositFunds(req *pb.DepositFundsRequest) (*pb.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
	resp, err := c.client.DepositFunds(ctx, req)
	if err != nil {
		slog.Error("error creating user", "error", err.Error())
		return nil, err
	}
	return resp, nil
}


func (c *ApiClient) WithdrawFunds(req *pb.WithdrawFundsRequest) (*pb.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
	resp, err := c.client.WithdrawFunds(ctx, req)
	if err != nil {
		slog.Error("error creating user", "error", err.Error())
		return nil, err
	}
	return resp, nil
}

func (c *ApiClient) TransferFunds(req *pb.TransferFundsRequest) (*pb.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := c.client.TransferFunds(ctx, req)
	if err != nil {
		slog.Error("error creating user", "error", err.Error())
		return nil, err
	}
	return resp, nil
}

func (c *ApiClient) ListTransactions(req *pb.ListTransactionsRequest) (*pb.ListTransactionsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := c.client.ListTransactions(ctx, req)
	if err != nil {
		slog.Error("error creating user", "error", err.Error())
		return nil, err
	}
	return resp, nil
}


func (c *ApiClient) GetAccountBalance(req *pb.GetAccountBalanceRequest) (*pb.AccountBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := c.client.GetAccountBalance(ctx, req)
	if err != nil {
		slog.Error("error creating user", "error", err.Error())
		return nil, err
	}
	return resp, nil
}
