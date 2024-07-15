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
		name             string
		amount           float64
		userId           string
		debitAccountId   string
		creditAccountId  string
		expectedIdLength int
		wantErr          bool
	}{
		{
			name:             "successful deposit",
			amount:           100,
			userId:           "usr_1",
			debitAccountId:   "acct_1",
			creditAccountId:  "acct_2",
			expectedIdLength: 20,
			wantErr:          false,
		},
		{
			name:             "insert error,  debit account does not exist",
			amount:           100,
			userId:           "usr_2",
			debitAccountId:   "acct_3",
			creditAccountId:  "acct_4",
			expectedIdLength: 0,
			wantErr:          true,
		},
		{
			name:             "insert error, credit account does not exist",
			amount:           100,
			userId:           "usr_3",
			debitAccountId:   "acct_5",
			creditAccountId:  "acct_6",
			expectedIdLength: 0,
			wantErr:          true,
		},
	}

	repo := NewTransactionRepository(db, "txn_", "le_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txnId, err := repo.DepositFunds(context.Background(), tt.amount, tt.userId, tt.debitAccountId, tt.creditAccountId)

			if len(txnId) != int(tt.expectedIdLength) {
				t.Errorf("expected result %d, got %d", tt.expectedIdLength, len(txnId))
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIdLength, len(txnId))
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
		name             string
		amount           float64
		userId           string
		debitAccountId   string
		creditAccountId  string
		expectedIdLength int
		wantErr          bool
	}{
		{
			name:             "successful withdraw",
			amount:           100,
			userId:           "usr_1",
			debitAccountId:   "acct_1",
			creditAccountId:  "acct_2",
			expectedIdLength: 20,
			wantErr:          false,
		},
		{
			name:             "insufficient funds",
			amount:           100,
			userId:           "usr_2",
			debitAccountId:   "acct_3",
			creditAccountId:  "acct_4",
			expectedIdLength: 0,
			wantErr:          true,
		},
		{
			name:             "credit account does not exist",
			amount:           100,
			userId:           "usr_3",
			debitAccountId:   "acct_5",
			creditAccountId:  "acct_6",
			expectedIdLength: 0,
			wantErr:          true,
		},
	}

	repo := NewTransactionRepository(db, "txn_", "le_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txnId, err := repo.WithdrawFunds(context.Background(), tt.amount, tt.userId, tt.debitAccountId, tt.creditAccountId)

			if len(txnId) != tt.expectedIdLength {
				t.Errorf("expected result %v, got %v", tt.expectedIdLength, len(txnId))
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIdLength, len(txnId))
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
		name             string
		amount           float64
		userId           string
		debitAccountId   string
		creditAccountId  string
		expectedIdLength int
		wantErr          bool
	}{
		{
			name:             "successful transfer",
			amount:           100,
			userId:           "usr_1",
			debitAccountId:   "acct_1",
			creditAccountId:  "acct_2",
			expectedIdLength: 20,
			wantErr:          false,
		},
		{
			name:             "insert error debit insufficient funds",
			amount:           100,
			userId:           "usr_2",
			debitAccountId:   "acct_3",
			creditAccountId:  "acct_4",
			expectedIdLength: 0,
			wantErr:          true,
		},
		{
			name:             "insert error debit insufficient funds",
			amount:           100,
			userId:           "usr_1",
			debitAccountId:   "acct_5",
			creditAccountId:  "acct_6",
			expectedIdLength: 0,
			wantErr:          true,
		},
	}

	repo := NewTransactionRepository(db, "txn_", "le_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txnId, err := repo.TransferFunds(context.Background(), tt.amount, tt.userId, tt.debitAccountId, tt.creditAccountId)

			if len(txnId) != tt.expectedIdLength {
				t.Errorf("expected result %v, got %v", tt.expectedIdLength, len(txnId))
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIdLength, len(txnId))
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
		expectedResult []Transaction
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
			expectedResult: []Transaction{
				{Id: "txn_11", Amount: 210, Status: "success"},
				{Id: "txn_12", Amount: 220, Status: "success"},
				{Id: "txn_14", Amount: 240, Status: "success"},
				{Id: "txn_15", Amount: 250, Status: "success"},
				{Id: "txn_17", Amount: 270, Status: "success"},
			},
			expectedCursor: "txn_17",
			wantErr:        false,
		},
		{
			name: "successful list for acct_2",
			filter: TransactionFilter{
				AccountID: strPtr("acct_2"),
				Cursor:    nil,
				Limit:     intPtr(5),
			},
			expectedResult: []Transaction{
				{Id: "txn_1", Amount: 110, Status: "success"},
				{Id: "txn_10", Amount: 200, Status: "success"},
				{Id: "txn_12", Amount: 220, Status: "success"},
				{Id: "txn_13", Amount: 230, Status: "success"},
				{Id: "txn_15", Amount: 250, Status: "success"},
			},
			expectedCursor: "txn_15",
			wantErr:        false,
		},
		{
			name: "successful list for acct_3",
			filter: TransactionFilter{
				AccountID: strPtr("acct_3"),
				Cursor:    nil,
				Limit:     intPtr(5),
			},
			expectedResult: []Transaction{
				{Id: "txn_1", Amount: 110, Status: "success"},
				{Id: "txn_10", Amount: 200, Status: "success"},
				{Id: "txn_11", Amount: 210, Status: "success"},
				{Id: "txn_13", Amount: 230, Status: "success"},
				{Id: "txn_14", Amount: 240, Status: "success"},
			},
			expectedCursor: "txn_14",
			wantErr:        false,
		},
		{
			name: "successful list with cursor",
			filter: TransactionFilter{
				AccountID: strPtr("acct_1"),
				Cursor:    nil,
				Limit:     intPtr(3),
			},
			expectedResult: []Transaction{
				{Id: "txn_11", Amount: 210, Status: "success"},
				{Id: "txn_12", Amount: 220, Status: "success"},
				{Id: "txn_14", Amount: 240, Status: "success"},
			},
			expectedCursor: "txn_14",
			wantErr:        false,
		},
		{
			name: "no transactions found",
			filter: TransactionFilter{
				AccountID: strPtr("non_existent_account"),
				Cursor:    nil,
				Limit:     intPtr(5),
			},
			expectedResult: []Transaction{},
			expectedCursor: "",
			wantErr:        false,
		},
	}

	repo := NewTransactionRepository(db, "txn_", "le_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, cursor, err := repo.ListTransactions(context.Background(), &tt.filter)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				for i, result := range results {
					assert.Equal(t, tt.expectedResult[i].Id, result.Id)
					assert.Equal(t, tt.expectedResult[i].Amount, result.Amount)
					assert.Equal(t, tt.expectedResult[i].Status, result.Status)
				}
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
