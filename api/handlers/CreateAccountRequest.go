package handlers

type CreateAccountRequest struct {
	AccountID      int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

func (car *CreateAccountRequest) Validate() *ErrorResponse {
	if car.AccountID == 0 {
		return NewErrorResponse(ErrBadRequest, "invalid AccountID")
	}
	if car.InitialBalance == "" {
		return NewErrorResponse(ErrBadRequest, "InitialBalance is empty")
	}
	return nil
}

// {
// 	"account_id": 123,
// 	"initial_balance": "100.23344"
// 	}
