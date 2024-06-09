package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	transactionservice "aeshanw.com/accountApi/api/services/TransactionService"

	"aeshanw.com/accountApi/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionService is a mock implementation of the TransactionServiceInt interface
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(ctx context.Context, db *sql.DB, req models.CreateTransactionRequest) (*transactionservice.TransactionModel, error) {
	args := m.Called(ctx, db, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactionservice.TransactionModel), args.Error(1)
}

func TestCreateTransaction(t *testing.T) {
	validTransactionModel := transactionservice.TransactionModel{
		ID:                   1,
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               100.50,
	}
	tests := []struct {
		name           string
		requestBody    models.CreateTransactionRequest
		mockSetup      func(m *MockTransactionService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful creation",
			requestBody: models.CreateTransactionRequest{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               "100.50",
			},
			mockSetup: func(m *MockTransactionService) {
				m.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).
					Return(&validTransactionModel, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"source_account_id":1,"destination_account_id":2,"amount":"100.50"}`,
		},
		{
			name: "invalid request body",
			requestBody: models.CreateTransactionRequest{
				SourceAccountID:      0,
				DestinationAccountID: 2,
				Amount:               "50.00",
			},
			mockSetup: func(m *MockTransactionService) {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"status\":400,\"detail\":\"bad_request\",\"message\":\"invalid SourceAccountID\"}", // Expected error message for invalid body
		},
		{
			name: "service error",
			requestBody: models.CreateTransactionRequest{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               "100.50",
			},
			mockSetup: func(m *MockTransactionService) {
				m.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("service error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"status\":400,\"detail\":\"bad_request\",\"message\":\"service error\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockDB := new(sql.DB)
			mockService := new(MockTransactionService)
			tt.mockSetup(mockService)

			// Create an instance of the handler
			handler := NewTransactionHandler(mockDB, mockService)

			// Create a request body
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
			assert.NoError(t, err)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create a handler function and serve the request
			http.HandlerFunc(handler.CreateTransaction).ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body
			assert.Equal(t, strings.TrimSpace(tt.expectedBody), strings.TrimSpace(rr.Body.String()))

			// Assert that the mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
