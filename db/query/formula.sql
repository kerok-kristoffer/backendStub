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
SELECT f.id as "formula_id", f.name as "formula_name", f.default_amount, f.default_amount_oz,
       f.description, f.user_id,
       p.name as "phase_name", p.description as "phase_description",
       fi.phase_id, fi.id as "formula_ingredient_id", fi.ingredient_id, fi.cost, fi.percentage,
       i.name as "ingredient_name", i.inci, i.function_id
FROM formulas f
         JOIN phases p on f.id = p.formula_id
         JOIN formula_ingredients fi on p.id = fi.phase_id
         JOIN ingredients i on fi.ingredient_id = i.id
WHERE f.id = $1 ORDER BY fi.phase_id;

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