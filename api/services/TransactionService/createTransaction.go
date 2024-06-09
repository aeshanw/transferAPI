package transaction_service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"aeshanw.com/accountApi/api/models"
)

// TransactionServiceInt defines the methods for interacting with the account service.
type TransactionServiceInt interface {
	// Define methods for interacting with the database
	CreateTransaction(ctx context.Context, db *sql.DB, req models.CreateTransactionRequest) (*TransactionModel, error)
}

type TransactionModel struct {
	ID                   int64
	SourceAccountID      int64
	DestinationAccountID int64
	Amount               float64
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func NewTransactionModel() *TransactionModel {
	return &TransactionModel{}
}

func (tm *TransactionModel) Render(w http.ResponseWriter, r *http.Request) error {
	// TODO Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func (tm *TransactionModel) SetFromRequest(req models.CreateTransactionRequest) error {
	tm.SourceAccountID = req.SourceAccountID
	tm.DestinationAccountID = req.DestinationAccountID

	if tm.SourceAccountID == tm.DestinationAccountID {
		return fmt.Errorf("sourceAccountID and destinationAccountID cannot be the same")
	}

	floatAmount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return err
	}

	if floatAmount < 0 {
		return fmt.Errorf("inital_balance cannot be less than 0, input:%v", floatAmount)
	}

	tm.Amount = floatAmount

	return nil
}

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

// Mutex is required to handle race-conditions where 2 threads compete to UPDATE a account row in the DB
var mutex sync.Mutex

func (ts *TransactionService) CreateTransaction(ctx context.Context, db *sql.DB, req models.CreateTransactionRequest) (*TransactionModel, error) {
	//Mutex-lock to avoid race-cases
	mutex.Lock()
	defer mutex.Unlock()

	transaction := NewTransactionModel()
	if err := transaction.SetFromRequest(req); err != nil {
		return nil, fmt.Errorf("invalid create-transaction-request due to:%w", err)
	}

	fmt.Printf("transaction-model: %v\n", transaction)

	// Begin a transaction with the specified options
	txn, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("txn for createTransaction fail:%w", err)
	}

	//Confirm the account exists
	sqlCheckForAccounts := `SELECT COUNT (id) FROM accounts WHERE id IN ($1,$2)`
	sqlCheckSourceBalance := `SELECT balance FROM accounts WHERE id=$1`
	sqlDebitSourceAccountBalance := `UPDATE accounts SET balance = balance - $1 WHERE id=$2`
	sqlCreditDestinationAccountBalance := `UPDATE accounts SET balance = balance + $1 WHERE id=$2`
	sqlInsertNewTransaction := `INSERT INTO transactions(source_account_id,destination_account_id,amount) VALUES ($1,$2,$3) RETURNING id`

	var count int
	if err := txn.QueryRow(sqlCheckForAccounts, req.SourceAccountID, req.DestinationAccountID).Scan(&count); err != nil {
		txn.Rollback()
		return nil, fmt.Errorf("check for existing account:%w", err)
	}

	log.Println("count done!")

	if count != 2 {
		//Both accounts must exist
		txn.Rollback()
		return nil, fmt.Errorf("account-count != 2 count:%d", count)
	}

	log.Println("count check ok!")

	//Check SourceBalance
	var sourceAccountBalance float64
	if err := txn.QueryRow(sqlCheckSourceBalance, req.SourceAccountID).Scan(&sourceAccountBalance); err != nil {
		txn.Rollback()
		return nil, fmt.Errorf("check for source account balance:%w", err)
	}

	finalSourceAccountBalance := sourceAccountBalance - transaction.Amount

	if finalSourceAccountBalance < 0 {
		//balance cannot fall below 0
		txn.Rollback()
		return nil, fmt.Errorf("source account has insufficent funds: finalSourceAccountBalance:%v", finalSourceAccountBalance)
	}

	//Debit Source
	if _, err = txn.Exec(sqlDebitSourceAccountBalance, transaction.Amount, transaction.SourceAccountID); err != nil {
		txn.Rollback()
		return nil, fmt.Errorf("unable to debit source account due to :%w", err)
	}

	//Credit Destination
	if _, err = txn.Exec(sqlCreditDestinationAccountBalance, transaction.Amount, transaction.DestinationAccountID); err != nil {
		txn.Rollback()
		return nil, fmt.Errorf("unable to credit destination account due to :%w", err)
	}

	//No other issues can proceed to lock-in the transaction
	if err = txn.QueryRow(sqlInsertNewTransaction, transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount).Scan(&transaction.ID); err != nil {
		txn.Rollback()
		return nil, fmt.Errorf("unable to insert new account due to :%w", err)
	}

	if err = txn.Commit(); err != nil {
		txn.Rollback()
		return nil, fmt.Errorf("unable to commit account-creation txn due to :%w", err)
	}

	return transaction, nil
}
