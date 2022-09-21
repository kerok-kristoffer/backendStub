-- name: CreateIngredientFunction :one
INSERT INTO ingredient_functions (
    name,
    user_id
) VALUES (
             $1, $2
         ) RETURNING *;

-- name: GetIngredientFunction :one
SELECT * FROM ingredient_functions
WHERE id = $1 LIMIT 1;

-- name: ListIngredientFunctionsByUserId :many
SELECT * FROM ingredient_functions
WHERE user_id = $1
ORDER BY id
LIMIT $2
    OFFSET $3;

-- name: UpdateIngredientFunction :one
UPDATE ingredient_functions
SET (name,
     user_id) =
        ($2, $3)
WHERE id = $1
RETURNING *;

-- name: DeleteIngredientFunction :exec
DELETE FROM ingredient_functions
WHERE id = $1;