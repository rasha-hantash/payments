package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name         string
		input        *Account
		mockSetup    func(mock sqlmock.Sqlmock)
		expectedID   string
		expectedErr  error
	}{
		{
			name: "successful insert",
			input: &Account{
				Id:           "acc-123",
				AccountState: "active",
				AccountType:  "savings",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO accounts`).
					WithArgs("acc-123", "active", "savings").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("acc-123"))
			},
			expectedID:  "acc-123",
			expectedErr: nil,
		},
		{
			name: "insert error",
			input: &Account{
				Id:           "acc-456",
				AccountState: "inactive",
				AccountType:  "checking",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO accounts`).
					WithArgs("acc-456", "inactive", "checking").
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

			repo := NewAccountRepository(db, "acct")
			id, err := repo.CreateAccount(context.Background(), tt.input)

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
