package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/shop"
	"github.com/RomanAgaltsev/avito-shop/internal/config"
	"github.com/RomanAgaltsev/avito-shop/internal/model"
	"github.com/RomanAgaltsev/avito-shop/internal/pkg/auth"
)

const (
	contentTypeJSON = "application/json"

	argError = "error"

	msgNewJWTToken = "new JWT token"
	msgUserAuth    = "user auth"
	msgSendCoins   = "send coins"
	msgBuyItem     = "buy item"
	msgUserInfo    = "user info"
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
	var user model.User
	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}

	// Get context from request
	ctx := r.Context()

	// Auth user
	err := h.service.UserAuth(ctx, user)
	//
	if err != nil && errors.Is(err, shop.ErrWrongUserNamePassword) {
		// There is a problem with login/password
		slog.Info(msgUserAuth, argError, err.Error())
		_ = render.Render(w, r, ErrWrongLoginPassword)
		return
	}
	//
	if err != nil && !errors.Is(err, shop.ErrWrongUserNamePassword) {
		// Something has gone wrong
		slog.Info(msgUserAuth, argError, err.Error())
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}

	// Generate JWT token
	ja := auth.NewAuth(h.cfg.SecretKey)
	_, tokenString, err := auth.NewJWTToken(ja, user.UserName)
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

	// Render the response
	if err = render.Render(w, r, &authResponse); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}
}

// SendCoins handles send coins request.
func (h *Handler) SendCoins(w http.ResponseWriter, r *http.Request) {
	// Get context from request
	ctx := r.Context()

	// Get user from request
	fromUser, err := auth.UserFromRequest(r, h.cfg.SecretKey)
	if err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}

	var coinsSending model.CoinsSending
	if err = render.Bind(r, &coinsSending); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}

	toUser := model.User{
		UserName: coinsSending.ToUser,
	}
	amount := coinsSending.Amount

	// Send coins
	err = h.service.SendCoins(ctx, fromUser, toUser, amount)
	if err != nil {
		// Something has gone wrong
		slog.Info(msgSendCoins, argError, err.Error())
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// BuyItem handles buy item request.
func (h *Handler) BuyItem(w http.ResponseWriter, r *http.Request) {
	// Get context from request
	ctx := r.Context()

	// Get user from request
	user, err := auth.UserFromRequest(r, h.cfg.SecretKey)
	if err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}

	item := model.InventoryItem{
		Type:     chi.URLParam(r, "item"),
		Quantity: 1,
	}

	err = h.service.BuyItem(ctx, user, item)
	if err != nil {
		// Something has gone wrong
		slog.Info(msgBuyItem, argError, err.Error())
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// Info handles info request.
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	// Get context from request
	ctx := r.Context()

	// Get user from request
	user, err := auth.UserFromRequest(r, h.cfg.SecretKey)
	if err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}

	info, err := h.service.UserInfo(ctx, user)
	if err != nil {
		// Something has gone wrong
		slog.Info(msgUserInfo, argError, err.Error())
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}

	// Set header
	w.Header().Set("Content-type", contentTypeJSON)
	render.Status(r, http.StatusOK)

	// Render the list of orders to response
	if err = render.Render(w, r, &info); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}
}
