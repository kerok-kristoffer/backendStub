-- name: CreateUser :one
INSERT INTO users (
                   full_name,
                   hash,
                   user_name,
                   email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUserName :one
SELECT * FROM users
WHERE user_name = $1 LIMIT 1;

-- name: GetUserByUserEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUserName :one
UPDATE users
SET full_name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserHash :one
UPDATE users
SET hash = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

