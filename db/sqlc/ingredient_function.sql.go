// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: ingredient_function.sql

package db

import (
	"context"
)

const createIngredientFunction = `-- name: CreateIngredientFunction :one
INSERT INTO ingredient_functions (
    name,
    user_id
) VALUES (
             $1, $2
         ) RETURNING id, name, user_id, created_at, updated_at
`

type CreateIngredientFunctionParams struct {
	Name   string `json:"name"`
	UserID int64  `json:"userID"`
}

func (q *Queries) CreateIngredientFunction(ctx context.Context, arg CreateIngredientFunctionParams) (IngredientFunction, error) {
	row := q.db.QueryRowContext(ctx, createIngredientFunction, arg.Name, arg.UserID)
	var i IngredientFunction
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteIngredientFunction = `-- name: DeleteIngredientFunction :exec
DELETE FROM ingredient_functions
WHERE id = $1
`

func (q *Queries) DeleteIngredientFunction(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteIngredientFunction, id)
	return err
}

const getIngredientFunction = `-- name: GetIngredientFunction :one
SELECT id, name, user_id, created_at, updated_at FROM ingredient_functions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetIngredientFunction(ctx context.Context, id int64) (IngredientFunction, error) {
	row := q.db.QueryRowContext(ctx, getIngredientFunction, id)
	var i IngredientFunction
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listIngredientFunctionsByUserId = `-- name: ListIngredientFunctionsByUserId :many
SELECT id, name, user_id, created_at, updated_at FROM ingredient_functions
WHERE user_id = $1
ORDER BY id
LIMIT $2
    OFFSET $3
`

type ListIngredientFunctionsByUserIdParams struct {
	UserID int64 `json:"userID"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListIngredientFunctionsByUserId(ctx context.Context, arg ListIngredientFunctionsByUserIdParams) ([]IngredientFunction, error) {
	rows, err := q.db.QueryContext(ctx, listIngredientFunctionsByUserId, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []IngredientFunction{}
	for rows.Next() {
		var i IngredientFunction
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateIngredientFunction = `-- name: UpdateIngredientFunction :one
UPDATE ingredient_functions
SET (name,
     user_id) =
        ($2, $3)
WHERE id = $1
RETURNING id, name, user_id, created_at, updated_at
`

type UpdateIngredientFunctionParams struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	UserID int64  `json:"userID"`
}

func (q *Queries) UpdateIngredientFunction(ctx context.Context, arg UpdateIngredientFunctionParams) (IngredientFunction, error) {
	row := q.db.QueryRowContext(ctx, updateIngredientFunction, arg.ID, arg.Name, arg.UserID)
	var i IngredientFunction
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}