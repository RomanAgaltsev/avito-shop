package httpserver

import (
	"fmt"
	"net/http"

	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/api"
	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/shop"
	"github.com/RomanAgaltsev/avito-shop/internal/config"
)

var ErrRunAddressIsEmpty = fmt.Errorf("configuration: HTTP server run address is empty")

// New creates new http server with middleware and routes.
func New(cfg *config.Config, service shop.Service) (*http.Server, error) {
	if cfg.RunAddress == "" {
		return nil, ErrRunAddressIsEmpty
	}

	// Create handler
	handle := api.NewHandler(cfg, service)

	// Create router
	router := api.NewRouter(cfg, handle)

	// Return *http.Server
	return &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}, nil
}
