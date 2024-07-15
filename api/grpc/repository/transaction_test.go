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
				AccountID: strPtr("acct_1"),
				Cursor:    nil,
				Limit:     intPtr(5),
			},
			expectedResult: []*Transaction{
				{Id: "txn_1", Amount: 110, Status: "success"},
				{Id: "txn_3", Amount: 130, Status: "success"},
				{Id: "txn_4", Amount: 140, Status: "success"},
				{Id: "txn_6", Amount: 160, Status: "success"},
				{Id: "txn_7", Amount: 170, Status: "success"},
			},
			expectedCursor: "txn_7",
			wantErr:        false,
		},
		{
			name: "successful list for acct_2",
			filter: TransactionFilter{
				AccountID: strPtr("acct_2"),
				Cursor:    nil,
				Limit:     intPtr(5),
			},
			expectedResult: []*Transaction{
				{Id: "txn_1", Amount: 110, Status: "success"},
				{Id: "txn_2", Amount: 120, Status: "success"},
				{Id: "txn_4", Amount: 140, Status: "success"},
				{Id: "txn_5", Amount: 150, Status: "success"},
				{Id: "txn_7", Amount: 170, Status: "success"},
			},
			expectedCursor: "txn_7",
			wantErr:        false,
		},
		{
			name: "successful list for acct_3",
			filter: TransactionFilter{
				AccountID: strPtr("acct_3"),
				Cursor:    nil,
				Limit:     intPtr(5),
			},
			expectedResult: []*Transaction{
				{Id: "txn_2", Amount: 120, Status: "success"},
				{Id: "txn_3", Amount: 130, Status: "success"},
				{Id: "txn_5", Amount: 150, Status: "success"},
				{Id: "txn_6", Amount: 160, Status: "success"},
				{Id: "txn_8", Amount: 180, Status: "success"},
			},
			expectedCursor: "txn_8",
			wantErr:        false,
		},
		{
			name: "successful list with cursor",
			filter: TransactionFilter{
				AccountID: strPtr("acct_1"),
				Cursor:    strPtr("txn_7"),
				Limit:     intPtr(3),
			},
			expectedResult: []*Transaction{
				{Id: "txn_9", Amount: 190, Status: "success"},
				{Id: "txn_10", Amount: 200, Status: "success"},
				{Id: "txn_12", Amount: 220, Status: "success"},
			},
			expectedCursor: "txn_12",
			wantErr:        false,
		},
		{
			name: "no transactions found",
			filter: TransactionFilter{
				AccountID: strPtr("non_existent_account"),
				Cursor:    nil,
				Limit:     intPtr(5),
			},
			expectedResult: []*Transaction{},
			expectedCursor: "",
			wantErr:        false,
		},
	}

	repo := NewTransactionRepository(db, "txn_", "le_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, cursor, err := repo.ListTransactions(context.Background(), &tt.filter)

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

// Helper functions to create pointers
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}