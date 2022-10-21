-- name: CreatePhase :one
INSERT INTO phases (
    name,
    description,
    formula_id,
    update_id
) VALUES (
             $1, $2, $3, $4
         ) RETURNING *;

-- name: GetPhase :one
SELECT * FROM phases
WHERE id = $1 LIMIT 1;

-- name: ListPhasesByFormulaId :many
SELECT * FROM phases
WHERE formula_id = $1
ORDER BY id;

-- name: UpdatePhase :one
UPDATE phases
SET (name,
     description,
     formula_id,
     update_id
     ) =
        ($2, $3, $4, $5)
WHERE id = $1
RETURNING *;

-- name: GetPhasesNotInUpdate :many
SELECT * FROM phases
WHERE phases.formula_id IN ($1)
  AND update_id NOT IN ($2);

-- name: DeletePhase :exec
DELETE FROM phases
WHERE id = $1;