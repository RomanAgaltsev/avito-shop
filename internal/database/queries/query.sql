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

-- name: GetMerch :one
SELECT id, type, price
FROM merch
WHERE type = $1 LIMIT 1;

-- name: CreateHistoryRecord :one
INSERT INTO history (username, from_user, to_user, amount)
VALUES ($1, $2, $3, $4) RETURNING id;

-- name: CreateInventory :one
INSERT INTO inventory (username, type, quantity)
VALUES ($1, $2, 1) RETURNING id;

-- name: GetInventory :many
SELECT type, SUM(quantity) AS quantity
FROM inventory
WHERE username = $1
GROUP BY type;

-- name: GetHistory :many
SELECT from_user, to_user, SUM(amount) AS amount
FROM history
WHERE username = $1
GROUP BY from_user, to_user;