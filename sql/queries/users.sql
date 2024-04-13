-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, api_key)
VALUES ($1, $2, $3, $4, encode(digest(random()::text, 'sha256'), 'hex'))
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE api_key = $1;