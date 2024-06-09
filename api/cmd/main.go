package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"aeshanw.com/accountApi/api/handlers"
	accountservice "aeshanw.com/accountApi/api/services/AccountService"
	transactionservice "aeshanw.com/accountApi/api/services/TransactionService"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL is empty")
	}

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	as := accountservice.NewAccountService()
	accHandler := handlers.NewAccountHandler(db, as)

	ts := transactionservice.NewTransactionService()
	trHandler := handlers.NewTransactionHandler(db, ts)

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
		r.Post("/", accHandler.CreateAccount)                // POST /accounts
		r.Get("/{account_id}", accHandler.GetAccountDetails) // GET /accounts/{account_id}
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", trHandler.CreateTransaction) // POST /transactions
	})

	log.Println("API running at :3000 port")
	http.ListenAndServe(":3000", r)
}
