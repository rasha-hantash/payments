package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"

	"github.com/rasha-hantash/chariot-takehome/api/pkgs/test"
	"github.com/testcontainers/testcontainers-go"
)

func TestTransactionRepository_DepositFunds(t *testing.T) {
	db, container := test.SetupAndFillDatabaseContainer("seed_transactions_deposit_funds.sql")
	defer func(container testcontainers.Container) {
		err := test.TeardownDatabaseContainer(container)
		if err != nil {
			log.Fatalf("failed to close container down: %v\n", err)
		}
	}(container)
	defer db.Close()

	tests := []struct {
		name            string
		amount          float64
		userId          string
		debitAccountId  string
		creditAccountId string
		expectedResult  string
		wantErr         bool
	}{
		{
			name:            "successful deposit",
			amount:          100,
			userId:          "usr_1",
			debitAccountId:  "acct_1",
			creditAccountId: "acct_2",
			expectedResult:  "success",
			wantErr:         false,
		},
		{
			name:            "insert error,  debit account does not exist",
			amount:          100,
			userId:          "usr_2",
			debitAccountId:  "acct_3",
			creditAccountId: "acct_4",
			expectedResult:  "",
			wantErr:         true,
		},
		{
			name:            "insert error, credit account does not exist",
			amount:          100,
			userId:          "usr_3",
			debitAccountId:  "acct_5",
			creditAccountId: "acct_6",
			expectedResult:  "",
			wantErr:         true,
		},
	}

	repo := NewTransactionRepository(db, "txn_", "le_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.DepositFunds(context.Background(), tt.amount, tt.userId, tt.debitAccountId, tt.creditAccountId)

			if result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

		})
	}
}

func TestTransactionRepository_WithdrawFunds(t *testing.T) {
	db, container := test.SetupAndFillDatabaseContainer("seed_transactions_withdraw_funds.sql")
	defer func(container testcontainers.Container) {
		err := test.TeardownDatabaseContainer(container)
		if err != nil {
			log.Fatalf("failed to close container down: %v\n", err)
		}
	}(container)
	defer db.Close()

	tests := []struct {
		name            string
		amount          float64
		userId          string
		debitAccountId  string
		creditAccountId string
		expectedResult  string
		wantErr         bool
	}{
		{
			name:            "successful withdraw",
			amount:          100,
			userId:          "usr_1",
			debitAccountId:  "acct_1",
			creditAccountId: "acct_2",
			expectedResult:  "success",
			wantErr:         false,
		},
		{
			name:            "insufficient funds",
			amount:          100,
			userId:          "usr_2",
			debitAccountId:  "acct_3",
			creditAccountId: "acct_4",
			expectedResult:  "",
			wantErr:         true,
		},
		{
			name:            "credit account does not exist",
			amount:          100,
			userId:          "usr_3",
			debitAccountId:  "acct_5",
			creditAccountId: "acct_6",
			expectedResult:  "",
			wantErr:         true,
		},
	}

	repo := NewTransactionRepository(db, "txn_", "le_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.WithdrawFunds(context.Background(), tt.amount, tt.userId, tt.debitAccountId, tt.creditAccountId)

			if result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestTransactionRepository_TransferFunds(t *testing.T) {
	db, container := test.SetupAndFillDatabaseContainer("seed_transactions_withdraw_funds.sql")
	defer func(container testcontainers.Container) {
		err := test.TeardownDatabaseContainer(container)
		if err != nil {
			log.Fatalf("failed to close container down: %v\n", err)
		}
	}(container)
	defer db.Close()

	tests := []struct {
		name            string
		amount          float64
		userId          string
		debitAccountId  string
		creditAccountId string
		expectedResult  string
		wantErr         bool
	}{
		{
			name:            "successful transfer",
			amount:          100,
			userId:          "usr_1",
			debitAccountId:  "acct_1",
			creditAccountId: "acct_2",
			expectedResult:  "success",
			wantErr:         false,
		},
		{
			name:            "insert error debit insufficient funds",
			amount:          100,
			userId:          "usr_2",
			debitAccountId:  "acct_3",
			creditAccountId: "acct_4",
			expectedResult:  "",
			wantErr:         true,
		},
		{
			name:            "insert error debit insufficient funds",
			amount:          100,
			userId:          "usr_1",
			debitAccountId:  "acct_5",
			creditAccountId: "acct_6",
			expectedResult:  "",
			wantErr:         true,
		},
	}

	repo := NewTransactionRepository(db, "txn_", "le_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.TransferFunds(context.Background(), tt.amount, tt.userId, tt.debitAccountId, tt.creditAccountId)

			if result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestTransactionRepository_ListTransactions(t *testing.T) {
	db, container := test.SetupAndFillDatabaseContainer("seed_transactions_list_transactions.sql")
	defer func(container testcontainers.Container) {
		err := test.TeardownDatabaseContainer(container)
		if err != nil {
			log.Fatalf("failed to close container down: %v\n", err)
		}
	}(container)
	defer db.Close()

	tests := []struct {
		name           string
		filter         TransactionFilter
		expectedResult []*Transaction
		expectedCursor string
		wantErr        bool
	}{
		{
			name: "successful list for acct_1",
			filter: TransactionFilter{
				AccountID: "acct_1",
				Cursor:    "",
				Limit:     5,
			},
			expectedResult: []*Transaction{
				{Id: "txn_1", AccountId: "acct_1", Amount: 110, Direction: "debit"},
				{Id: "txn_3", AccountId: "acct_1", Amount: 130, Direction: "credit"},
				{Id: "txn_4", AccountId: "acct_1", Amount: 140, Direction: "debit"},
				{Id: "txn_6", AccountId: "acct_1", Amount: 160, Direction: "credit"},
				{Id: "txn_7", AccountId: "acct_1", Amount: 170, Direction: "debit"},
			},
			expectedCursor: "txn_9",
			wantErr:        false,
		},
		{
			name: "successful list for acct_2",
			filter: TransactionFilter{
				AccountID: "acct_2",
				Cursor:    "",
				Limit:     5,
			},
			expectedResult: []*Transaction{
				{Id: "txn_1", AccountId: "acct_2", Amount: 110, Direction: "credit"},
				{Id: "txn_2", AccountId: "acct_2", Amount: 120, Direction: "debit"},
				{Id: "txn_4", AccountId: "acct_2", Amount: 140, Direction: "credit"},
				{Id: "txn_5", AccountId: "acct_2", Amount: 150, Direction: "debit"},
				{Id: "txn_7", AccountId: "acct_2", Amount: 170, Direction: "credit"},
			},
			expectedCursor: "txn_8",
			wantErr:        false,
		},
		{
			name: "successful list for acct_3",
			filter: TransactionFilter{
				AccountID: "acct_3",
				Cursor:    "",
				Limit:     5,
			},
			expectedResult: []*Transaction{
				{Id: "txn_2", AccountId: "acct_3", Amount: 120, Direction: "credit"},
				{Id: "txn_3", AccountId: "acct_3", Amount: 130, Direction: "debit"},
				{Id: "txn_5", AccountId: "acct_3", Amount: 150, Direction: "credit"},
				{Id: "txn_6", AccountId: "acct_3", Amount: 160, Direction: "debit"},
				{Id: "txn_8", AccountId: "acct_3", Amount: 180, Direction: "credit"},
			},
			expectedCursor: "txn_9",
			wantErr:        false,
		},
		{
			name: "successful list with cursor",
			filter: TransactionFilter{
				AccountID: "acct_1",
				Cursor:    "txn_7",
				Limit:     3,
			},
			expectedResult: []*Transaction{
				{Id: "txn_9", AccountId: "acct_1", Amount: 190, Direction: "credit"},
				{Id: "txn_10", AccountId: "acct_1", Amount: 200, Direction: "debit"},
				{Id: "txn_12", AccountId: "acct_1", Amount: 220, Direction: "credit"},
			},
			expectedCursor: "txn_13",
			wantErr:        false,
		},
		{
			name: "no transactions found",
			filter: TransactionFilter{
				AccountID: "non_existent_account",
				Cursor:    "",
				Limit:     5,
			},
			expectedResult: []*Transaction{},
			expectedCursor: "",
			wantErr:        false,
		},
	}

	repo := NewTransactionRepository(db, "txn_", "le_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, cursor, err := repo.ListTransactions(context.Background(), tt.filter)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
				assert.Equal(t, tt.expectedCursor, cursor)
			}
		})
	}
}
