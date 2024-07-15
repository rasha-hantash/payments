package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/rasha-hantash/chariot-takehome/api/pkgs/test"
	"github.com/testcontainers/testcontainers-go"
	"log"
)

func TestCreateAccount(t *testing.T) {
	db, container := test.SetupAndFillDatabaseContainer("")
	defer func(container testcontainers.Container) {
		err := test.TeardownDatabaseContainer(container)
		if err != nil {
			log.Fatalf("failed to close container down: %v\n", err)
		}
	}(container)
	defer db.Close()

	tests := []struct {
		name        string
		input       *Account
		expectedErr error
	}{
		{
			name: "successful insert",
			input: &Account{
				AccountState: "active",
				AccountType:  "savings",
			},
			expectedErr: nil,
		},
	}

	repo := NewAccountRepository(db, "acct_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CreateAccount(context.Background(), tt.input)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestGetAccountBalance(t *testing.T) {
	db, container := test.SetupAndFillDatabaseContainer("seed_accounts_get_balance.sql")
	defer func(container testcontainers.Container) {
		err := test.TeardownDatabaseContainer(container)
		if err != nil {
			log.Fatalf("failed to close container down: %v\n", err)
		}
	}(container)
	defer db.Close()

	tests := []struct {
		name        string
		accountId   string
		expectedBal int64
		expectedErr bool
	}{
		{
			name:        "Valid internal balance",
			accountId:   "acct_1",
			expectedBal: 450,
			expectedErr: false,
		},
		{
			name:        "Valid external balance",
			accountId:   "acct_2",
			expectedBal: -450,
			expectedErr: false,
		},
		{
			name:        "No balance",
			accountId:   "account_3",
			expectedBal: 0,
			expectedErr: false,
		},
	}

	repo := &AccountRepository{db: db}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance, err := repo.GetAccountBalance(context.Background(), tt.accountId)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBal, balance)
			}
		})
	}
}
