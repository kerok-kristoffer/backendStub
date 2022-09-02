-- name: CreateRecipeIngredient :one
INSERT INTO recipe_ingredients (
                                ingredient_id,
                                percentage,
                                phase_id,
                                description
) VALUES (
          $1, $2, $3, $4
         ) RETURNING *;

-- name: GetRecipeIngredient :one
SELECT * FROM recipe_ingredients
    WHERE id = $1 LIMIT 1;

-- name: ListRecipeIngredientsByUserId :many
SELECT * FROM recipe_ingredients
    WHERE phase_id = $1
    ORDER BY id;

-- name: UpdateRecipeIngredient :one
UPDATE recipe_ingredients
SET (ingredient_id,
     percentage,
     phase_id,
    description) =
        ($2, $3, $4, $5)
WHERE id = $1
RETURNING *;

-- name: DeleteRecipeIngredient :exec
DELETE FROM recipe_ingredients
WHERE id = $1;