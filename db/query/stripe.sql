-- name: CreateStripeEntry :one
INSERT INTO stripe (
    id,
    user_id,
    stripe_customer_id,
    stripe_plan_id
) VALUES (
          $1, $2, $3, $4
         ) RETURNING *;

-- name: GetStripeByUserId :one
SELECT * FROM stripe
WHERE user_id = $1 LIMIT 1;

-- name: UpdateStripeByUserId :one
UPDATE stripe
SET (stripe_customer_id,
     stripe_plan_id,
     updated_at) =
    ($2, $3, CURRENT_TIMESTAMP)
WHERE user_id = $1
RETURNING *;