-- +goose Up
CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;