-- name: CreateBatch :one
INSERT INTO batches (
                     user_id,
                     name,
                     description,
                     production_date,
                     expiry_date
) VALUES (
             $1, $2, $3, $4, $5
         ) RETURNING *;

-- name: GetBatch :one
SELECT * FROM batches
WHERE id = $1 LIMIT 1;

-- name: ListBatchesByUserId :many
SELECT * FROM batches
WHERE user_id = $1
ORDER BY id
LIMIT $2
    OFFSET $3;

-- name: UpdateBatch :one
UPDATE batches
SET (user_id,
     name,
     description,
     production_date,
     expiry_date) =
        ($2, $3, $4, $5, $6)
WHERE id = $1
RETURNING *;

-- name: DeleteBatch :exec
DELETE FROM batches
WHERE id = $1;