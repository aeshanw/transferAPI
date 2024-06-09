package account_service

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"aeshanw.com/accountApi/api/models"
)

// MockDB is a mock database connection
type MockDB struct {
	mock.Mock
}

func (m *MockDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

// MockTx is a mock database transaction
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit() error {
	return m.Called().Error(0)
}

func (m *MockTx) Rollback() error {
	return m.Called().Error(0)
}

func (m *MockTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	argList := m.Called(ctx, query, args)
	return argList.Get(0).(*sql.Row)
}

func (m *MockTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	argList := m.Called(ctx, query, args)
	return argList.Get(0).(sql.Result), argList.Error(1)
}

// MockRow is a mock for sql.Row
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name                 string
		req                  models.CreateAccountRequest
		mockSetup            func(sqlmock.Sqlmock)
		expectError          bool
		expectedErrorMessage string
	}{
		{
			name: "successful account creation",
			req: models.CreateAccountRequest{
				AccountID:      1,
				InitialBalance: "100.0",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT (id) FROM accounts WHERE id=$1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count(id)"}).AddRow(0))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO accounts(id,balance) VALUES ($1,$2)")).WithArgs(1, 100.0).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError:          false,
			expectedErrorMessage: "",
		},
		{
			name: "account already exists",
			req: models.CreateAccountRequest{
				AccountID:      1,
				InitialBalance: "100.0",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT (id) FROM accounts WHERE id=$1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count(id)"}).AddRow(1))
				mock.ExpectRollback()
			},
			expectError:          true,
			expectedErrorMessage: "account already exists",
		},
		{
			name: "failure - database error",
			req: models.CreateAccountRequest{
				AccountID:      1,
				InitialBalance: "100.0",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT (id) FROM accounts WHERE id=$1")).WithArgs(1).WillReturnError(errors.New("database_error"))
				mock.ExpectRollback()
			},
			expectError:          true,
			expectedErrorMessage: "check for existing account:database_error",
		},
		{
			name: "initial balance is not a valid number",
			req: models.CreateAccountRequest{
				AccountID:      1,
				InitialBalance: "invalid",
			},
			mockSetup:            func(mock sqlmock.Sqlmock) {},
			expectError:          true,
			expectedErrorMessage: "invalid initial_balance format",
		},
		{
			name: "initial balance cannot be negative",
			req: models.CreateAccountRequest{
				AccountID:      1,
				InitialBalance: "-40.23",
			},
			mockSetup:            func(mock sqlmock.Sqlmock) {},
			expectError:          true,
			expectedErrorMessage: "invalid create-account-request due to:inital_balance cannot be less than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			// Setup mock expectations based on the test case
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			// Create a new account service with the mock database
			as := NewAccountService()

			// Invoke the CreateAccount method
			err = as.CreateAccount(context.Background(), db, tt.req)

			// Assert error if expected
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMessage)
			} else {
				assert.NoError(t, err)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
