package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"net/http"

	"github.com/RomanAgaltsev/avito-shop/internal/config"
	"github.com/RomanAgaltsev/avito-shop/internal/logger"
)

const ContentTypeJSON = "application/json"

func NewRouter(cfg *config.Config, handle *Handler) *chi.Mux {
	// Create router
	router := chi.NewRouter()

	// Enable common middleware
	router.Use(logger.NewRequestLogger())
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5, ContentTypeJSON))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// Replace default handlers
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)

	//	Set routes

	// Public routes
	router.Group(func(r chi.Router) {
		r.Post("/api/auth", handle.Auth)
	})
	// Protected routes
	router.Group(func(r chi.Router) {
		//		tokenAuth := auth.NewAuth(cfg.SecretKey)
		//		r.Use(jwtauth.Verifier(tokenAuth))
		//		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/api/sendCoin", handle.SendCoin)
		r.Get("/api/buy/{item}", handle.BuyItem)
		r.Get("/api/info", handle.Info)
	})

	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", ContentTypeJSON)
	w.WriteHeader(405)
	_ = render.Render(w, r, ErrMethodNotAllowed)
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", ContentTypeJSON)
	w.WriteHeader(400)
	_ = render.Render(w, r, ErrNotFound)
}
