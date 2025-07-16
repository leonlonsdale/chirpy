-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    created_at timestamp,
    updated_at timestamp,
    email text UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;

