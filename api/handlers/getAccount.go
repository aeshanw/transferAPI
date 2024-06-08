package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func GetAccountDetails(w http.ResponseWriter, r *http.Request) {
	// var req.URL
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	render.Status(r, http.StatusBadRequest)
	// 	render.Render(w, r, NewAccountResponse(BadRequest))
	// 	return
	// }

	//TODO ServiceMethod to Validate & Save Account to DB

	render.Status(r, http.StatusOK)
	//TODO populate AccountDetailsResponse
}
