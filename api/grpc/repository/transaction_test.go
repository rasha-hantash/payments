package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTransactionRepository_DepositFunds(t *testing.T) {
	tests := []struct {
		name            string
		amount          float64
		userId          string
		debitAccountId  string
		creditAccountId string
		mockSetup       func(mock sqlmock.Sqlmock)
		expectedResult  string
		expectedErr     error
	}{
		{
			name:            "successful deposit",
			amount:          100,
			userId:          "user-1",
			debitAccountId:  "debit-1",
			creditAccountId: "credit-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("debit-1", 100, "debit", "user-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("credit-1", 100, "credit", "user-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedResult: "success",
			expectedErr:    nil,
		},
		{
			name:            "insert error on debit",
			amount:          100,
			userId:          "user-1",
			debitAccountId:  "debit-1",
			creditAccountId: "credit-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("debit-1", 100, "debit", "user-1").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedResult: "",
			expectedErr:    sql.ErrConnDone,
		},
		{
			name:            "insert error on credit",
			amount:          100,
			userId:          "user-1",
			debitAccountId:  "debit-1",
			creditAccountId: "credit-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("debit-1", 100, "debit", "user-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("credit-1", 100, "credit", "user-1").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedResult: "",
			expectedErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := NewTransactionRepository(db, "txn")
			result, err := repo.DepositFunds(context.Background(), tt.amount, tt.userId, tt.debitAccountId, tt.creditAccountId)

			if result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTransactionRepository_WithdrawFunds(t *testing.T) {
	tests := []struct {
		name            string
		amount          float64
		userId          string
		debitAccountId  string
		creditAccountId string
		mockSetup       func(mock sqlmock.Sqlmock)
		expectedResult  string
		expectedErr     error
	}{
		{
			name:            "successful withdraw",
			amount:          100,
			userId:          "user-1",
			debitAccountId:  "debit-1",
			creditAccountId: "credit-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("debit-1", 100, "debit", "user-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("credit-1", 100, "credit", "user-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedResult: "success",
			expectedErr:    nil,
		},
		{
			name:            "insert error on debit",
			amount:          100,
			userId:          "user-1",
			debitAccountId:  "debit-1",
			creditAccountId: "credit-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("debit-1", 100, "debit", "user-1").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedResult: "",
			expectedErr:    sql.ErrConnDone,
		},
		{
			name:            "insert error on credit",
			amount:          100,
			userId:          "user-1",
			debitAccountId:  "debit-1",
			creditAccountId: "credit-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("debit-1", 100, "debit", "user-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("credit-1", 100, "credit", "user-1").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedResult: "",
			expectedErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := NewTransactionRepository(db, "txn")
			result, err := repo.WithdrawFunds(context.Background(), tt.amount, tt.userId, tt.debitAccountId, tt.creditAccountId)

			if result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTransactionRepository_TransferFunds(t *testing.T) {
	tests := []struct {
		name            string
		amount          float64
		userId          string
		debitAccountId  string
		creditAccountId string
		mockSetup       func(mock sqlmock.Sqlmock)
		expectedResult  string
		expectedErr     error
	}{
		{
			name:            "successful transfer",
			amount:          100.50,
			userId:          "user-1",
			debitAccountId:  "debit-1",
			creditAccountId: "credit-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("debit-1", 100.50, "debit", "user-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("credit-1", 100.50, "credit", "user-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedResult: "success",
			expectedErr:    nil,
		},
		{
			name:            "insert error on debit",
			amount:          100.50,
			userId:          "user-1",
			debitAccountId:  "debit-1",
			creditAccountId: "credit-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("debit-1", 100.50, "debit", "user-1").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedResult: "",
			expectedErr:    sql.ErrConnDone,
		},
		{
			name:            "insert error on credit",
			amount:          100.50,
			userId:          "user-1",
			debitAccountId:  "debit-1",
			creditAccountId: "credit-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("debit-1", 100.50, "debit", "user-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs("credit-1", 100.50, "credit", "user-1").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedResult: "",
			expectedErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := NewTransactionRepository(db, "txn")
			result, err := repo.TransferFunds(context.Background(), tt.amount, tt.userId, tt.debitAccountId, tt.creditAccountId)

			if result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTransactionRepository_ListTransactions(t *testing.T) {
	tests := []struct {
		name           string
		filter         TransactionFilter
		mockSetup      func(mock sqlmock.Sqlmock)
		expectedResult []*Transaction
		expectedCursor string
		expectedErr    error
	}{
		{
			name: "successful list with limit",
			filter: TransactionFilter{
				AccountID: "acc-1",
				Cursor:    "",
				Limit:     2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "direction"}).
					AddRow("1", "acc-1", 100, "debit").
					AddRow("2", "acc-1", 200, "credit").
					AddRow("3", "acc-1", 300, "debit")
				mock.ExpectQuery("SELECT id, account_id, amount, direction FROM transactions WHERE 1=1 AND account_id = \\? ORDER BY id ASC LIMIT \\?").
					WithArgs("acc-1", 3).
					WillReturnRows(rows)
			},
			expectedResult: []*Transaction{
				{Id: "1", AccountId: "acc-1", Amount: 100, Direction: "debit"},
				{Id: "2", AccountId: "acc-1", Amount: 200, Direction: "credit"},
			},
			expectedCursor: "3",
			expectedErr:    nil,
		},
		{
			name: "no transactions found",
			filter: TransactionFilter{
				AccountID: "acc-2",
				Cursor:    "",
				Limit:     2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "direction"})
				mock.ExpectQuery("SELECT id, account_id, amount, direction FROM transactions WHERE 1=1 AND account_id = \\? ORDER BY id ASC LIMIT \\?").
					WithArgs("acc-2", 3).
					WillReturnRows(rows)
			},
			expectedResult: []*Transaction{},
			expectedCursor: "",
			expectedErr:    nil,
		},
		{
			name: "query error",
			filter: TransactionFilter{
				AccountID: "acc-3",
				Cursor:    "",
				Limit:     2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, account_id, amount, direction FROM transactions WHERE 1=1 AND account_id = \\? ORDER BY id ASC LIMIT \\?").
					WithArgs("acc-3", 3).
					WillReturnError(sql.ErrConnDone)
			},
			expectedResult: nil,
			expectedCursor: "",
			expectedErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := NewTransactionRepository(db, "txn")
			result, cursor, err := repo.ListTransactions(context.Background(), tt.filter)

			if len(result) != len(tt.expectedResult) {
				t.Errorf("expected result length %v, got %v", len(tt.expectedResult), len(result))
			} else {
				for i := range result {
					if *result[i] != *tt.expectedResult[i] {
						t.Errorf("expected result %v, got %v", tt.expectedResult[i], result[i])
					}
				}
			}
			if cursor != tt.expectedCursor {
				t.Errorf("expected cursor %v, got %v", tt.expectedCursor, cursor)
			}
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
