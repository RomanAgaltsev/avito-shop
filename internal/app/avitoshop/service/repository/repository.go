package repository

import "github.com/jackc/pgx/v5/pgxpool"

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
