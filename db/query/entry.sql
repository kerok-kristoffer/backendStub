-- name: CreateEntry :one
INSERT INTO entries (
                     user_id,
                     amount
) VALUES (
          $1, $2
         ) RETURNING *;

-- name: GetEntryByUserId :one
SELECT * FROM entries
WHERE user_id = $1 LIMIT 1;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;