package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateAccount tests the CreateAccount handler
func TestCreateAccount(t *testing.T) {
	// Test cases
	tests := []struct {
		name               string
		inputTestPath      string
		expectedStatusCode int
		expectedResponse   *ErrorResponse
	}{
		{
			name:               "Invalid JSON",
			inputTestPath:      "testdata/invalidJSON.json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   NewDefaultErrorResponse(ErrBadRequest),
		},
		{
			name:               "Validation Error: Invalid AccountID",
			inputTestPath:      "testdata/invalidCreateAccountRequest_AccountID.json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid AccountID",
			},
		},
		{
			name:               "Validation Error: Invalid Balance",
			inputTestPath:      "testdata/invalidCreateAccountRequest_Balance.json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "InitialBalance is empty",
			},
		},
		{
			name:               "Successful Request",
			inputTestPath:      "testdata/validCreateAccountRequest.json",
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   nil,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the provided input

			// Read the JSON file
			file, err := os.Open(tc.inputTestPath)
			if err != nil {
				t.Fatalf("Failed to open file: %v", err)
			}
			defer file.Close()

			jsonData, err := io.ReadAll(file)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}
			req, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Call the handler
			handler := http.HandlerFunc(CreateAccount)
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tc.expectedStatusCode, rr.Code)

			// Check the response body if an error is expected
			if tc.expectedResponse != nil {
				var response ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				assert.Equal(t, tc.expectedResponse.Message, response.Message)
				assert.Equal(t, tc.expectedResponse.StatusCode, response.StatusCode)
			}
		})
	}
}
