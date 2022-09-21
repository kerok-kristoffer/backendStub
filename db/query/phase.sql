-- name: CreatePhase :one
INSERT INTO phases (
    name,
    description,
    formula_id
) VALUES (
             $1, $2, $3
         ) RETURNING *;

-- name: GetPhase :one
SELECT * FROM phases
WHERE id = $1 LIMIT 1;

-- name: ListPhasesByRecipeId :many
SELECT * FROM phases
WHERE formula_id = $1
ORDER BY id;

-- name: UpdatePhase :one
UPDATE phases
SET (name,
     description,
     formula_id
     ) =
        ($2, $3, $4)
WHERE id = $1
RETURNING *;

-- name: DeletePhase :exec
DELETE FROM phases
WHERE id = $1;