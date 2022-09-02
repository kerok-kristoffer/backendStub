-- name: CreateCurrency :one
INSERT INTO currencies (
    name
) VALUES (
             $1
         ) RETURNING *;

-- name: GetCurrency :one
SELECT * FROM currencies
WHERE id = $1 LIMIT 1;

-- name: ListCurrencies :many
SELECT * FROM currencies
ORDER BY id
LIMIT $1
    OFFSET $2;

-- name: UpdateCurrency :one
UPDATE currencies
SET (name) =
        ($2)
WHERE id = $1
RETURNING *;

-- name: DeleteCurrency :exec
DELETE FROM currencies
WHERE id = $1;