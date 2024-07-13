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

// func NewGrpcService(userRepo *repository.UserRepository) *GrpcService {
// 	return &GrpcService{
// 		UserRepo: 
// 		userRepo}
// }

func (g *GrpcService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user := &repository.User{
		Email: req.Email,
		//FirstName: req.FirstName,
		//LastName:  req.LastName,
		// Set IntLedgerAccountId and ExtLedgerAccountId if needed
	}

	id, err := g.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id: id,
	}, nil
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

	return &pb.Account{
		Id: id,
	}, nil
}

func (g *GrpcService) DepositFunds(ctx context.Context, req *pb.DepositFunds) (*pb.Transaction, error) {
	// Get the account
	// Get the transaction
	// Update the account
	// Update the transaction
	// Return the response

	// id, err := g.TransactionRepo.DepositFunds(ctx, req.userId, req.debitAccount, req.CrediAccount, req.Amount)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

func (g *GrpcService) WithdrawFunds(ctx context.Context, req *pb.WithdrawFundsRequest) (*pb.Transaction, error) {
	
	

	id, err := g.TransactionRepo.WithdrawFunds(ctx, req.userId, req.debitAccount, req.CrediAccount, req.Amount)
	if err != nil {
		return nil, err
	}

	return &pb.Transaction{
		Id: id,
	}, nil
}

func (g *GrpcService) TransferFunds(ctx context.Context, req *pb.TransferFundsRequest) (*pb.Transaction, error) {
	// Get the account
	// Get the transaction
	// Update the account
	// Update the transaction
	// Return the response
	return nil, nil
}
