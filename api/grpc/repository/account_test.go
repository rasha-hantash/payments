package repository

import (
	"context"
	"database/sql"
	"testing"
	"github.com/stretchr/testify/assert"

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

func TestGetAccountBalance(t *testing.T) {
	tests := []struct {
		name         string
		accountId    string
		mockQuery    func(mock sqlmock.Sqlmock, accountId string)
		expectedBal  int64
		expectedErr  bool
	}{
		{
			name:      "Valid balance",
			accountId: "account-1",
			mockQuery: func(mock sqlmock.Sqlmock, accountId string) {
				mock.ExpectQuery("SELECT SUM\\(CASE WHEN direction = 'credit' THEN amount ELSE -amount END\\) as balance FROM transactions WHERE account_id = \\?").
					WithArgs(accountId).
					WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(1000))
			},
			expectedBal: 1000,
			expectedErr: false,
		},
		{
			name:      "No transactions",
			accountId: "account-2",
			mockQuery: func(mock sqlmock.Sqlmock, accountId string) {
				mock.ExpectQuery("SELECT SUM\\(CASE WHEN direction = 'credit' THEN amount ELSE -amount END\\) as balance FROM transactions WHERE account_id = \\?").
					WithArgs(accountId).
					WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(0))
			},
			expectedBal: 0,
			expectedErr: false,
		},
		{
			name:      "Query error",
			accountId: "account-3",
			mockQuery: func(mock sqlmock.Sqlmock, accountId string) {
				mock.ExpectQuery("SELECT SUM\\(CASE WHEN direction = 'credit' THEN amount ELSE -amount END\\) as balance FROM transactions WHERE account_id = \\?").
					WithArgs(accountId).
					WillReturnError(sql.ErrNoRows)
			},
			expectedBal: 0,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := &AccountRepository{db: db}

			tt.mockQuery(mock, tt.accountId)

			balance, err := repo.GetAccountBalance(context.Background(), tt.accountId)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBal, balance)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

