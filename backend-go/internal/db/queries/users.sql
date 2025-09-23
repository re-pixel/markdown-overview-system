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

-- name: CreateSession :one
INSERT INTO user_sessions (user_id, session_token, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, session_token, created_at, expires_at;

-- name: GetSession :one
SELECT id, user_id, session_token, created_at, expires_at
FROM user_sessions
WHERE session_token = $1;

-- name: DeleteSession :exec
DELETE
FROM user_sessions
WHERE session_token = $1;