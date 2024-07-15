package repository

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/rasha-hantash/chariot-takehome/api/pkgs/test"
	"github.com/testcontainers/testcontainers-go"
)

func TestCreateUser(t *testing.T) {
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
		input       *User
		expectedErr bool
	}{
		{
			name: "successful insert",
			input: &User{
				Email:              "test@example.com",
				Name:               "Test User",
				IntLedgerAccountId: sql.NullString{String: "int-ledger-id", Valid: true},
				ExtLedgerAccountId: sql.NullString{String: "ext-ledger-id", Valid: true},
			},
			expectedErr: false,
		},
		{
			name: "insert error",
			input: &User{
				Email:              "error@example.com",
				Name:               "Error User",
				IntLedgerAccountId: sql.NullString{String: "int-ledger-id", Valid: true},
				ExtLedgerAccountId: sql.NullString{String: "ext-ledger-id", Valid: true},
			},
			expectedErr: true,
		},
	}

	accountRepo := NewAccountRepository(db, "acct_")
	userRepo := NewUserRepository(db, accountRepo, "usr_")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := userRepo.CreateUser(context.Background(), tt.input)

			// if res.Id != tt.expectedID {
			// 	t.Errorf("expected id %v, got %v", tt.expectedID, res.Id)
			// }

			if !tt.expectedErr {
				assert.NotEmpty(t,res.Id)
			}
			if err != nil && !tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

		})
	}
}
