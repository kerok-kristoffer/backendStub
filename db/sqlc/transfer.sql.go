// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: transfer.sql

package db

import (
	"context"
	"database/sql"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
                       from_user_id,
                       to_user_id,
                       amount
) VALUES (
          $1, $2, $3
         ) RETURNING id, from_user_id, to_user_id, amount, created_at
`

type CreateTransferParams struct {
	FromUserID sql.NullInt64 `json:"fromUserID"`
	ToUserID   sql.NullInt64 `json:"toUserID"`
	Amount     sql.NullInt64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.FromUserID, arg.ToUserID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTransfer = `-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1
`

func (q *Queries) DeleteTransfer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransfer, id)
	return err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_user_id, to_user_id, amount, created_at FROM transfers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
