// Package repository provides for interaction with DB.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RomanAgaltsev/avito-shop/internal/database/queries"
	"github.com/RomanAgaltsev/avito-shop/internal/model"
)

var (
	// ErrNoData error means that query has returned no data from DB.
	ErrNoData = fmt.Errorf("no data")

	// ErrConflict error means that there was a data conflict while execute a query.
	ErrConflict = fmt.Errorf("data conflict")

	// ErrNegativeBalance error means that user balance would become neagive if transaction being commited.
	ErrNegativeBalance = fmt.Errorf("negative balance")

	// DefaultBackOff - default backoff parameters.
	DefaultBackOff = NewDefaultBackOff()
)

// conflictUser contains confict user and an error.
type conflictUser struct {
	user model.User
	err  error
}

// PgxPool needs to mock pgxpool in tests.
type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (c *pgxpool.Conn, err error)
	AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	Reset()
	Config() *pgxpool.Config
	Stat() *pgxpool.Stat
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

// New creates new repository.
func New(dbpool PgxPool) (*Repository, error) {
	// Return Repository struct with new queries
	return &Repository{
		db: dbpool,
		q:  queries.New(dbpool),
	}, nil
}

// Repository is the repository structure.
type Repository struct {
	db PgxPool
	q  *queries.Queries
}

// CreateUser creates new user in the repository.
func (r *Repository) CreateUser(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) (model.User, error) {
	// PG error to catch the conflict
	var pgErr *pgconn.PgError

	// Create a function to wrap user creation with exponential backoff
	f := func() (conflictUser, error) {
		var cu conflictUser
		// Try to create new user
		_, errCreate := r.q.CreateUser(ctx, queries.CreateUserParams{
			Username: user.UserName,
			Password: user.Password,
		})

		// Check if there is a conflict
		if errors.As(errCreate, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			existingUser, errGet := backoff.RetryWithData(func() (queries.User, error) {
				// Return existing user
				return r.q.GetUser(ctx, user.UserName)
			}, bo)

			// Something has gone wrong
			if errGet != nil {
				return cu, errGet
			}

			return conflictUser{
				user: model.User{
					UserName: existingUser.Username,
					Password: existingUser.Password,
				},
				err: ErrConflict,
			}, nil
		}

		return cu, errCreate
	}

	// Call the wrapping function
	existingUser, err := backoff.RetryWithData(f, bo)
	if err != nil {
		return model.User{}, err
	}

	// There is a conflict
	if errors.Is(existingUser.err, ErrConflict) {
		return existingUser.user, existingUser.err
	}

	return user, nil
}

// CreateBalance creates new user balance in the repository.
func (r *Repository) CreateBalance(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) error {
	// Create new user balance in DB
	_, err := backoff.RetryWithData(func() (int32, error) {
		return r.q.CreateBalance(ctx, user.UserName)
	}, bo)

	if err != nil {
		return err
	}

	return nil
}

// SendCoins transfer given amount of coins from one user to another.
func (r *Repository) SendCoins(ctx context.Context, bo *backoff.ExponentialBackOff, fromUser model.User, toUser model.User, amount int) error {
	// Get user from DB
	_, err := backoff.RetryWithData(func() (queries.User, error) {
		return r.q.GetUser(ctx, toUser.UserName)
	}, bo)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoData
	}

	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Defer transaction rollback
	defer func() { _ = tx.Rollback(ctx) }()

	// Create query with transaction
	qtx := r.q.WithTx(tx)

	// Withdraw from the balance of user that sends
	fromUserBalance, err := backoff.RetryWithData(func() (int32, error) {
		return qtx.UpdateBalance(ctx, queries.UpdateBalanceParams{
			Username: fromUser.UserName,
			Coins:    int32(-amount),
		})
	}, bo)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// If the balance has become negative after withdrawal,
	// rollback the transaction and return the negative balance error
	if fromUserBalance < 0 {
		_ = tx.Rollback(ctx)
		return ErrNegativeBalance
	}

	// Balance enough to withdraw - create new history record for user that sends
	_, err = backoff.RetryWithData(func() (int32, error) {
		return qtx.CreateHistoryRecord(ctx, queries.CreateHistoryRecordParams{
			Username: fromUser.UserName,
			ToUser:   toUser.UserName,
			Amount:   int32(amount),
		})
	}, bo)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// Update the balance of user that receives
	_, err = backoff.RetryWithData(func() (int32, error) {
		return qtx.UpdateBalance(ctx, queries.UpdateBalanceParams{
			Username: toUser.UserName,
			Coins:    int32(amount),
		})
	}, bo)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// Create new history record for user that receives
	_, err = backoff.RetryWithData(func() (int32, error) {
		return qtx.CreateHistoryRecord(ctx, queries.CreateHistoryRecordParams{
			Username: toUser.UserName,
			FromUser: fromUser.UserName,
			Amount:   int32(amount),
		})
	}, bo)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// BuyItem register purhcase of inventory item (merch) for a given user.
func (r *Repository) BuyItem(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User, item model.InventoryItem) error {
	// Get merch from DB
	merch, err := backoff.RetryWithData(func() (queries.Merch, error) {
		return r.q.GetMerch(ctx, item.Type)
	}, bo)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoData
	}

	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Defer transaction rollback
	defer func() { _ = tx.Rollback(ctx) }()

	// Create query with transaction
	qtx := r.q.WithTx(tx)

	// Withdraw merch price from the balance of the user
	userBalance, err := backoff.RetryWithData(func() (int32, error) {
		return qtx.UpdateBalance(ctx, queries.UpdateBalanceParams{
			Username: user.UserName,
			Coins:    -merch.Price,
		})
	}, bo)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// If the balance has become negative after withdrawal,
	// rollback the transaction and return the negative balance error
	if userBalance < 0 {
		_ = tx.Rollback(ctx)
		return ErrNegativeBalance
	}

	// Balance enough to withdraw - add item to the user inventory
	_, err = backoff.RetryWithData(func() (int32, error) {
		return qtx.CreateInventory(ctx, queries.CreateInventoryParams{
			Username: user.UserName,
			Type:     item.Type,
		})
	}, bo)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// GetBalance returns users coins balance.
func (r *Repository) GetBalance(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) (int, error) {
	// Get user balance from DB
	userBalance, err := backoff.RetryWithData(func() (queries.Balance, error) {
		return r.q.GetBalance(ctx, user.UserName)
	}, bo)

	if err != nil {
		return 0, err
	}

	return int(userBalance.Coins), nil
}

// GetInventory returns users inventory.
func (r *Repository) GetInventory(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) ([]model.InventoryItem, error) {
	// Get user inventory from DB
	inventoryQuery, err := backoff.RetryWithData(func() ([]queries.GetInventoryRow, error) {
		return r.q.GetInventory(ctx, user.UserName)
	}, bo)

	if err != nil {
		return nil, err
	}

	// Fill the slice of inventory to return
	inventory := make([]model.InventoryItem, 0, len(inventoryQuery))
	for _, item := range inventoryQuery {
		inventory = append(inventory, model.InventoryItem{
			Type:     item.Type,
			Quantity: int(item.Quantity),
		})
	}

	return inventory, nil
}

// GetHistory returns users coins transaction history.
func (r *Repository) GetHistory(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) (model.CoinsHistory, error) {
	// Get history of user transactions from DB
	historyQuery, err := backoff.RetryWithData(func() ([]queries.GetHistoryRow, error) {
		return r.q.GetHistory(ctx, user.UserName)
	}, bo)

	if err != nil {
		return model.CoinsHistory{}, err
	}

	// There is a compromise between the number of database accesses and the memory allocations for slices capacity.
	received := make([]model.CoinsReceiving, 0, len(historyQuery))
	sent := make([]model.CoinsSending, 0, len(historyQuery))

	for _, rec := range historyQuery {
		if rec.FromUser != "" {
			received = append(received, model.CoinsReceiving{
				FromUser: rec.FromUser,
				Amount:   int(rec.Amount),
			})
			continue
		}
		if rec.ToUser != "" {
			sent = append(sent, model.CoinsSending{
				ToUser: rec.ToUser,
				Amount: int(rec.Amount),
			})
		}
	}

	coinHistory := model.CoinsHistory{
		Received: received,
		Sent:     sent,
	}

	return coinHistory, nil
}

// NewDefaultBackOff returns default backoff parameters.
func NewDefaultBackOff() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 50 * time.Millisecond
	bo.RandomizationFactor = 0.1
	bo.Multiplier = 2.0
	bo.MaxInterval = 1 * time.Second
	bo.MaxElapsedTime = 5 * time.Second
	bo.Reset()
	return bo
}
