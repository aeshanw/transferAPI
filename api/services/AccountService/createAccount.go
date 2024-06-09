package account_service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"aeshanw.com/accountApi/api/models"
)

// AccountService defines the methods for interacting with the account service.
type AccountServiceInt interface {
	// Define methods for interacting with the database
	CreateAccount(ctx context.Context, db *sql.DB, req models.CreateAccountRequest) error
	GetAccount(ctx context.Context, db *sql.DB, accountID int64) (*AccountModel, error)
}

type AccountModel struct {
	ID        int64
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccountModel() *AccountModel {
	return &AccountModel{}
}

func (am *AccountModel) SetFromRequest(req models.CreateAccountRequest) error {
	am.ID = req.AccountID

	floatIntialBalance, err := strconv.ParseFloat(req.InitialBalance, 64)
	if err != nil {
		return fmt.Errorf("invalid initial_balance format due to:%w", err)
	}

	if floatIntialBalance < 0 {
		return fmt.Errorf("inital_balance cannot be less than 0, input:%v", floatIntialBalance)
	}

	am.Balance = floatIntialBalance

	return nil
}

type AccountService struct{}

func NewAccountService() *AccountService {
	return &AccountService{}
}

// Mutex is required to handle race-conditions where 2 threads compete to UPDATE a account row in the DB
var mutex sync.Mutex

func (as *AccountService) CreateAccount(ctx context.Context, db *sql.DB, req models.CreateAccountRequest) error {
	//Mutex-lock to avoid race-cases
	mutex.Lock()
	defer mutex.Unlock()

	account := NewAccountModel()
	if err := account.SetFromRequest(req); err != nil {
		return fmt.Errorf("invalid create-account-request due to:%w", err)
	}

	fmt.Printf("account-model: %v\n", account)

	// Begin a transaction with the specified options
	txn, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("txn for createAccount fail:%w", err)
	}

	//Confirm the account exists
	sqlCheckForAccount := `SELECT COUNT (id) FROM accounts WHERE id=$1`
	sqlInsertNewAccount := `INSERT INTO accounts(id,balance) VALUES ($1,$2)`

	var count int
	if err := txn.QueryRowContext(ctx, sqlCheckForAccount, req.AccountID).Scan(&count); err != nil {
		txn.Rollback()
		return fmt.Errorf("check for existing account:%w", err)
	}

	log.Println("count done!")

	if count > 0 {
		//No existing account must exist
		txn.Rollback()
		return errors.New("account already exists")
	}

	log.Println("count check ok!")

	//Race conditions unlikely for this resource as the unique PK index ensures the 2nd try will fail hence data-consistency is maintained
	if _, err = txn.ExecContext(ctx, sqlInsertNewAccount, account.ID, account.Balance); err != nil {
		txn.Rollback()
		return fmt.Errorf("unable to insert new account due to :%w", err)
	}

	if err = txn.Commit(); err != nil {
		return fmt.Errorf("unable to commit account-creation txn due to :%w", err)
	}

	return nil
}
