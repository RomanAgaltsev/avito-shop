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
