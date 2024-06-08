package handlers

import "aeshanw.com/accountApi/api/models"

// type CreateAccountRequest struct {
// 	AccountID      int64  `json:"account_id"`
// 	InitialBalance string `json:"initial_balance"`
// }

func ValidateCreateAccountRequest(req models.CreateAccountRequest) *ErrorResponse {
	if req.AccountID == 0 {
		return NewErrorResponse(ErrBadRequest, "invalid AccountID")
	}
	if req.InitialBalance == "" {
		return NewErrorResponse(ErrBadRequest, "InitialBalance is empty")
	}
	return nil
}

// {
// 	"account_id": 123,
// 	"initial_balance": "100.23344"
// 	}
