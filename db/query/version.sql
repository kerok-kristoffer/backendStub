-- name: CreateVersion :one
INSERT INTO versions (
                      number
) VALUES (
          $1
         ) RETURNING *;

-- name: GetLatestVersion :one
SELECT * FROM versions
ORDER BY number
LIMIT 1;