package handlers

import "aeshanw.com/accountApi/api/models"

func ValidateCreateAccountRequest(req models.CreateAccountRequest) *ErrorResponse {
	if req.AccountID == 0 {
		return NewErrorResponse(ErrBadRequest, "invalid AccountID")
	}
	if req.InitialBalance == "" {
		return NewErrorResponse(ErrBadRequest, "InitialBalance is empty")
	}
	return nil
}
