package api

import (
	"net/http"

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

// Auth handles user registration and authentication.
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {

}

// SendCoin handles send coins request.
func (h *Handler) SendCoin(w http.ResponseWriter, r *http.Request) {

}

// BuyItem handles buy item request.
func (h *Handler) BuyItem(w http.ResponseWriter, r *http.Request) {

}

// Info handles info request.
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {

}
