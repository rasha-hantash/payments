// grpc/service.go
package grpc

import (
	"context"
	"log/slog"

	pb "github.com/rasha-hantash/chariot-takehome/api/grpc/proto"
	"github.com/rasha-hantash/chariot-takehome/api/grpc/repository"
	lg "github.com/rasha-hantash/chariot-takehome/api/pkgs/logger"
)

type GrpcService struct {
	UserRepo        *repository.UserRepository
	AccountRepo     *repository.AccountRepository
	TransactionRepo *repository.TransactionRepository
	pb.UnimplementedApiServiceServer
}

func (g *GrpcService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	ctx = lg.AppendCtx(ctx, slog.String("email", req.Email), slog.String("name", req.Name)) 
	slog.InfoContext(ctx, "creating user")
	
	user := &repository.User{
		Name:  req.Name,
		Email: req.Email,
	}

	res, err := g.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Id:                 res.Id,
		IntLedgerAccountId: res.IntLedgerAccountId.String,
		ExtLedgerAccountId: res.ExtLedgerAccountId.String,
	}, nil
}

func (g *GrpcService) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
	ctx = lg.AppendCtx(ctx, slog.String("account_state", req.AccountState), slog.String("account_type", req.AccountType))
	slog.InfoContext(ctx, "creating account")
	
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
	ctx = lg.AppendCtx(ctx, slog.Float64("amount", req.Amount), slog.String("user_id", req.UserId), slog.String("debit_account_id", req.DebitAccountId), slog.String("credit_account_id", req.CreditAccountId))
	slog.InfoContext(ctx, "depositing funds")

	id, err := g.TransactionRepo.DepositFunds(ctx, req.Amount, req.UserId, req.DebitAccountId, req.CreditAccountId)
	if err != nil {
		return nil, err
	}
	return &pb.Transaction{Id: id}, nil
}

func (g *GrpcService) WithdrawFunds(ctx context.Context, req *pb.WithdrawFundsRequest) (*pb.Transaction, error) {
	ctx = lg.AppendCtx(ctx, slog.Float64("amount", req.Amount), slog.String("user_id", req.UserId), slog.String("debit_account_id", req.DebitAccountId), slog.String("credit_account_id", req.CreditAccountId))
	slog.InfoContext(ctx, "withdrawing funds")
	
	id, err := g.TransactionRepo.WithdrawFunds(ctx, req.Amount, req.UserId, req.DebitAccountId, req.CreditAccountId)
	if err != nil {
		return nil, err
	}
	return &pb.Transaction{Id: id}, nil
}

func (g *GrpcService) TransferFunds(ctx context.Context, req *pb.TransferFundsRequest) (*pb.Transaction, error) {
	ctx = lg.AppendCtx(ctx, slog.Float64("amount", req.Amount), slog.String("user_id", req.UserId), slog.String("debit_account_id", req.DebitAccountId), slog.String("credit_account_id", req.CreditAccountId))
	slog.InfoContext(ctx, "transferring funds")
	
	id, err := g.TransactionRepo.TransferFunds(ctx, req.Amount, req.UserId, req.DebitAccountId, req.CreditAccountId)
	if err != nil {
		return nil, err
	}
	return &pb.Transaction{Id: id}, nil
}

func (g *GrpcService) GetAccountBalance(ctx context.Context, req *pb.GetAccountBalanceRequest) (*pb.AccountBalance, error) {
	ctx = lg.AppendCtx(ctx, slog.String("account_id", req.AccountId), slog.Time("timestamp", req.AtTime.AsTime()))
	slog.InfoContext(ctx, "getting account balance")
	
	balance, err := g.AccountRepo.GetAccountBalance(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}
	if balance == 0 {
		return &pb.AccountBalance{Balance: 0}, nil
	}
	return &pb.AccountBalance{Balance: float64(balance / 100)}, nil
}

func (g *GrpcService) ListTransactions(ctx context.Context, req *pb.ListTransactionsRequest) (*pb.ListTransactionsResponse, error) {
	ctx = lg.AppendCtx(ctx, slog.String("account_id", req.AccountId), slog.String("cursor", req.Cursor), slog.Int("page_size", int(req.PageSize)))
	slog.InfoContext(ctx, "listing transactions")
	
	var filter repository.TransactionFilter

	filter.AccountID = &req.AccountId
	filter.Cursor = &req.Cursor
	filter.Limit = nil // todo change this to req.PageSize

	transactions, nextCursor, err := g.TransactionRepo.ListTransactions(ctx, &filter)
	if err != nil {
		return nil, err
	}
	var pbTransactions []*pb.Transaction
	for _, t := range transactions {
		pbTransactions = append(pbTransactions, &pb.Transaction{
			Id:        t.Id,
			Amount:    float64(t.Amount),
			AccountId: t.AccountId,
			Direction: t.Direction,
		})
	}

	return &pb.ListTransactionsResponse{
		Transactions: pbTransactions,
		NextCursor:   nextCursor,
	}, nil
}
