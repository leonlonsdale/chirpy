-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
    VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2)
RETURNING
    *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = $1,
    hashed_password = $2,
    updated_at = NOW()
WHERE
    id = $3
RETURNING
    id,
    email,
    created_at,
    updated_at,
    is_chirpy_red;

-- name: UpgradeToChirpyRed :execrows
UPDATE
    users
SET
    is_chirpy_red = TRUE
WHERE
    id = $1;

