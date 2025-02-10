package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/repository"
	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/shop"
	"github.com/RomanAgaltsev/avito-shop/internal/config"
	"github.com/RomanAgaltsev/avito-shop/internal/database"
	"github.com/RomanAgaltsev/avito-shop/internal/logger"
	"github.com/RomanAgaltsev/avito-shop/internal/pkg/httpserver"
)

// Run runs the whole application.
func Run(cfg *config.Config) error {
	// Logger
	err := logger.Initialize()
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Repository
	// Create connection pool
	dbpool, err := database.NewConnectionPool(ctx, cfg.DatabaseURI)
	if err != nil {
		return fmt.Errorf("connection pool: %w", err)
	}
	defer dbpool.Close()

	// Create repository
	repo, err := repository.New(dbpool)
	if err != nil {
		return err
	}

	// Create shop service
	shopService, err := shop.NewService(repo, cfg)
	if err != nil {
		return nil
	}

	// HTTP server
	server, err := httpserver.New(cfg, shopService)
	if err != nil {
		return err
	}

	// Create channels for graceful shutdown
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	// Interrupt signal
	signal.Notify(quit, os.Interrupt)

	// Graceful shutdown executes in a goroutine
	go func() {
		<-quit

		slog.Info("shutting down HTTP server")

		// Shutdown HTTP server
		if err = server.Shutdown(ctx); err != nil {
			slog.Error("HTTP server shutdown error", slog.String("error", err.Error()))
		}

		close(done)
	}()

	slog.Info("starting HTTP server", "addr", server.Addr)

	// Run HTTP server
	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server error", slog.String("error", err.Error()))
		return err
	}

	<-done
	slog.Info("HTTP server stopped")
	return nil
}
