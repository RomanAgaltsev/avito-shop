package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/RomanAgaltsev/avito-shop/migrations"
)

// NewConnectionPool creates new pgx connection pool and runs migrations.
func NewConnectionPool(ctx context.Context, databaseURI string) (*pgxpool.Pool, error) {
	// Create new connection pool
	dbpool, err := pgxpool.New(ctx, databaseURI)
	if err != nil {
		slog.Error("new DB connection", slog.String("error", err.Error()))
		return nil, err
	}

	// Ping DB
	if err = dbpool.Ping(ctx); err != nil {
		slog.Error("ping DB", slog.String("error", err.Error()))
		return nil, err
	}

	// Do migrations
	Migrate(ctx, dbpool, databaseURI)

	return dbpool, nil
}

// Migrate runs migrations.
func Migrate(ctx context.Context, dbpool *pgxpool.Pool, databaseURI string) {
	// Set migrations directory
	goose.SetBaseFS(migrations.Migrations)

	// Set dialect
	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("goose: set dialect", slog.String("error", err.Error()))
	}

	// Open connection from db pool
	db := stdlib.OpenDBFromPool(dbpool)

	// Up migrations
	if err := goose.UpContext(ctx, db, "."); err != nil {
		slog.Error("goose: run migrations", slog.String("error", err.Error()))
	}

	// Close connection
	if err := db.Close(); err != nil {
		slog.Error("goose: close connection", slog.String("error", err.Error()))
	}
}
