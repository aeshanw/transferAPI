package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"aeshanw.com/accountApi/api/mocks"
	accountservice "aeshanw.com/accountApi/api/services/AccountService"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAccountDetails(t *testing.T) {
	mockDB := new(sql.DB)
	mockAccountService := new(mocks.MockAccountService)

	accountID := int64(1)
	accountModel := &accountservice.AccountModel{
		ID:      accountID,
		Balance: 100.23344,
	}

	mockAccountService.On("GetAccount", mock.Anything, mock.Anything, accountID).Return(accountModel, nil)

	accountHandler := &AccountHandler{
		db:             mockDB,
		accountservice: mockAccountService,
	}

	r := chi.NewRouter()
	r.Get("/accounts/{account_id}", accountHandler.GetAccountDetails)

	req, err := http.NewRequest("GET", "/accounts/"+strconv.FormatInt(accountID, 10), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expectedResponse := `{"account_id":1,"balance":"100.23344"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
}

func TestGetAccountDetails_InvalidAccountID(t *testing.T) {
	mockDB := new(sql.DB)
	mockAccountService := new(mocks.MockAccountService)

	accountID := int64(1)
	accountModel := &accountservice.AccountModel{
		ID:      accountID,
		Balance: 100.23344,
	}

	mockAccountService.On("GetAccount", mock.Anything, mock.Anything, accountID).Return(accountModel, nil)

	accountHandler := &AccountHandler{
		db:             mockDB,
		accountservice: mockAccountService,
	}

	r := chi.NewRouter()
	r.Get("/accounts/{account_id}", accountHandler.GetAccountDetails)

	req, err := http.NewRequest("GET", "/accounts/invalid", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "account_id parameter must be an integer")
}

func TestGetAccountDetails_AccountNotFound(t *testing.T) {
	mockDB := new(sql.DB)
	mockAccountService := new(mocks.MockAccountService)

	accountID := int64(1)
	mockAccountService.On("GetAccount", mock.Anything, mockDB, accountID).Return(nil, errors.New("account not found"))

	accountHandler := &AccountHandler{
		db:             mockDB,
		accountservice: mockAccountService,
	}

	r := chi.NewRouter()
	r.Get("/accounts/{account_id}", accountHandler.GetAccountDetails)

	req, err := http.NewRequest("GET", "/accounts/"+strconv.FormatInt(accountID, 10), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "account not found")
}
