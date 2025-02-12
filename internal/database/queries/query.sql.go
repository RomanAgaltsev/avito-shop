// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package queries

import (
	"context"
)

const createBalance = `-- name: CreateBalance :one
INSERT INTO balance (username, coins)
VALUES ($1, 1000) RETURNING id
`

func (q *Queries) CreateBalance(ctx context.Context, username string) (int32, error) {
	row := q.db.QueryRow(ctx, createBalance, username)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createHistoryRecord = `-- name: CreateHistoryRecord :one
INSERT INTO history (username, from_user, to_user, amount)
VALUES ($1, $2, $3, $4) RETURNING id
`

type CreateHistoryRecordParams struct {
	Username string
	FromUser string
	ToUser   string
	Amount   int32
}

func (q *Queries) CreateHistoryRecord(ctx context.Context, arg CreateHistoryRecordParams) (int32, error) {
	row := q.db.QueryRow(ctx, createHistoryRecord,
		arg.Username,
		arg.FromUser,
		arg.ToUser,
		arg.Amount,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createInventory = `-- name: CreateInventory :one
INSERT INTO inventory (username, type, quantity)
VALUES ($1, $2, quantity+1) RETURNING id
`

type CreateInventoryParams struct {
	Username string
	Type     string
}

func (q *Queries) CreateInventory(ctx context.Context, arg CreateInventoryParams) (int32, error) {
	row := q.db.QueryRow(ctx, createInventory, arg.Username, arg.Type)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2) RETURNING id
`

type CreateUserParams struct {
	Username string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int32, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.Password)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getBalance = `-- name: GetBalance :one
SELECT id, username, coins
FROM balance
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetBalance(ctx context.Context, username string) (Balance, error) {
	row := q.db.QueryRow(ctx, getBalance, username)
	var i Balance
	err := row.Scan(&i.ID, &i.Username, &i.Coins)
	return i, err
}

const getInventory = `-- name: GetInventory :many
SELECT id, username, type, quantity, bought_at
FROM inventory
WHERE username = $1
ORDER BY bought_at
`

func (q *Queries) GetInventory(ctx context.Context, username string) ([]Inventory, error) {
	rows, err := q.db.Query(ctx, getInventory, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Inventory
	for rows.Next() {
		var i Inventory
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Type,
			&i.Quantity,
			&i.BoughtAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT id, username, password, created_at
FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const updateBalance = `-- name: UpdateBalance :one
UPDATE balance
SET coins = coins + $2
WHERE username = $1 RETURNING coins
`

type UpdateBalanceParams struct {
	Username string
	Coins    int32
}

func (q *Queries) UpdateBalance(ctx context.Context, arg UpdateBalanceParams) (int32, error) {
	row := q.db.QueryRow(ctx, updateBalance, arg.Username, arg.Coins)
	var coins int32
	err := row.Scan(&coins)
	return coins, err
}

const withdrawMerchFromBalance = `-- name: WithdrawMerchFromBalance :one
UPDATE balance
SET coins = coins - m.price
FROM merch AS m
LEFT JOIN balance AS b ON m.type = $2
WHERE b.username = $1 RETURNING b.coins
`

type WithdrawMerchFromBalanceParams struct {
	Username string
	Type     string
}

func (q *Queries) WithdrawMerchFromBalance(ctx context.Context, arg WithdrawMerchFromBalanceParams) (int32, error) {
	row := q.db.QueryRow(ctx, withdrawMerchFromBalance, arg.Username, arg.Type)
	var coins int32
	err := row.Scan(&coins)
	return coins, err
}
