package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"aeshanw.com/accountApi/api/models"
	accountservice "aeshanw.com/accountApi/api/services/AccountService"
	"github.com/go-chi/render"
)

// AccountService defines the methods for interacting with the account service.
// type AccountService interface {
// 	// Define methods for interacting with the database
// 	CreateAccount(ctx context.Context, req models.CreateAccountRequest) (*AccountModel, error)
// }

// Handlers contains the HTTP handlers and dependencies.
type AccountHandler struct {
	db             *sql.DB
	accountservice accountservice.AccountServiceInt
}

// NewAccountHandler creates a new instance of Handlers with the provided dependencies.
func NewAccountHandler(db *sql.DB, as accountservice.AccountServiceInt) *AccountHandler {
	return &AccountHandler{
		db:             db,
		accountservice: as,
	}
}

func (ah *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, NewDefaultErrorResponse(ErrBadRequest))
		return
	}

	if errRes := ValidateCreateAccountRequest(req); errRes != nil {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, errRes)
		return
	}

	ctx := r.Context()

	//ServiceMethod to Validate & Save Account to DB
	err := ah.accountservice.CreateAccount(ctx, ah.db, req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, NewErrorResponse(ErrBadRequest, err.Error()))
		return
	}
	//TODO replace with logger
	fmt.Printf("account-model created: %d\n", req.AccountID)

	w.WriteHeader(http.StatusCreated) //Empty response is ok
}
