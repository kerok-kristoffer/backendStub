// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: phase.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createPhase = `-- name: CreatePhase :one
INSERT INTO phases (
    name,
    description,
    formula_id,
    update_id
) VALUES (
             $1, $2, $3, $4
         ) RETURNING id, name, description, formula_id, update_id, created_at, updated_at
`

type CreatePhaseParams struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FormulaID   int64     `json:"formulaID"`
	UpdateID    uuid.UUID `json:"updateID"`
}

func (q *Queries) CreatePhase(ctx context.Context, arg CreatePhaseParams) (Phase, error) {
	row := q.db.QueryRowContext(ctx, createPhase,
		arg.Name,
		arg.Description,
		arg.FormulaID,
		arg.UpdateID,
	)
	var i Phase
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.FormulaID,
		&i.UpdateID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePhase = `-- name: DeletePhase :exec
DELETE FROM phases
WHERE id = $1
`

func (q *Queries) DeletePhase(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePhase, id)
	return err
}

const getPhase = `-- name: GetPhase :one
SELECT id, name, description, formula_id, update_id, created_at, updated_at FROM phases
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPhase(ctx context.Context, id int64) (Phase, error) {
	row := q.db.QueryRowContext(ctx, getPhase, id)
	var i Phase
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.FormulaID,
		&i.UpdateID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPhasesNotInUpdate = `-- name: GetPhasesNotInUpdate :many
SELECT id, name, description, formula_id, update_id, created_at, updated_at FROM phases
WHERE phases.formula_id IN ($1)
  AND update_id NOT IN ($2)
`

type GetPhasesNotInUpdateParams struct {
	FormulaID int64     `json:"formulaID"`
	UpdateID  uuid.UUID `json:"updateID"`
}

func (q *Queries) GetPhasesNotInUpdate(ctx context.Context, arg GetPhasesNotInUpdateParams) ([]Phase, error) {
	rows, err := q.db.QueryContext(ctx, getPhasesNotInUpdate, arg.FormulaID, arg.UpdateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Phase{}
	for rows.Next() {
		var i Phase
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.FormulaID,
			&i.UpdateID,
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

const listPhasesByFormulaId = `-- name: ListPhasesByFormulaId :many
SELECT id, name, description, formula_id, update_id, created_at, updated_at FROM phases
WHERE formula_id = $1
ORDER BY created_at
`

func (q *Queries) ListPhasesByFormulaId(ctx context.Context, formulaID int64) ([]Phase, error) {
	rows, err := q.db.QueryContext(ctx, listPhasesByFormulaId, formulaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Phase{}
	for rows.Next() {
		var i Phase
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.FormulaID,
			&i.UpdateID,
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

const updatePhase = `-- name: UpdatePhase :one
UPDATE phases
SET (name,
     description,
     formula_id,
     update_id,
     updated_at
     ) =
        ($2, $3, $4, $5, CURRENT_TIMESTAMP)
WHERE id = $1
RETURNING id, name, description, formula_id, update_id, created_at, updated_at
`

type UpdatePhaseParams struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FormulaID   int64     `json:"formulaID"`
	UpdateID    uuid.UUID `json:"updateID"`
}

func (q *Queries) UpdatePhase(ctx context.Context, arg UpdatePhaseParams) (Phase, error) {
	row := q.db.QueryRowContext(ctx, updatePhase,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.FormulaID,
		arg.UpdateID,
	)
	var i Phase
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.FormulaID,
		&i.UpdateID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
