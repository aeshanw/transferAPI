package transaction_service

import (
	"context"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"aeshanw.com/accountApi/api/models"
)

func TestCreateTransaction(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new transaction service with the mock database
	transactionService := NewTransactionService()

	// Define the input and expected output
	req := models.CreateTransactionRequest{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               "100.50",
	}
	expectedID := int64(1)

	// Convert the Amount string to float64
	amountFloat, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		// Handle error
		assert.Fail(t, "unable to parse amount")
	}

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
		WillReturnRows(rows)

	// Expect Commit method to be called
	mock.ExpectCommit()

	// Invoke the CreateTransaction method
	actualTransaction, err := transactionService.CreateTransaction(context.Background(), db, req)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the returned transaction matches the expected transaction
	assert.NotNil(t, actualTransaction)
	assert.Equal(t, expectedID, actualTransaction.ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
