// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: formula.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createFormula = `-- name: CreateFormula :one
INSERT INTO formulas (
                     name,
                     default_amount,
                     default_amount_oz,
                     description,
                     user_id
) VALUES (
          $1, $2, $3, $4, $5
         ) RETURNING id, name, default_amount, default_amount_oz, description, user_id, created_at, updated_at
`

type CreateFormulaParams struct {
	Name            string  `json:"name"`
	DefaultAmount   float32 `json:"defaultAmount"`
	DefaultAmountOz float32 `json:"defaultAmountOz"`
	Description     string  `json:"description"`
	UserID          int64   `json:"userID"`
}

func (q *Queries) CreateFormula(ctx context.Context, arg CreateFormulaParams) (Formula, error) {
	row := q.db.QueryRowContext(ctx, createFormula,
		arg.Name,
		arg.DefaultAmount,
		arg.DefaultAmountOz,
		arg.Description,
		arg.UserID,
	)
	var i Formula
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.DefaultAmount,
		&i.DefaultAmountOz,
		&i.Description,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFormula = `-- name: DeleteFormula :exec
DELETE FROM formulas
WHERE id = $1
`

func (q *Queries) DeleteFormula(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteFormula, id)
	return err
}

const getFormula = `-- name: GetFormula :one
SELECT id, name, default_amount, default_amount_oz, description, user_id, created_at, updated_at FROM formulas
    WHERE id = $1 LIMIT 1
`

func (q *Queries) GetFormula(ctx context.Context, id int64) (Formula, error) {
	row := q.db.QueryRowContext(ctx, getFormula, id)
	var i Formula
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.DefaultAmount,
		&i.DefaultAmountOz,
		&i.Description,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFullFormula = `-- name: GetFullFormula :many
SELECT f.id as "formula_id", f.name as "formula_name", f.default_amount, f.default_amount_oz,
       f.description, f.user_id, f.created_at, f.updated_at,
       p.name as "phase_name", p.description as "phase_description",
       fi.phase_id, fi.id as "formula_ingredient_id", fi.ingredient_id, fi.cost, fi.percentage,
       i.name as "ingredient_name", i.inci, i.function_id
FROM formulas f
         JOIN phases p on f.id = p.formula_id
         JOIN formula_ingredients fi on p.id = fi.phase_id
         JOIN ingredients i on fi.ingredient_id = i.id
WHERE f.id = $1 ORDER BY p.id, fi.phase_id
`

type GetFullFormulaRow struct {
	FormulaID           int64           `json:"formulaID"`
	FormulaName         string          `json:"formulaName"`
	DefaultAmount       float32         `json:"defaultAmount"`
	DefaultAmountOz     float32         `json:"defaultAmountOz"`
	Description         string          `json:"description"`
	UserID              int64           `json:"userID"`
	CreatedAt           time.Time       `json:"createdAt"`
	UpdatedAt           time.Time       `json:"updatedAt"`
	PhaseName           string          `json:"phaseName"`
	PhaseDescription    string          `json:"phaseDescription"`
	PhaseID             int64           `json:"phaseID"`
	FormulaIngredientID int64           `json:"formulaIngredientID"`
	IngredientID        int64           `json:"ingredientID"`
	Cost                sql.NullFloat64 `json:"cost"`
	Percentage          float32         `json:"percentage"`
	IngredientName      string          `json:"ingredientName"`
	Inci                string          `json:"inci"`
	FunctionID          sql.NullInt64   `json:"functionID"`
}

func (q *Queries) GetFullFormula(ctx context.Context, id int64) ([]GetFullFormulaRow, error) {
	rows, err := q.db.QueryContext(ctx, getFullFormula, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFullFormulaRow{}
	for rows.Next() {
		var i GetFullFormulaRow
		if err := rows.Scan(
			&i.FormulaID,
			&i.FormulaName,
			&i.DefaultAmount,
			&i.DefaultAmountOz,
			&i.Description,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PhaseName,
			&i.PhaseDescription,
			&i.PhaseID,
			&i.FormulaIngredientID,
			&i.IngredientID,
			&i.Cost,
			&i.Percentage,
			&i.IngredientName,
			&i.Inci,
			&i.FunctionID,
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

const listFormulasByUserId = `-- name: ListFormulasByUserId :many
SELECT id, name, default_amount, default_amount_oz, description, user_id, created_at, updated_at FROM formulas
    WHERE user_id = $1
    ORDER BY name
    LIMIT $2
    OFFSET $3
`

type ListFormulasByUserIdParams struct {
	UserID int64 `json:"userID"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListFormulasByUserId(ctx context.Context, arg ListFormulasByUserIdParams) ([]Formula, error) {
	rows, err := q.db.QueryContext(ctx, listFormulasByUserId, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Formula{}
	for rows.Next() {
		var i Formula
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.DefaultAmount,
			&i.DefaultAmountOz,
			&i.Description,
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

const updateFormula = `-- name: UpdateFormula :one
UPDATE formulas
    SET (name,
        default_amount,
        default_amount_oz,
        description,
        user_id,
        updated_at) =
        ($2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
    WHERE id = $1
    RETURNING id, name, default_amount, default_amount_oz, description, user_id, created_at, updated_at
`

type UpdateFormulaParams struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	DefaultAmount   float32 `json:"defaultAmount"`
	DefaultAmountOz float32 `json:"defaultAmountOz"`
	Description     string  `json:"description"`
	UserID          int64   `json:"userID"`
}

func (q *Queries) UpdateFormula(ctx context.Context, arg UpdateFormulaParams) (Formula, error) {
	row := q.db.QueryRowContext(ctx, updateFormula,
		arg.ID,
		arg.Name,
		arg.DefaultAmount,
		arg.DefaultAmountOz,
		arg.Description,
		arg.UserID,
	)
	var i Formula
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.DefaultAmount,
		&i.DefaultAmountOz,
		&i.Description,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
