package models

import "net/http"

type CreateAccountRequest struct {
	AccountID      int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

type CreateTransactionRequest struct {
	SourceAccountID      int64  `json:"source_account_id"`
	DestinationAccountID int64  `json:"destination_account_id"`
	Amount               string `json:"amount"`
}

func (ctr CreateTransactionRequest) Render(w http.ResponseWriter, r *http.Request) error {
	// TODO Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// {
// 	"source_account_id": 123,
// 	"destination_account_id": 456,
// 	"amount": "100.12345"
// 	}
