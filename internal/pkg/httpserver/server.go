package httpserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"

	"github.com/RomanAgaltsev/avito-shop/internal/config"
	"github.com/RomanAgaltsev/avito-shop/internal/logger"
)

const ContentTypeJSON = "application/json"

var ErrRunAddressIsEmpty = fmt.Errorf("configuration: HTTP server run address is empty")

// New creates new http server with middleware and routes.
func New(cfg *config.Config) (*http.Server, error) {
	if cfg.RunAddress == "" {
		return nil, ErrRunAddressIsEmpty
	}

	// Create handler
	//handle := api.NewHandler(cfg)

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

	/*
		Set routes
	*/

	// Public routes
	router.Group(func(r chi.Router) {

	})
	// Protected routes
	router.Group(func(r chi.Router) {
		//tokenAuth := auth.NewAuth(cfg.SecretKey)
		//r.Use(jwtauth.Verifier(tokenAuth))
		//r.Use(jwtauth.Authenticator(tokenAuth))

	})

	return &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}, nil
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
