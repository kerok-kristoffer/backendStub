-- name: CreateFormulaIngredient :one
INSERT INTO formula_ingredients (
                                ingredient_id,
                                percentage,
                                phase_id,
                                cost,
                                description,
                                update_id
) VALUES (
          $1, $2, $3, $4, $5, $6
         ) RETURNING *;

-- name: GetFormulaIngredient :one
SELECT * FROM formula_ingredients
    WHERE id = $1 LIMIT 1;

-- name: ListFormulaIngredientsByPhaseId :many
SELECT * FROM formula_ingredients
    WHERE phase_id = $1
    ORDER BY id;

-- name: UpdateFormulaIngredient :one
UPDATE formula_ingredients
SET (ingredient_id,
     percentage,
     phase_id,
     cost,
     description,
     update_id,
     updated_at) =
        ($2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
WHERE id = $1
RETURNING *;

-- name: DeleteFormulaIngredient :exec
DELETE FROM formula_ingredients
WHERE id = $1;

-- name: DeleteFormulaIngredientsNotInUpdate :exec
DELETE FROM formula_ingredients
WHERE phase_id IN ($1)
AND update_id NOT IN ($2);
