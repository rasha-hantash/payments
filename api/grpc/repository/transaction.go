package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/rasha-hantash/chariot-takehome/api/pkgs/identifier"
)

type TransactionRepository struct {
	db       *sql.DB
	txnID    identifier.ID
	ledgerID identifier.ID
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

func NewTransactionRepository(db *sql.DB, txnPrefix, ledgerPrefix string) *TransactionRepository {
	return &TransactionRepository{db: db, txnID: identifier.ID(txnPrefix), ledgerID: identifier.ID(ledgerPrefix)}
}

func (t *TransactionRepository) ListTransactions(ctx context.Context, filter TransactionFilter) ([]*Transaction, string, error) {
	query := "SELECT id, account_id, amount, direction FROM ledger_entries WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if filter.AccountID != "" {
		query += fmt.Sprintf(" AND account_id = $%d", argCount)
		args = append(args, filter.AccountID)
		argCount++
	}

	if filter.Cursor != "" {
		query += fmt.Sprintf(" AND id > $%d", argCount)
		args = append(args, filter.Cursor)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY id ASC LIMIT $%d", argCount)
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
	return t.addDoubleEntryTransactionFromExternal(ctx, amount, debitAccountId, creditAccountId, userId)
}

// TODO: Add Comment
func (t *TransactionRepository) WithdrawFunds(ctx context.Context, amount float64, userId, debitAccountId, creditAccountId string) (string, error) {
	return t.addDoubleEntryTransaction(ctx, amount, debitAccountId, creditAccountId, userId)
}

// TODO: Add Comment
func (t *TransactionRepository) TransferFunds(ctx context.Context, amount float64, userId, debitAccountId, creditAccountId string) (string, error) {
	return t.addDoubleEntryTransaction(ctx, amount, debitAccountId, creditAccountId, userId)
}

func (t *TransactionRepository) addDoubleEntryTransactionFromExternal(ctx context.Context, amount float64, debitedAccountId, creditedAccountId, userId string) (string, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// note in this mvp, we are not checking to see if a user has enough funds in their external account to deposit funds
	// in the next iteration i would rely on the third party api to determine that

	txnId := t.txnID.New()
	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (id, amount, status, created_by) VALUES ($1, $2, $3, $4)",
		txnId, amount*100, "success", userId)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating transaction", "error", err)
		return "", err
	}

	ledgerId1 := t.ledgerID.New()
	_, err = tx.ExecContext(ctx, "INSERT INTO ledger_entries (id, transaction_id, account_id, amount, direction, created_by) VALUES ($1, $2, $3, $4, $5, $6)",
		ledgerId1, txnId, debitedAccountId, amount*100, "debit", userId)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating credit ledger entry", "error", err)
		return "", err
	}

	ledgerId2 := t.ledgerID.New()
	_, err = tx.ExecContext(ctx, "INSERT INTO ledger_entries (id,  transaction_id, account_id, amount, direction, created_by) VALUES ($1, $2, $3, $4, $5, $6)",
		ledgerId2, txnId, creditedAccountId, amount*100, "credit", userId)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating credit ledger entry", "error", err)
		return "", err
	}

	return "success", nil
}

func (t *TransactionRepository) addDoubleEntryTransaction(ctx context.Context, amount float64, debitedAccountId, creditedAccountId, userId string) (string, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
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

	txnId := t.txnID.New()
	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (id, amount, status, created_by) VALUES ($1, $2, $3, $4)",
		txnId, amount*100, "success", userId)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating transaction", "error", err)
		return "", err
	}

	ledgerId1 := t.ledgerID.New()
	_, err = tx.ExecContext(ctx, "INSERT INTO ledger_entries (id, transaction_id, account_id, amount, direction, created_by) VALUES ($1, $2, $3, $4, $5, $6)",
		ledgerId1, txnId, debitedAccountId, amount*100, "debit", userId)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating credit ledger entry", "error", err)
		return "", err
	}

	ledgerId2 := t.ledgerID.New()
	_, err = tx.ExecContext(ctx, "INSERT INTO ledger_entries (id,  transaction_id, account_id, amount, direction, created_by) VALUES ($1, $2, $3, $4, $5, $6)",
		ledgerId2, txnId, creditedAccountId, amount*100, "credit", userId)
	if err != nil {
		slog.ErrorContext(ctx, "error while creating credit ledger entry", "error", err)
		return "", err
	}

	return "success", nil
}

func (t *TransactionRepository) checkSufficientBalance(ctx context.Context, tx *sql.Tx, accountId string, amount float64) (bool, error) {
	var balance float64
	err := tx.QueryRowContext(ctx, `
        SELECT 
            COALESCE(SUM(CASE WHEN direction = 'credit' THEN amount ELSE -amount END), 0) AS balance
        FROM
            ledger_entries
        WHERE
            account_id = $1
    `, accountId).Scan(&balance)
	if err != nil {
		return false, err
	}
	return balance >= amount, nil
}
