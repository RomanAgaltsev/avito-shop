package repository

import (
	"context"
	"github.com/RomanAgaltsev/avito-shop/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
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
func (r *Repository) CreateUser(ctx context.Context, user model.User) error {
	return nil
}

// GetUser returns a user from repository.
func (r *Repository) GetUser(ctx context.Context, login string) (model.User, error) {
	return model.User{}, nil
}
