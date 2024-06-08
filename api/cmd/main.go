package main

import (
	"log"
	"net/http"

	"aeshanw.com/accountApi/api/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome")) //TODO consider removing this unused route after testing
	})
	// RESTy routes for "accounts" resource
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", handlers.CreateAccount)                // POST /accounts
		r.Get("/{account_id}", handlers.GetAccountDetails) // GET /accounts/{account_id}
	})

	log.Println("API running at :3000 port")
	http.ListenAndServe(":3000", r)
}
