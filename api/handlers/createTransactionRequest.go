package handlers

import "aeshanw.com/accountApi/api/models"

func ValidateCreateTransactionRequest(req models.CreateTransactionRequest) *ErrorResponse {
	if req.SourceAccountID == 0 {
		return NewErrorResponse(ErrBadRequest, "invalid SourceAccountID")
	}
	if req.DestinationAccountID == 0 {
		return NewErrorResponse(ErrBadRequest, "invalid DestinationAccountID")
	}
	if req.Amount == "" {
		return NewErrorResponse(ErrBadRequest, "Amount is empty")
	}
	return nil
}
