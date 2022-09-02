-- name: CreateRecipe :one
INSERT INTO recipes (
                     name,
                     default_amount,
                     description,
                     user_id
) VALUES (
          $1, $2, $3, $4
         ) RETURNING *;

-- name: GetRecipe :one
SELECT * FROM recipes
    WHERE id = $1 LIMIT 1;

-- name: ListRecipesByUserId :many
SELECT * FROM recipes
    WHERE user_id = $1
    ORDER BY id
    LIMIT $2
    OFFSET $3;

-- name: UpdateRecipe :one
UPDATE recipes
    SET (name,
        default_amount,
        description,
        user_id) =
        ($2, $3, $4, $5)
    WHERE id = $1
    RETURNING *;

-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE id = $1;