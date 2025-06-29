-- +goose Up
CREATE TABLE posts(
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    post_url TEXT UNIQUE NOT NULL,
    post_description TEXT,
    published_at TIMESTAMP NOT NULL,
    feed_id INTEGER NOT NULL REFERENCES feeds ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;