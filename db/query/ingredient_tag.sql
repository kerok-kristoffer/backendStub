-- name: CreateIngredientTag :one
INSERT INTO ingredient_tags (
                             user_id,
                             ingredient_id,
                             label
) VALUES (
             $1, $2, $3
         ) RETURNING *;

-- name: GetIngredientTag :one
SELECT * FROM ingredient_tags
WHERE id = $1 LIMIT 1;

-- name: ListIngredientTagsByUserId :many
SELECT * FROM ingredient_tags
WHERE user_id = $1
ORDER BY id
LIMIT $2
    OFFSET $3;

-- name: UpdateIngredientTag :one
UPDATE ingredient_tags
SET (user_id,
     ingredient_id,
     label) =
        ($2, $3, $4)
WHERE id = $1
RETURNING *;

-- name: DeleteIngredientTag :exec
DELETE FROM ingredient_tags
WHERE id = $1;