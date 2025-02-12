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

-- name: WithdrawMerchFromBalance :one
UPDATE balance
SET coins = coins - m.price
FROM merch AS m
LEFT JOIN balance AS b ON m.type = $2
WHERE b.username = $1 RETURNING b.coins;

-- name: CreateHistoryRecord :one
INSERT INTO history (username, from_user, to_user, amount)
VALUES ($1, $2, $3, $4) RETURNING id;

-- name: CreateInventory :one
INSERT INTO inventory (username, type, quantity)
VALUES ($1, $2, quantity+1) RETURNING id;

-- name: GetInventory :many
SELECT id, username, type, quantity, bought_at
FROM inventory
WHERE username = $1
ORDER BY bought_at;