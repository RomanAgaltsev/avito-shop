package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RomanAgaltsev/avito-shop/internal/database/queries"
	"github.com/RomanAgaltsev/avito-shop/internal/model"
)

var (
	ErrConflict        = fmt.Errorf("data conflict")
	ErrNegativeBalance = fmt.Errorf("negative balance")

	DefaultBackOff = backoff.NewExponentialBackOff()
)

// conflictUser contains confict user and an error.
type conflictUser struct {
	user model.User
	err  error
}

// New creates new repository.
func New(dbpool *pgxpool.Pool) (*Repository, error) {
	// Return Repository struct with new queries
	return &Repository{
		db: dbpool,
		q:  queries.New(dbpool),
	}, nil
}

// Repository is the repository structure.
type Repository struct {
	db *pgxpool.Pool
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
		return qtx.WithdrawMerchFromBalance(ctx, queries.WithdrawMerchFromBalanceParams{
			Username: user.UserName,
			Type:     item.Type,
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

func (r *Repository) GetInventory(ctx context.Context, user model.User) ([]model.InventoryItem, error) {
	return nil, nil
}

func (r *Repository) GetHistory(ctx context.Context, user model.User) (model.CoinsHistory, error) {
	return model.CoinsHistory{}, nil
}
