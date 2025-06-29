-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, feed_name, feed_url, user_id)
VALUES ($1, $2, $2, $3, $4, $5)
RETURNING *;