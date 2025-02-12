-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2) RETURNING id;

-- name: GetUser :one
SELECT id, username, password, created_at
FROM users
WHERE username = $1 LIMIT 1;

-- name: CreateBalance :one
INSERT INTO balance (username, coins)
VALUES ($1, 1000) RETURNING id;

-- name: GetBalance :one
SELECT id, username, coins
FROM balance
WHERE username = $1 LIMIT 1;

-- name: UpdateBalance :one
UPDATE balance
SET coins = coins + $2
WHERE username = $1 RETURNING coins;

-- name: CreateHistoryRecord :one
INSERT INTO history (username, from_user, to_user, amount)
VALUES ($1, $2, $3, $4) RETURNING id;