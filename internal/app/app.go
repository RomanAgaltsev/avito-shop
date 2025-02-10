package app

import (
	"log"

	"github.com/RomanAgaltsev/avito-shop/internal/config"
	"github.com/RomanAgaltsev/avito-shop/internal/logger"
)

// Run run the whole application.
func Run(cfg *config.Config) {
	// Logger
	err := logger.Initialize()
	if err != nil {
		log.Fatalf("logger initialization: %s", err)
	}

	// Repository

	// Service

	// HTTP server

	// Shutdown

}
