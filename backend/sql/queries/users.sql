-- name: CreateUser :exec
INSERT INTO users (id, discord_user_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),  -- Generates a new UUID
    $1,
    NOW(),              -- Sets created_at to the current timestamp
    NOW()              -- Sets updated_at to the current timestamp                  
);

-- name: GetUser :one
SELECT * FROM users where discord_user_id=$1;