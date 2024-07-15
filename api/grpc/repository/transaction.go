package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

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
	Status    string
	CreatedAt string
	// Direction string
}

type TransactionFilter struct {
	AccountID *string
	Cursor    *string
	Limit     *int
}

func NewTransactionRepository(db *sql.DB, txnPrefix, ledgerPrefix string) *TransactionRepository {
	return &TransactionRepository{db: db, txnID: identifier.ID(txnPrefix), ledgerID: identifier.ID(ledgerPrefix)}
}

func (t *TransactionRepository) ListTransactions(ctx context.Context, filter *TransactionFilter) ([]Transaction, string, error) {
	query := `
	SELECT DISTINCT t.id, le.account_id, t.amount, t.status, t.created_at
	FROM transactions t
	JOIN ledger_entries le ON t.id = le.transaction_id
	WHERE 1=1
`

	var args []interface{}
	var conditions []string

	if filter.AccountID != nil && *filter.AccountID != "" {
		conditions = append(conditions, "le.account_id = $"+fmt.Sprint(len(args)+1))
		args = append(args, *filter.AccountID)
	}

	if filter.Cursor != nil && *filter.Cursor != "" {
		conditions = append(conditions, "t.id > $"+fmt.Sprint(len(args)+1))
		args = append(args, *filter.Cursor)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.id"

	if filter.Limit != nil {
		query += " LIMIT $" + fmt.Sprint(len(args)+1)
		args = append(args, *filter.Limit)
	}

	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("error querying transactions: %v", err)
	}
	defer rows.Close()

	var transactions []Transaction
	var lastID string

	for rows.Next() {
		var txn Transaction
		err := rows.Scan(&txn.Id, &txn.AccountId, &txn.Amount, &txn.Status, &txn.CreatedAt)
		if err != nil {
			return nil, "", fmt.Errorf("error scanning transaction: %v", err)
		}
		txn.AccountId = *filter.AccountID
		transactions = append(transactions, txn)
		lastID = txn.Id
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating transactions: %v", err)
	}

	var nextCursor string
	if filter.Limit != nil && len(transactions) == *filter.Limit {
		nextCursor = lastID
	}

	return transactions, nextCursor, nil
}

// DepositFunds deposits funds into an account
func (t *TransactionRepository) DepositFunds(ctx context.Context, amount float64, userId, debitAccountId, creditAccountId string) (string, error) {
	return t.addDoubleEntryTransactionFromExternal(ctx, amount, debitAccountId, creditAccountId, userId)
}

// WithdrawFunds withdraws funds from an account
func (t *TransactionRepository) WithdrawFunds(ctx context.Context, amount float64, userId, debitAccountId, creditAccountId string) (string, error) {
	return t.addDoubleEntryTransaction(ctx, amount, debitAccountId, creditAccountId, userId)
}

// TransferFunds transfers funds from one account to another
func (t *TransactionRepository) TransferFunds(ctx context.Context, amount float64, userId, debitAccountId, creditAccountId string) (string, error) {
	return t.addDoubleEntryTransaction(ctx, amount, debitAccountId, creditAccountId, userId)
}

// addDoubleEntryTransactionFromExternal adds a transaction with a double ledger entry
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

	return string(txnId), nil
}

// addDoubleEntryTransaction adds a transaction with a double ledger entry
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

	return string(txnId), nil
}

// checkSufficientBalance checks if the account has sufficient balance to withdraw the amount
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
