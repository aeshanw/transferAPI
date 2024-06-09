package account_service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetAccount(t *testing.T) {
	tests := []struct {
		name         string
		accountID    int64
		mockSetup    func(sqlmock.Sqlmock)
		expectedErr  error
		expectedAcct *AccountModel
	}{
		{
			name:      "successfully retrieve account",
			accountID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
					AddRow(1, 100.23, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT id,balance,created_at,updated_at FROM accounts WHERE id=\$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedErr: nil,
			expectedAcct: &AccountModel{
				ID:      1,
				Balance: 100.23,
			},
		},
		{
			name:      "account not found",
			accountID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id,balance,created_at,updated_at FROM accounts WHERE id=\$1`).
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectedErr:  fmt.Errorf("unable to fetch account due to: %w", sql.ErrNoRows),
			expectedAcct: nil,
		},
		{
			name:      "database error",
			accountID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id,balance,created_at,updated_at FROM accounts WHERE id=\$1`).
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			expectedErr:  errors.New("database error"),
			expectedAcct: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			as := NewAccountService()
			account, err := as.GetAccount(context.Background(), db, tt.accountID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAcct.ID, account.ID)
				assert.Equal(t, tt.expectedAcct.Balance, account.Balance)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
