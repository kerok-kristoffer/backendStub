-- name: CreateIngredient :one
INSERT INTO ingredients (
                        name,
                        inci,
                        hash,
                        user_id,
                        function_id
) VALUES (
          $1, $2, $3, $4, $5
         ) RETURNING *;

-- name: GetIngredient :one
SELECT * FROM ingredients
WHERE id = $1 LIMIT 1;

-- name: GetIngredientCount :one
SELECT COUNT(*) FROM ingredients
WHERE user_id = $1;

-- name: ListIngredients :many
SELECT * FROM ingredients
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListIngredientsByUserId :many
SELECT * FROM ingredients
    WHERE user_id = $1
    ORDER BY name
    LIMIT $2
    OFFSET $3;

-- name: UpdateIngredient :one
UPDATE ingredients
SET (name,
     inci,
     hash,
     cost,
     user_id,
     function_id
    ) = (
                 $2, $3, $4, $5, $6, $7)
WHERE id = $1
RETURNING *;

-- name: DeleteIngredient :exec
DELETE FROM ingredients
WHERE id = $1;

-- name: DeleteIngredientsByUserId :exec
DELETE FROM ingredients
WHERE user_id = $1;
