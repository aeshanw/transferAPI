package mocks

import (
	"context"
	"database/sql"

	"aeshanw.com/accountApi/api/models"
	accountservice "aeshanw.com/accountApi/api/services/AccountService"
	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) GetAccount(ctx context.Context, db *sql.DB, accountID int64) (*accountservice.AccountModel, error) {
	args := m.Called(ctx, db, accountID)
	if args.Get(0) != nil {
		return args.Get(0).(*accountservice.AccountModel), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAccountService) CreateAccount(ctx context.Context, db *sql.DB, req models.CreateAccountRequest) error {
	args := m.Called(ctx, db, req)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return args.Error(0)
}
