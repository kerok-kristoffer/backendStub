-- name: CreateBatchItem :one
INSERT INTO batch_items (
                         amount,
                         inventory_item_id,
                         recipe_ingredient_id,
                         batch_id,
                         user_id,
                         description

) VALUES (
             $1, $2, $3, $4, $5, $6
         ) RETURNING *;

-- name: GetBatchItem :one
SELECT * FROM batch_items
WHERE id = $1 LIMIT 1;

-- name: ListBatchItemsByUserId :many
SELECT * FROM batch_items
WHERE batch_id = $1
ORDER BY id;

-- name: UpdateBatchItem :one
UPDATE batch_items
SET (amount,
     inventory_item_id,
     recipe_ingredient_id,
     batch_id,
     user_id,
     description) =
        ($2, $3, $4, $5, $6, $7)
WHERE id = $1
RETURNING *;

-- name: DeleteBatchItem :exec
DELETE FROM batch_items
WHERE id = $1;