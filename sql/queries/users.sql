-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, user_name)
VALUES( $1, $2, $2, $3)
RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users
WHERE user_name = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;