package http

import (
	"net/http"

	"aarambha-agency/internal/delivery/http/middleware"
)

// registerAuthRoutes configures user authentication, session termination, and profile management routes.
func registerAuthRoutes(mux *http.ServeMux, cfg RouterConfig) {
	if cfg.AuthHandler == nil {
		return
	}

	// Public guest routes
	mux.HandleFunc("GET /login", cfg.AuthHandler.LoginPage)
	mux.HandleFunc("POST /login", cfg.AuthHandler.LoginSubmit)
	mux.HandleFunc("GET /logout", cfg.AuthHandler.Logout)
	mux.HandleFunc("POST /logout", cfg.AuthHandler.Logout)

	// Protected user profile routes
	if cfg.AuthService != nil {
		requireAuth := middleware.RequireAuth(cfg.AuthService)
		mux.Handle("GET /profile", requireAuth(http.HandlerFunc(cfg.AuthHandler.ProfilePage)))
		mux.Handle("POST /profile", requireAuth(http.HandlerFunc(cfg.AuthHandler.ProfileUpdate)))
	} else {
		mux.HandleFunc("GET /profile", cfg.AuthHandler.ProfilePage)
		mux.HandleFunc("POST /profile", cfg.AuthHandler.ProfileUpdate)
	}
}
