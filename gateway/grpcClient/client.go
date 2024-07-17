package grpcClient

import (
	"context"
	"log/slog"

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

func (c *ApiClient) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user, err := c.client.CreateUser(ctx, req)
	if err != nil {
		slog.Error("error creating user", "error", err.Error())
		return nil, err
	}
	return user, nil
}

func (c *ApiClient) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
	account, err := c.client.CreateAccount(ctx, req)
	if err != nil {
		slog.Error("error creating account", "error", err.Error())
		return nil, err
	}
	return account, nil
}

func (c *ApiClient) DepositFunds(ctx context.Context, req *pb.DepositFundsRequest) (*pb.Transaction, error) {
	resp, err := c.client.DepositFunds(ctx, req)
	if err != nil {
		slog.Error("error depositing funds", "error", err.Error())
		return nil, err
	}
	return resp, nil
}

func (c *ApiClient) WithdrawFunds(ctx context.Context, req *pb.WithdrawFundsRequest) (*pb.Transaction, error) {
	resp, err := c.client.WithdrawFunds(ctx, req)
	if err != nil {
		slog.Error("error withdrawing funds", "error", err.Error())
		return nil, err
	}
	return resp, nil
}

func (c *ApiClient) TransferFunds(ctx context.Context, req *pb.TransferFundsRequest) (*pb.Transaction, error) {
	resp, err := c.client.TransferFunds(ctx, req)
	if err != nil {
		slog.Error("error transferring funds", "error", err.Error())
		return nil, err
	}
	return resp, nil
}

func (c *ApiClient) ListTransactions(ctx context.Context, req *pb.ListTransactionsRequest) (*pb.ListTransactionsResponse, error) {
	resp, err := c.client.ListTransactions(ctx, req)
	if err != nil {
		slog.Error("error listing transactions", "error", err.Error())
		return nil, err
	}
	return resp, nil
}

func (c *ApiClient) GetAccountBalance(ctx context.Context, req *pb.GetAccountBalanceRequest) (*pb.AccountBalance, error) {
	resp, err := c.client.GetAccountBalance(ctx, req)
	if err != nil {
		slog.Error("error getting account balance", "error", err.Error())
		return nil, err
	}
	return resp, nil
}
