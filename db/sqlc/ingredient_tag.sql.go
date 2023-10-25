// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: ingredient_tag.sql

package db

import (
	"context"
)

const createIngredientTag = `-- name: CreateIngredientTag :one
INSERT INTO ingredient_tags (
                             name, user_id
) VALUES ($1, $2) RETURNING id, name, user_id, created_at, updated_at
`

type CreateIngredientTagParams struct {
	Name   string `json:"name"`
	UserID int64  `json:"userID"`
}

func (q *Queries) CreateIngredientTag(ctx context.Context, arg CreateIngredientTagParams) (IngredientTag, error) {
	row := q.db.QueryRowContext(ctx, createIngredientTag, arg.Name, arg.UserID)
	var i IngredientTag
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteIngredientTag = `-- name: DeleteIngredientTag :exec
DELETE FROM ingredient_tags
WHERE id = $1
`

func (q *Queries) DeleteIngredientTag(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteIngredientTag, id)
	return err
}

const getIngredientTag = `-- name: GetIngredientTag :one
SELECT id, name, user_id, created_at, updated_at FROM ingredient_tags
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetIngredientTag(ctx context.Context, id int64) (IngredientTag, error) {
	row := q.db.QueryRowContext(ctx, getIngredientTag, id)
	var i IngredientTag
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getIngredientTagByName = `-- name: GetIngredientTagByName :one
SELECT id, name, user_id, created_at, updated_at FROM ingredient_tags
WHERE name ILIKE $1 AND user_id = $2 LIMIT 1
`

type GetIngredientTagByNameParams struct {
	Name   string `json:"name"`
	UserID int64  `json:"userID"`
}

func (q *Queries) GetIngredientTagByName(ctx context.Context, arg GetIngredientTagByNameParams) (IngredientTag, error) {
	row := q.db.QueryRowContext(ctx, getIngredientTagByName, arg.Name, arg.UserID)
	var i IngredientTag
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listIngredientTags = `-- name: ListIngredientTags :many
SELECT id, name, user_id, created_at, updated_at FROM ingredient_tags
WHERE user_id = $1
`

func (q *Queries) ListIngredientTags(ctx context.Context, userID int64) ([]IngredientTag, error) {
	rows, err := q.db.QueryContext(ctx, listIngredientTags, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []IngredientTag{}
	for rows.Next() {
		var i IngredientTag
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

const listIngredientTagsByIngredientId = `-- name: ListIngredientTagsByIngredientId :many
SELECT i."name" as "ingredient_name",
       i.id as "ingredient_id",
       it."name" as "ingredient_tag",
       it.id as "ingredient_tags_id",
       itm.id as "ingredient_tag_maps_id"
FROM "public".ingredient_tag_maps itm
         JOIN ingredients i on i.id = itm.ingredient_id
         JOIN ingredient_tags it on it.id = itm.ingredient_tag_id
WHERE i.id = $1
ORDER BY it.name
`

type ListIngredientTagsByIngredientIdRow struct {
	IngredientName      string `json:"ingredientName"`
	IngredientID        int64  `json:"ingredientID"`
	IngredientTag       string `json:"ingredientTag"`
	IngredientTagsID    int64  `json:"ingredientTagsID"`
	IngredientTagMapsID int64  `json:"ingredientTagMapsID"`
}

func (q *Queries) ListIngredientTagsByIngredientId(ctx context.Context, id int64) ([]ListIngredientTagsByIngredientIdRow, error) {
	rows, err := q.db.QueryContext(ctx, listIngredientTagsByIngredientId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListIngredientTagsByIngredientIdRow{}
	for rows.Next() {
		var i ListIngredientTagsByIngredientIdRow
		if err := rows.Scan(
			&i.IngredientName,
			&i.IngredientID,
			&i.IngredientTag,
			&i.IngredientTagsID,
			&i.IngredientTagMapsID,
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

const updateIngredientTag = `-- name: UpdateIngredientTag :one
UPDATE ingredient_tags
SET (name) =
        ($2)
WHERE id = $1
RETURNING id, name, user_id, created_at, updated_at
`

func (q *Queries) UpdateIngredientTag(ctx context.Context, id int64) (IngredientTag, error) {
	row := q.db.QueryRowContext(ctx, updateIngredientTag, id)
	var i IngredientTag
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
