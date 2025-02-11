package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"

	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/shop"
	"github.com/RomanAgaltsev/avito-shop/internal/config"
	"github.com/RomanAgaltsev/avito-shop/internal/model"
	"github.com/RomanAgaltsev/avito-shop/internal/pkg/auth"
)

const (
	contentTypeJSON = "application/json"

	argError = "error"

	msgNewJWTToken    = "new JWT token"
	msgUserLogin      = "user login"
	msgNewUserBalance = "new user balance"
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
	// Get user from request
	var usr model.User
	if err := render.Bind(r, &usr); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}

	// Get context from request
	ctx := r.Context()

	// Login user
	err := h.service.UserLogin(ctx, usr)
	if err != nil {
		// Something has gone wrong
		slog.Info(msgUserLogin, argError, err.Error())
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}

	// Create a balance for the user
	err = h.service.UserBalance(ctx, usr)
	if err != nil {
		slog.Info(msgNewUserBalance, argError, err.Error())
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}

	// Generate JWT token
	ja := auth.NewAuth(h.cfg.SecretKey)
	_, tokenString, err := auth.NewJWTToken(ja, usr.UserName)
	if err != nil {
		// Something has gone wrong
		slog.Info(msgNewJWTToken, argError, err.Error())
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)

	authResponse := model.AuthResponse{
		Token: tokenString,
	}

	// Render the list of orders to response
	if err = render.Render(w, r, &authResponse); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}
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
