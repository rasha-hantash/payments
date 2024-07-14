// grpc/service.go
package grpc

import (
	"context"
	pb "github.com/rasha-hantash/chariot-takehome/api/grpc/proto"
	"github.com/rasha-hantash/chariot-takehome/api/grpc/repository"
)

type GrpcService struct {
	UserRepo *repository.UserRepository
	AccountRepo *repository.AccountRepository
	TransactionRepo *repository.TransactionRepository
	pb.UnimplementedApiServiceServer
}

func (g *GrpcService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user := &repository.User{
		Name: req.Name,
		Email: req.Email,
	}

	id, err := g.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &pb.User{Id: id}, nil
}

func (g *GrpcService) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
	account := &repository.Account{
		AccountState: req.AccountState,
		AccountType:  req.AccountType,
	}

	id, err := g.AccountRepo.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	return &pb.Account{Id: id}, nil
}

func (g *GrpcService) DepositFunds(ctx context.Context, req *pb.DepositFundsRequest) (*pb.Transaction, error) {
	id, err := g.TransactionRepo.DepositFunds(ctx, req.Amount, req.UserId, req.DebitAccountId, req.CreditAccountId)
	if err != nil {
		return nil, err
	}
	return &pb.Transaction{Id: id}, nil
}

func (g *GrpcService) WithdrawFunds(ctx context.Context, req *pb.WithdrawFundsRequest) (*pb.Transaction, error) {
	id, err := g.TransactionRepo.WithdrawFunds(ctx, req.Amount, req.UserId, req.DebitAccountId, req.CreditAccountId)
	if err != nil {
		return nil, err
	}
	return &pb.Transaction{Id: id}, nil
}

func (g *GrpcService) TransferFunds(ctx context.Context, req *pb.TransferFundsRequest) (*pb.Transaction, error) {
	id, err := g.TransactionRepo.TransferFunds(ctx, req.Amount, req.UserId, req.DebitAccountId, req.CreditAccountId)
	if err != nil {
		return nil, err
	}
	return &pb.Transaction{Id: id}, nil
}

func (g *GrpcService) GetAccountBalance(ctx context.Context, req *pb.GetAccountBalanceRequest) (*pb.AccountBalance, error) {
	balance, err := g.AccountRepo.GetAccountBalance(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}
	return &pb.AccountBalance{Balance: float64(balance / 100)}, nil
}
