// grpc/repository/user.go
package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/rasha-hantash/chariot-takehome/api/pkgs/identifier"
)

type User struct {
	Id                 string
	Email              string
	Name               string
	IntLedgerAccountId sql.NullString
	ExtLedgerAccountId sql.NullString
}

type UserRepository struct {
	db          *sql.DB
	accountRepo *AccountRepository
	ID          identifier.ID
}

func NewUserRepository(db *sql.DB, accountRepo *AccountRepository, prefix string) *UserRepository {
	return &UserRepository{db: db, accountRepo: accountRepo, ID: identifier.ID(prefix)}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var userId string
	// Create internal and external ledger accounts
	intAccount := &Account{
		AccountState: "open",
		AccountType:  "debit",
	}
	extAccount := &Account{
		AccountState: "open",
		AccountType:  "credit",
	}

	intAccountId, err := r.accountRepo.CreateAccount(ctx, intAccount)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating internal account", "error", err)
		return nil, err
	}
	extAccountId, err := r.accountRepo.CreateAccount(ctx, extAccount)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating external account", "error", err)
		return nil, err
	}

	idUser := r.ID.New()
	user.IntLedgerAccountId = sql.NullString{String: intAccountId, Valid: true}
	user.ExtLedgerAccountId = sql.NullString{String: extAccountId, Valid: true}

	err = r.db.QueryRowContext(ctx, `
        INSERT INTO users (id, email, name, int_ledger_account_id, ext_ledger_account_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, idUser, user.Email, user.Name, user.IntLedgerAccountId, user.ExtLedgerAccountId).Scan(&userId)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating user", "error", err)
		return nil, err
	}

	return &User{
		Id:                 userId,
		IntLedgerAccountId: user.IntLedgerAccountId,
		ExtLedgerAccountId: user.ExtLedgerAccountId,
	}, err
}
