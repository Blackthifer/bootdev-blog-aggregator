-- name: CreateFeedFollow :one
WITH feed_follows_row AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $2, $3, $4)
    RETURNING *
) SELECT feed_follows_row.*, users.user_name, feeds.feed_name
FROM feed_follows_row
INNER JOIN users ON feed_follows_row.user_id = users.id
INNER JOIN feeds ON feed_follows_row.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, users.user_name, feeds.feed_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1 AND feed_id IN (
    SELECT id FROM feeds WHERE feed_url = $2
);
