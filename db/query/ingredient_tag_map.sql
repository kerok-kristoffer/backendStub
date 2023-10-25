-- name: CreateIngredientTagMap :one
INSERT INTO ingredient_tag_maps (
    ingredient_tag_id, ingredient_id
) VALUES ($1, $2) RETURNING *;

-- name: GetIngredientTagMapByIngredientTagId :one
SELECT * FROM ingredient_tag_maps
WHERE ingredient_tag_id = $1 LIMIT 1;

-- name: GetIngredientTagMapByIngredientId :one
SELECT * FROM ingredient_tag_maps
WHERE ingredient_id = $1 LIMIT 1;

-- name: GetTagMap :one
SELECT * FROM ingredient_tag_maps
WHERE ingredient_id = $1 AND ingredient_tag_id = $2
LIMIT 1;

-- name: UpdateIngredientTagMap :one
UPDATE ingredient_tag_maps
SET (ingredient_id, ingredient_tag_id) =
        ($2, $3)
WHERE id = $1
RETURNING *;

-- name: DeleteIngredientTagMap :exec
DELETE FROM ingredient_tag_maps
WHERE id = $1;

-- name: DeleteIngredientTagMapByIngredientId :exec
DELETE FROM ingredient_tag_maps
WHERE ingredient_id = $1;