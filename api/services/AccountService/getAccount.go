package account_service

import (
	"context"
	"database/sql"
)

func GetAccount(ctx context.Context, db *sql.DB, accountID int64) (*AccountModel, error) {
	return nil, nil
}
