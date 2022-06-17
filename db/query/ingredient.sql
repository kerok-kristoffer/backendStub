-- name: CreateIngredient :one
INSERT INTO ingredients (
                        name,
                        hash,
                        user_id
) VALUES (
          $1, $2, $3
         ) RETURNING *;

-- name: GetIngredient :one
SELECT * FROM ingredients
WHERE id = $1 LIMIT 1;

-- name: ListIngredients :many
SELECT * FROM ingredients
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListIngredientsByUserId :many
SELECT * FROM ingredients
    WHERE user_id = $1
    ORDER BY id
    LIMIT $2
    OFFSET $3;

-- name: UpdateIngredientName :one
UPDATE ingredients
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteIngredient :exec
DELETE FROM ingredients
WHERE id = $1;

-- name: DeleteIngredientsByUserId :exec
DELETE FROM ingredients
WHERE user_id = $1;

/*
 "id" bigserial PRIMARY KEY,
                               "name" varchar NOT NULL,
                               "user_id" bigint,
                               "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
 */