package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	accountservice "aeshanw.com/accountApi/api/services/AccountService"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type GetAccountDetailsResponse struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"balance"`
}

func (gadr *GetAccountDetailsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// TODO Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewGetAccountDetailsResponse(am *accountservice.AccountModel) (*GetAccountDetailsResponse, error) {
	if am == nil {
		return nil, errors.New("accountModel is nil")
	}

	formattedBalance := fmt.Sprintf("%.5f", am.Balance)

	return &GetAccountDetailsResponse{
		AccountID: am.ID,
		Balance:   formattedBalance,
	}, nil
}

func (ah *AccountHandler) GetAccountDetails(w http.ResponseWriter, r *http.Request) {
	//Validate input
	accountIDStr := chi.URLParam(r, "account_id")
	if accountIDStr == "" {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, NewErrorResponse(ErrBadRequest, "account_id query parameter is required"))
		return
	}

	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, NewErrorResponse(ErrBadRequest, "account_id parameter must be an integer"))
		return
	}

	//ServiceMethod to Validate & Get AccountDetails from DB
	accountModel, err := ah.accountservice.GetAccount(r.Context(), ah.db, int64(accountID))
	if err != nil {
		fmt.Printf("TEST err:%v\n", err)
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, NewErrorResponse(ErrBadRequest, err.Error()))
		return
	}

	resp, err := NewGetAccountDetailsResponse(accountModel)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.Render(w, r, NewErrorResponse(ErrInternalServerError, err.Error()))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, resp)
}
