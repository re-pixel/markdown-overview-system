-- name: CreateUser :one
INSERT INTO users (email, pass)
VALUES ($1, $2)
RETURNING id, email, pass, created_at;

-- name: GetUserByEmail :one
SELECT id, email, pass, created_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, email, pass, created_at
FROM users
ORDER BY created_at DESC;
