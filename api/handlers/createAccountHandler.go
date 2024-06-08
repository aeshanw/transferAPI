package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, NewDefaultErrorResponse(ErrBadRequest))
		return
	}

	if errRes := req.Validate(); errRes != nil {
		render.Status(r, http.StatusBadRequest)
		render.Render(w, r, errRes)
		return
	}

	fmt.Printf("req: %v\n", req)

	//TODO ServiceMethod to Validate & Save Account to DB

	w.WriteHeader(http.StatusCreated) //Empty response is ok
}
