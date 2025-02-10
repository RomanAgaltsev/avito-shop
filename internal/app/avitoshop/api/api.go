package api

import (
	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/shop"
	"github.com/RomanAgaltsev/avito-shop/internal/config"
)

// Handler handles all HTTP requests.
type Handler struct {
	cfg     *config.Config
	service shop.Service
}

// NewHandler is a Handler constructor.
func NewHandler(cfg *config.Config, service shop.Service) *Handler {
	return &Handler{
		cfg:     cfg,
		service: service,
	}
}
