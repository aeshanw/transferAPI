package account_service

import (
	"context"
	"database/sql"
	"fmt"
)

func (as *AccountService) GetAccount(ctx context.Context, db *sql.DB, accountID int64) (*AccountModel, error) {
	sqlGetAccount := `SELECT id,balance,created_at,updated_at FROM accounts WHERE id=$1`

	var account AccountModel
	err := db.QueryRowContext(ctx, sqlGetAccount, accountID).Scan(&account.ID, &account.Balance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to fetch account due to: %w", err)
		}
		return nil, err
	}

	fmt.Printf("account: %v\n", account)

	return &account, nil
}
