-- name: CreateUser :one
INSERT INTO users (username, email, pass)
VALUES ($1, $2, $3)
RETURNING id, username, email, pass, created_at;

-- name: GetUserByEmail :one
SELECT id, username, email, pass, created_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, username, email, pass, created_at
FROM users
ORDER BY created_at DESC;

-- name: GetUserIdByUsername :one
SELECT id
FROM users
WHERE username = $1;