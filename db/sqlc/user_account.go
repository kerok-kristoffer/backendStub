package db

import (
	"context"
	"database/sql"
	"fmt"
)

type UserAccount struct {
	// todo corresponds to store in tut
	*Queries
	db *sql.DB
}

func NewUserAccount(db *sql.DB) *UserAccount {
	return &UserAccount{
		db:      db,
		Queries: New(db),
	}
}

// DB transaction execution example, not in use currently
func (userAccount *UserAccount) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := userAccount.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromUserID int64 `json:"fromUserID"`
	ToUserID   int64 `json:"toUserID"`
	Amount     int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer   Transfer `json:"transfer"`
	FromUserID User     `json:"fromUserID"`
	ToUserID   User     `json:"toUserID"`
	FromEntry  Entry    `json:"fromEntry"`
	ToEntry    Entry    `json:"toEntry"`
}

func (userAccount *UserAccount) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	// TODO kerok - trying out transactions, no current need in project, keep for future reference example
	err := userAccount.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromUserID: sql.NullInt64{Int64: arg.FromUserID, Valid: true},
			ToUserID:   sql.NullInt64{Int64: arg.ToUserID, Valid: true},
			Amount:     sql.NullInt64{Int64: arg.Amount, Valid: true},
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			UserID: sql.NullInt64{Int64: arg.FromUserID, Valid: true},
			Amount: sql.NullInt64{Int64: arg.Amount, Valid: true},
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			UserID: sql.NullInt64{Int64: arg.ToUserID, Valid: true},
			Amount: sql.NullInt64{Int64: arg.Amount, Valid: true},
		})
		if err != nil {
			return err
		}
		// TODO update accounts - probably will skip this since I'm not really implementing transfers in this way.

		return nil
	})

	return result, err
}
