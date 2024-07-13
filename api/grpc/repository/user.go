// grpc/repository/user.go
package repository

import (
	"context"
	"database/sql"
)

type User struct {
	Id                 string
	Email              string
	FirstName          string
	LastName           string
	IntLedgerAccountId sql.NullString
	ExtLedgerAccountId sql.NullString
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) (string, error) {
	var userId string
	// todo create a accounts (return internal and external ledger

	err := r.db.QueryRowContext(ctx, `
        INSERT INTO users (email, first_name, last_name, int_ledger_account_id, ext_ledger_account_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, user.Email, user.FirstName, user.LastName, user.IntLedgerAccountId, user.ExtLedgerAccountId).Scan(userId)
	if err != nil {
		return "", err
	}

	return userId, err
}
