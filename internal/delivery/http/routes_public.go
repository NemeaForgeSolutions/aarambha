package http

import (
	"io/fs"
	"log"
	"net/http"

	"aarambha-agency/internal/delivery/http/handler"
	"aarambha-agency/web"
)

// registerStaticAndMetaRoutes sets up static asset serving, favicons, manifests, SEO, and health checks.
func registerStaticAndMetaRoutes(mux *http.ServeMux, cfg RouterConfig) {
	staticSubFS, err := fs.Sub(web.StaticFS, "static")
	if err != nil {
		log.Fatalf("Failed to create static sub-filesystem: %v", err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS))))

	if cfg.SEOHandler != nil {
		mux.HandleFunc("GET /sitemap.xml", cfg.SEOHandler.SitemapXML)
		mux.HandleFunc("GET /robots.txt", cfg.SEOHandler.RobotsTxt)
	}

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/img/favicon.svg", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /site.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/site.webmanifest", http.StatusMovedPermanently)
	})

	healthHandler := cfg.HealthHandler
	if healthHandler == nil {
		healthHandler = handler.NewHealthHandler("development", "1.0.0", nil, nil)
	}
	mux.HandleFunc("GET /health", healthHandler.Health)
}

// registerPublicRoutes sets up the consumer-facing marketing and property catalog routes.
func registerPublicRoutes(mux *http.ServeMux, cfg RouterConfig) {
	if cfg.PageHandler != nil {
		mux.HandleFunc("GET /{$}", cfg.PageHandler.Home)
		mux.HandleFunc("GET /locations", cfg.PageHandler.Locations)
		mux.HandleFunc("GET /commercial", cfg.PageHandler.Commercial)
		mux.HandleFunc("GET /rent", cfg.PageHandler.Rent)
		mux.HandleFunc("GET /sell", cfg.PageHandler.Sell)
		mux.HandleFunc("GET /new-projects", cfg.PageHandler.NewProjects)
		mux.HandleFunc("GET /agents", cfg.PageHandler.Agents)
		mux.HandleFunc("GET /insights", cfg.PageHandler.Insights)
		mux.HandleFunc("GET /about", cfg.PageHandler.About)
		mux.HandleFunc("GET /contact", cfg.PageHandler.Contact)
	}

	if cfg.PropertyHandler != nil {
		mux.HandleFunc("GET /properties", cfg.PropertyHandler.List)
		mux.HandleFunc("GET /properties/search", cfg.PropertyHandler.Search)
		mux.HandleFunc("GET /properties/featured", cfg.PropertyHandler.Featured)
		mux.HandleFunc("GET /properties/{slug}", cfg.PropertyHandler.Detail)
	}

	// Catch-all: any unmatched route returns the branded 404 page.
	if cfg.PageHandler != nil {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// The mux already handles exact "/" via Home — only non-matching paths fall here.
			cfg.PageHandler.NotFound(w, r)
		})
	}
}
