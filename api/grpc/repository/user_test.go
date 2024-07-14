package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name         string
		input        *User
		mockSetup    func(mock sqlmock.Sqlmock)
		expectedID   string
		expectedErr  error
	}{
		{
			name: "successful insert",
			input: &User{
				Email:              "test@example.com",
				Name:               "Test User",
				IntLedgerAccountId: sql.NullString{String: "int-ledger-id", Valid: true},
				ExtLedgerAccountId: sql.NullString{String: "ext-ledger-id", Valid: true},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs("test@example.com", "Test User", sql.NullString{String: "int-ledger-id", Valid: true}, sql.NullString{String: "ext-ledger-id", Valid: true}).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("12345"))
			},
			expectedID:  "12345",
			expectedErr: nil,
		},
		{
			name: "insert error",
			input: &User{
				Email:              "error@example.com",
				Name:               "Error User",
				IntLedgerAccountId: sql.NullString{String: "int-ledger-id", Valid: true},
				ExtLedgerAccountId: sql.NullString{String: "ext-ledger-id", Valid: true},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs("error@example.com", "Error User", sql.NullString{String: "int-ledger-id", Valid: true}, sql.NullString{String: "ext-ledger-id", Valid: true}).
					WillReturnError(sql.ErrConnDone)
			},
			expectedID:  "",
			expectedErr: sql.ErrConnDone,
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

			accountRepo := NewAccountRepository(db, "acct_")
			userRepo := NewUserRepository(db, accountRepo, "usr_")
			
			id, err := userRepo.CreateUser(context.Background(), tt.input)

			if id != tt.expectedID {
				t.Errorf("expected id %v, got %v", tt.expectedID, id)
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
