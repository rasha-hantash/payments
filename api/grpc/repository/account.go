package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"fmt"

	"github.com/rasha-hantash/chariot-takehome/api/pkgs/identifier"
)

type Account struct {
	Id           string
	AccountState string
	AccountType  string
}

type AccountRepository struct {
	db *sql.DB
	ID identifier.ID
}

func NewAccountRepository(db *sql.DB, prefix string) *AccountRepository {
	return &AccountRepository{db: db, ID: identifier.ID(prefix)}
}

func (a *AccountRepository) CreateAccount(ctx context.Context, account *Account) (string, error) {
	var id string
	hrId := a.ID.New()

	err := a.db.QueryRowContext(ctx, "INSERT INTO accounts (id, account_state, account_type) VALUES ($1, $2, $3) RETURNING id",
		hrId, account.AccountState, account.AccountType).Scan(&id)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating account", "error", err)
		return "", err
	}
	return id, err
}

func (a *AccountRepository) GetAccountBalance(ctx context.Context, accountId string) (int64, error) {
	var balance int64
	fmt.Println("accountId", accountId)
	err := a.db.QueryRowContext(ctx, "SELECT COALESCE(SUM(CASE WHEN direction = 'credit' THEN amount ELSE -amount END), 0) as balance FROM ledger_entries WHERE account_id = $1", accountId).Scan(&balance)
	if err != nil {
		slog.ErrorContext(ctx, "error while getting account balance", "error", err)
		return 0, err
	}
	fmt.Println("balance", balance)
	return balance, nil
}
