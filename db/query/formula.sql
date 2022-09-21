-- name: CreateFormula :one
INSERT INTO formulas (
                     name,
                     default_amount,
                     default_amount_oz,
                     description,
                     user_id
) VALUES (
          $1, $2, $3, $4, $5
         ) RETURNING *;

-- name: GetFormula :one
SELECT * FROM formulas
    WHERE id = $1 LIMIT 1;

-- name: GetFullFormula :many
SELECT r.id as "formula_id", r.name as "formula_name", r.default_amount, r.default_amount_oz,
       r.description, r.user_id,
       p.name as "phase_name", p.description as "phase_description",
       ri.phase_id, ri.id as "formula_ingredient_id", ri.ingredient_id, ri.percentage,
       i.name as "ingredient_name", i.inci, i.function_id
FROM formulas r
         JOIN phases p on r.id = p.formula_id
         JOIN formula_ingredients ri on p.id = ri.phase_id
         JOIN ingredients i on ri.ingredient_id = i.id
WHERE r.id = $1;

-- name: ListFormulasByUserId :many
SELECT * FROM formulas
    WHERE user_id = $1
    ORDER BY id
    LIMIT $2
    OFFSET $3;

-- name: UpdateFormula :one
UPDATE formulas
    SET (name,
        default_amount,
        description,
        user_id) =
        ($2, $3, $4, $5)
    WHERE id = $1
    RETURNING *;

-- name: DeleteFormula :exec
DELETE FROM formulas
WHERE id = $1;