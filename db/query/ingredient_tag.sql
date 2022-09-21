-- name: CreateIngredientTag :one
INSERT INTO ingredient_tags (
                             name
) VALUES ($1) RETURNING *;

-- name: GetIngredientTag :one
SELECT * FROM ingredient_tags
WHERE id = $1 LIMIT 1;

-- name: ListIngredientTags :many
SELECT * FROM ingredient_tags;

-- name: UpdateIngredientTag :one
UPDATE ingredient_tags
SET (name) =
        ($2)
WHERE id = $1
RETURNING *;

-- name: DeleteIngredientTag :exec
DELETE FROM ingredient_tags
WHERE id = $1;