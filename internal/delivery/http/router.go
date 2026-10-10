package http

import (
	"net/http"

	"aarambha-agency/internal/delivery/http/handler"
	"aarambha-agency/internal/delivery/http/middleware"
	"aarambha-agency/internal/service"
)

type RouterConfig struct {
	PageHandler     *handler.PageHandler
	PropertyHandler *handler.PropertyHandler
	LeadHandler     *handler.LeadHandler
	AdminHandler    *handler.AdminHandler
	AuthHandler     *handler.AuthHandler
	AuthService     *service.AuthService
	SEOHandler      *handler.SEOHandler
	HealthHandler   *handler.HealthHandler
}

// NewRouter constructs and configures the HTTP request router with modular domain route registration.
func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	// Register route groups organized by domain in dedicated files:
	// - routes_public.go: static assets, SEO, health, catalog, and public pages
	// - routes_auth.go: authentication, user registration, and profile
	// - routes_admin.go: backoffice properties and agent administration
	// - routes_api.go: REST/JSON endpoints (leads, site visits)
	registerStaticAndMetaRoutes(mux, cfg)
	registerPublicRoutes(mux, cfg)
	registerAuthRoutes(mux, cfg)
	registerAdminRoutes(mux, cfg)
	registerAPIRoutes(mux, cfg)

	var handler http.Handler = mux

	// If AuthService is present, inject Authenticate middleware to populate user context globally
	if cfg.AuthService != nil {
		handler = middleware.Authenticate(cfg.AuthService)(handler)
	}

	handler = middleware.Compress(handler)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.Logger(handler)
	handler = middleware.Recovery(handler)

	return handler
}
