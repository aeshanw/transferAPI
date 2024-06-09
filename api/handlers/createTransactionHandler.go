package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"aeshanw.com/accountApi/api/models"

	transactionservice "aeshanw.com/accountApi/api/services/TransactionService"
	"github.com/go-chi/render"
)

type TransactionHandler struct {
	db                 *sql.DB
	transactionservice transactionservice.TransactionServiceInt
}

// NewTransactionHandler creates a new instance of Handlers with the provided dependencies.
func NewTransactionHandler(db *sql.DB, ts transactionservice.TransactionServiceInt) *TransactionHandler {
	return &TransactionHandler{
		db:                 db,
		transactionservice: ts,
	}
}

func (th *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, NewDefaultErrorResponse(ErrBadRequest))
		return
	}

	fmt.Printf("req: %v\n", req)

	if errRes := ValidateCreateTransactionRequest(req); errRes != nil {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, errRes)
		return
	}

	_, err := th.transactionservice.CreateTransaction(r.Context(), th.db, req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, NewErrorResponse(ErrBadRequest, err.Error()))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, req)
}
