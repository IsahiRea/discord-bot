-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    discord_user_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;