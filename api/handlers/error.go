package handlers

import (
	"net/http"
)

// ErrorResponse: Generic ErrorResponse
type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Error      string `json:"detail"`
	Message    string `json:"message,omitempty"`
}

// NewDefaultErrorResponse: for default errors that need no override of the message
func NewDefaultErrorResponse(err ErrorResponse) *ErrorResponse {
	return &err
}

// NewErrorResponse: customized error messages
func NewErrorResponse(err ErrorResponse, msg string) *ErrorResponse {
	if msg != "" {
		err.Message = msg
	}
	return &err
}

func (re *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// TODO Pre-processing before a response is marshalled and sent across the wire
	return nil
}

var (
	ErrBadRequest = ErrorResponse{
		StatusCode: http.StatusBadRequest,
		Error:      "bad_request",
		Message:    "The request could not be understood or was missing required parameters.",
	}
	ErrUnauthorized = ErrorResponse{
		StatusCode: http.StatusUnauthorized,
		Error:      "unauthorized",
		Message:    "Authentication failed or user does not have permissions for the requested operation.",
	}
	ErrNotFound = ErrorResponse{
		StatusCode: http.StatusNotFound,
		Error:      "not_found",
		Message:    "The requested resource could not be found.",
	}
	ErrInternalServerError = ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Error:      "internal_server_error",
		Message:    "An unexpected error occurred on the server.",
	}
	ErrServiceUnavailable = ErrorResponse{
		StatusCode: http.StatusServiceUnavailable,
		Error:      "service_unavailable",
		Message:    "The service is temporarily unavailable. Please try again later.",
	}
)
