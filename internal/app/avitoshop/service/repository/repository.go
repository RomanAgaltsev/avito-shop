package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RomanAgaltsev/avito-shop/internal/model"
)

var (
	ErrConflict        = fmt.Errorf("data conflict")
	ErrNegativeBalance = fmt.Errorf("negative balance")
)

// New creates new repository.
func New(dbpool *pgxpool.Pool) (*Repository, error) {
	// Return Repository struct with new queries
	return &Repository{
		db: dbpool,
	}, nil
}

// Repository is the repository structure.
type Repository struct {
	db *pgxpool.Pool
}

// CreateUser creates new user in the repository.
func (r *Repository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	return model.User{}, nil
}

// CreateBalance creates new user balance in the repository.
func (r *Repository) CreateBalance(ctx context.Context, user model.User) error {
	return nil
}

// SendCoins transfer given amount of coins from one user to another.
func (r *Repository) SendCoins(ctx context.Context, fromUser model.User, toUser model.User, amount int) error {
	return nil
}

// BuyItem register purhcase of inventory item (merch) for a given user.
func (r *Repository) BuyItem(ctx context.Context, user model.User, item model.InventoryItem) error {
	return nil
}
