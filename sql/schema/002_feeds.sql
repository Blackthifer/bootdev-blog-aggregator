-- +goose Up
CREATE TABLE feeds(
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    feed_name TEXT NOT NULL,
    feed_url TEXT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;