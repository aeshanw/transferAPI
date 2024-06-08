package handlers

type CreateAccountRequest struct {
	AccountID      int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

// {
// 	"account_id": 123,
// 	"initial_balance": "100.23344"
// 	}
