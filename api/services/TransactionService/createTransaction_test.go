package transaction_service

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"aeshanw.com/accountApi/api/models"
)

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name                 string
		req                  models.CreateTransactionRequest
		expectedID           int64
		mockSetup            func(sqlmock.Sqlmock, models.CreateTransactionRequest, float64, int64)
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name: "successful transaction",
			req: models.CreateTransactionRequest{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               "100.50",
			},
			expectedID: 1,
			mockSetup: func(mock sqlmock.Sqlmock, req models.CreateTransactionRequest, amountFloat float64, expectedID int64) {
				// Expect BeginTx method to be called and return a transaction
				mock.ExpectBegin()

				// Expect QueryRowContext method to be called for checking accounts
				rows := sqlmock.NewRows([]string{"count"}).AddRow(2)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT (id) FROM accounts WHERE id IN ($1,$2)")).
					WithArgs(req.SourceAccountID, req.DestinationAccountID).
					WillReturnRows(rows)

				// Expect QueryRowContext method to be called for checking source account balance
				rows = sqlmock.NewRows([]string{"balance"}).AddRow(200.0)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM accounts WHERE id=$1")).
					WithArgs(req.SourceAccountID).
					WillReturnRows(rows)

				// Expect ExecContext method to be called for debiting source account balance
				mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = balance - $1 WHERE id=$2")).
					WithArgs(amountFloat, req.SourceAccountID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Expect ExecContext method to be called for crediting destination account balance
				mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = balance + $1 WHERE id=$2")).
					WithArgs(amountFloat, req.DestinationAccountID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Expect QueryRowContext method to be called for inserting new transaction
				rows = sqlmock.NewRows([]string{"id"}).AddRow(expectedID)
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO transactions(source_account_id,destination_account_id,amount) VALUES ($1,$2,$3) RETURNING id")).
					WithArgs(req.SourceAccountID, req.DestinationAccountID, amountFloat).
					WillReturnRows(rows)

				// Expect Commit method to be called
				mock.ExpectCommit()
			},
			expectError:          false,
			expectedErrorMessage: "",
		},
		{
			name: "failed transaction: insufficent source balance",
			req: models.CreateTransactionRequest{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               "100.50",
			},
			expectedID: 1,
			mockSetup: func(mock sqlmock.Sqlmock, req models.CreateTransactionRequest, amountFloat float64, expectedID int64) {
				// Expect BeginTx method to be called and return a transaction
				mock.ExpectBegin()

				// Expect QueryRowContext method to be called for checking accounts
				rows := sqlmock.NewRows([]string{"count"}).AddRow(2)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT (id) FROM accounts WHERE id IN ($1,$2)")).
					WithArgs(req.SourceAccountID, req.DestinationAccountID).
					WillReturnRows(rows)

				// Expect QueryRowContext method to be called for checking source account balance
				rows = sqlmock.NewRows([]string{"balance"}).AddRow(20.0)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM accounts WHERE id=$1")).
					WithArgs(req.SourceAccountID).
					WillReturnRows(rows)

				// Expect Commit method to be called
				mock.ExpectRollback()
			},
			expectError:          true,
			expectedErrorMessage: "source account has insufficent funds: finalSourceAccountBalance:",
		},
		{
			name: "failed transaction: missing accounts",
			req: models.CreateTransactionRequest{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               "100.50",
			},
			expectedID: 1,
			mockSetup: func(mock sqlmock.Sqlmock, req models.CreateTransactionRequest, amountFloat float64, expectedID int64) {
				// Expect BeginTx method to be called and return a transaction
				mock.ExpectBegin()

				// Expect QueryRowContext method to be called for checking accounts
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT (id) FROM accounts WHERE id IN ($1,$2)")).
					WithArgs(req.SourceAccountID, req.DestinationAccountID).
					WillReturnRows(rows)

				// Expect Commit method to be called
				mock.ExpectRollback()
			},
			expectError:          true,
			expectedErrorMessage: "account-count != 2 count:",
		},
		{
			name: "failed transaction: database error",
			req: models.CreateTransactionRequest{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               "100.50",
			},
			expectedID: 1,
			mockSetup: func(mock sqlmock.Sqlmock, req models.CreateTransactionRequest, amountFloat float64, expectedID int64) {
				// Expect BeginTx method to be called and return a transaction
				mock.ExpectBegin()

				// Expect QueryRowContext method to be called for checking accounts
				mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT (id) FROM accounts WHERE id IN ($1,$2)")).
					WithArgs(req.SourceAccountID, req.DestinationAccountID).
					WillReturnError(errors.New("database error"))

				// Expect Commit method to be called
				mock.ExpectRollback()
			},
			expectError:          true,
			expectedErrorMessage: "check for existing account:database error",
		},
		// Add more test cases here if needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock database: %v", err)
			}
			defer db.Close()

			// Create a new transaction service with the mock database
			transactionService := NewTransactionService()

			// Convert the Amount string to float64
			amountFloat, err := strconv.ParseFloat(tt.req.Amount, 64)
			if err != nil {
				assert.Fail(t, "unable to parse amount")
				return
			}

			// Setup mock expectations based on the test case
			if tt.mockSetup != nil {
				tt.mockSetup(mock, tt.req, amountFloat, tt.expectedID)
			}

			// Invoke the CreateTransaction method
			actualTransaction, err := transactionService.CreateTransaction(context.Background(), db, tt.req)

			// Assert error if expected
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMessage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, actualTransaction)
				assert.Equal(t, tt.expectedID, actualTransaction.ID)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
