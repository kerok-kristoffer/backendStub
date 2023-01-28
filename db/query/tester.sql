-- name: CreateTester :one
INSERT INTO testers (
    email
) VALUES (
             $1
         ) RETURNING *;

-- name: GetTesterByEmail :one
SELECT * FROM testers
WHERE email = $1 LIMIT 1;

-- name: UpdateTester :one
UPDATE testers
SET (
     email,
     user_id,
     updated_at) =
        ($2, $3, CURRENT_TIMESTAMP)
WHERE id = $1
RETURNING *;
