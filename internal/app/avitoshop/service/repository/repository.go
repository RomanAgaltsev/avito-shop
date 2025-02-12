package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RomanAgaltsev/avito-shop/internal/model"
)

var (
	ErrConflict = fmt.Errorf("data conflict")
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
