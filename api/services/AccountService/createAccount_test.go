package account_service

import (
	"context"
	"database/sql"
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
	ctx := context.Background()

	// Mock the DB and Tx
	mockDB := new(MockDB)
	mockTx := new(MockTx)
	mockRow := new(MockRow)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT (id) FROM accounts WHERE id=$1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count(id)"}).AddRow(0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO accounts(id,balance) VALUES ($1,$2)")).WithArgs(1, 100.0).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Prepare the request
	req := models.CreateAccountRequest{
		AccountID:      1,
		InitialBalance: "100.0",
	}

	// Call the CreateAccount function
	err = CreateAccount(ctx, db, req)

	// Assertions
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}
