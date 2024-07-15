package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"fmt"

	"github.com/rasha-hantash/chariot-takehome/api/pkgs/identifier"
)

type TransactionRepository struct {
	db *sql.DB
	ID identifier.ID
}

type Transaction struct {
	Id        string
	AccountId string
	Amount    int64
	Direction string
}

type TransactionFilter struct {
	AccountID string
	Cursor    string
	Limit     int
}



func NewTransactionRepository(db *sql.DB, prefix string) *TransactionRepository {
	return &TransactionRepository{db: db, ID: identifier.ID(prefix)}
}

func (t *TransactionRepository) ListTransactions(ctx context.Context, filter TransactionFilter) ([]*Transaction, string, error) {
	query := "SELECT id, account_id, amount, direction FROM transactions WHERE 1=1"
	args := []interface{}{}

	if filter.AccountID != "" {
		query += " AND account_id = ?"
		args = append(args, filter.AccountID)
	}

	if filter.Cursor != "" {
		query += " AND id > ?"
		args = append(args, filter.Cursor)
	}

	query += " ORDER BY id ASC LIMIT ?"
	args = append(args, filter.Limit+1) // Fetch one extra to determine if there are more results

	rows, err := t.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "error while listing transactions", "error", err)
		return nil, "", err
	}
	defer rows.Close()

	transactions := make([]*Transaction, 0, filter.Limit)
	var nextCursor string

	for rows.Next() {
		if len(transactions) == filter.Limit {
			// We've reached the limit, so this row determines if there are more results
			var lastID string
			if err := rows.Scan(&lastID); err == nil {
				nextCursor = lastID
			}
			break
		}

		transaction := &Transaction{}
		err := rows.Scan(&transaction.Id, &transaction.AccountId, &transaction.Amount, &transaction.Direction)
		if err != nil {
			slog.ErrorContext(ctx, "error while scanning transaction", "error", err)
			return nil, "", err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "error after iterating rows", "error", err)
		return nil, "", err
	}

	return transactions, nextCursor, nil
}

// TODO: Add Comment
func (t *TransactionRepository) DepositFunds(ctx context.Context, amount float64, userId, debitAccountId, creditAccountId string) (string, error) {
	return t.addDoubleEntryTransaction(ctx, amount, debitAccountId, creditAccountId, userId)
}

// TODO: Add Comment
func (t *TransactionRepository) WithdrawFunds(ctx context.Context, amount float64, userId, debitAccountId, creditAccountId string) (string, error) {
	return t.addDoubleEntryTransaction(ctx, amount, debitAccountId, creditAccountId, userId)
}

// TODO: Add Comment
func (t *TransactionRepository) TransferFunds(ctx context.Context, amount float64, userId, debitAccountId, creditAccountId string) (string, error) {
	return t.addDoubleEntryTransaction(ctx, amount, debitAccountId, creditAccountId, userId)
}

func (t *TransactionRepository) addDoubleEntryTransaction(ctx context.Context, amount float64, debitedAccountId, creditedAccountId, userId string) (string, error) {
	tx, err := t.db.BeginTx(ctx,  &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Check if the debited account has sufficient balance
    sufficient, err := t.checkSufficientBalance(ctx, tx, debitedAccountId, amount)
    if err != nil {
		slog.Error("error checking balance", "error", err.Error())
        return "", fmt.Errorf("error checking balance: %w", err)
    }
    if !sufficient {
		slog.Error("insufficient balance", "account_id", debitedAccountId)
        return "", fmt.Errorf("insufficient balance in account %s", debitedAccountId)
    }
	
	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (account_id, amount, direction, created_by) VALUES ($1, $2, $3, $4, $5)",
	debitedAccountId, amount*100, "debit", userId)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating credit transaction", "error", err)
		return "", err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (account_id, amount, direction, created_by) VALUES ($1, $2, $3, $4, $5)",
	creditedAccountId, amount*100, "credit", userId)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating credit transaction", "error", err)
		return "", err
	}

	return "success", nil
}

func (t *TransactionRepository) checkSufficientBalance(ctx context.Context, tx *sql.Tx, accountId string, amount float64) (bool, error) {
    var balance float64
    err := tx.QueryRowContext(ctx, `
        SELECT 
			SUM(CASE WHEN direction = 'credit' THEN amount ELSE -amount END) AS balance
		FROM
			ledger_entries
		WHERE
			account_id = '%s'
    `, accountId).Scan(&balance)
    if err != nil {
        return false, err
    }
    return balance >= amount, nil
}



