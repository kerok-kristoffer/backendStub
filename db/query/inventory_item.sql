-- name: CreateInventoryItem :one
INSERT INTO inventory_items (
                            user_id,
                            ingredient_id,
                            amount_in_grams,
                            cost_per_gram,
                            currency_id,
                            expiry_date
) VALUES (
             $1, $2, $3, $4, $5, $6
         ) RETURNING *;

-- name: GetInventoryItem :one
SELECT * FROM inventory_items
WHERE id = $1 LIMIT 1;

-- name: ListInventoryItemsByUserId :many
SELECT * FROM inventory_items
WHERE user_id = $1
ORDER BY id
LIMIT $2
    OFFSET $3;

-- name: UpdateInventoryItem :one
UPDATE inventory_items
SET (user_id,
     ingredient_id,
     amount_in_grams,
     cost_per_gram,
     currency_id,
     expiry_date) =
        ($2, $3, $4, $5, $6, $7)
WHERE id = $1
RETURNING *;

-- name: DeleteInventoryItem :exec
DELETE FROM inventory_items
WHERE id = $1;