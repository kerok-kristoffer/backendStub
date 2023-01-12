-- name: CreateStripePlan :one
INSERT INTO stripe_plans (
                          id,
                          name,
                          user_access_id
) VALUES (
          $1, $2, $3
         ) RETURNING *;

-- name: ListStripePlans :many
SELECT * FROM stripe_plans
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetStripePlanByUserAccess :one
SELECT * FROM stripe_plans
WHERE user_access_id = $1
LIMIT 1;