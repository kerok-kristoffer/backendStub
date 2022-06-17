-- name: CreateTransfer :one
INSERT INTO transfers (
                       from_user_id,
                       to_user_id,
                       amount
) VALUES (
          $1, $2, $3
         ) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;

