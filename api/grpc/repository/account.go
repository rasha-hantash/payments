package repository

import (
	"context"
	"database/sql"
	"log/slog"
)

type Account struct {
	Id           string
	AccountState string
	AccountType  string
}

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (a *AccountRepository) CreateAccount(ctx context.Context, account *Account) (string, error) {
	var id string
	err := a.db.QueryRowContext(ctx, "INSERT INTO accounts (id, account_state, account_type) VALUES ($1, $2, $3)",
		account.Id, account.AccountState, account.AccountType).Scan(&id)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating account", "error", err)
		return "", err
	}
	return id, err
}
