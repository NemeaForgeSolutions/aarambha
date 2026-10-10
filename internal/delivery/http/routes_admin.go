package http

import (
	"net/http"

	"aarambha-agency/internal/delivery/http/middleware"
)

// registerAdminRoutes configures protected backoffice administration routes.
func registerAdminRoutes(mux *http.ServeMux, cfg RouterConfig) {
	if cfg.AdminHandler == nil {
		return
	}

	if cfg.AuthService != nil {
		requireAdmin := middleware.RequireAdmin(cfg.AuthService)

		// Property Inventory & Lead Management
		mux.Handle("GET /admin", requireAdmin(http.HandlerFunc(cfg.AdminHandler.Dashboard)))
		mux.Handle("GET /admin/properties", requireAdmin(http.HandlerFunc(cfg.AdminHandler.Dashboard)))
		mux.Handle("GET /admin/properties/new", requireAdmin(http.HandlerFunc(cfg.AdminHandler.NewProperty)))
		mux.Handle("POST /admin/properties/new", requireAdmin(http.HandlerFunc(cfg.AdminHandler.CreateProperty)))
		mux.Handle("GET /admin/properties/edit/{id}", requireAdmin(http.HandlerFunc(cfg.AdminHandler.EditProperty)))
		mux.Handle("POST /admin/properties/edit/{id}", requireAdmin(http.HandlerFunc(cfg.AdminHandler.UpdateProperty)))
		mux.Handle("POST /admin/properties/delete/{id}", requireAdmin(http.HandlerFunc(cfg.AdminHandler.DeleteProperty)))
		mux.Handle("POST /admin/properties/toggle-featured/{id}", requireAdmin(http.HandlerFunc(cfg.AdminHandler.ToggleFeatured)))

		// Sell Requests Management
		mux.Handle("POST /admin/sell-requests/{id}/status", requireAdmin(http.HandlerFunc(cfg.AdminHandler.UpdateSellRequestStatus)))

		// Advisory Agents Administration
		mux.Handle("GET /admin/agents", requireAdmin(http.HandlerFunc(cfg.AdminHandler.AgentsList)))
		mux.Handle("GET /admin/agents/new", requireAdmin(http.HandlerFunc(cfg.AdminHandler.NewAgent)))
		mux.Handle("POST /admin/agents/new", requireAdmin(http.HandlerFunc(cfg.AdminHandler.CreateAgent)))
		mux.Handle("GET /admin/agents/edit/{id}", requireAdmin(http.HandlerFunc(cfg.AdminHandler.EditAgent)))
		mux.Handle("POST /admin/agents/edit/{id}", requireAdmin(http.HandlerFunc(cfg.AdminHandler.UpdateAgent)))
		mux.Handle("POST /admin/agents/delete/{id}", requireAdmin(http.HandlerFunc(cfg.AdminHandler.DeleteAgent)))
	} else {
		// Fallback without auth service (e.g. testing)
		mux.HandleFunc("GET /admin", cfg.AdminHandler.Dashboard)
		mux.HandleFunc("GET /admin/properties", cfg.AdminHandler.Dashboard)
		mux.HandleFunc("GET /admin/properties/new", cfg.AdminHandler.NewProperty)
		mux.HandleFunc("POST /admin/properties/new", cfg.AdminHandler.CreateProperty)
		mux.HandleFunc("GET /admin/properties/edit/{id}", cfg.AdminHandler.EditProperty)
		mux.HandleFunc("POST /admin/properties/edit/{id}", cfg.AdminHandler.UpdateProperty)
		mux.HandleFunc("POST /admin/properties/delete/{id}", cfg.AdminHandler.DeleteProperty)
		mux.HandleFunc("POST /admin/properties/toggle-featured/{id}", cfg.AdminHandler.ToggleFeatured)

		// Sell Requests Management
		mux.HandleFunc("POST /admin/sell-requests/{id}/status", cfg.AdminHandler.UpdateSellRequestStatus)

		// Advisory Agents Administration
		mux.HandleFunc("GET /admin/agents", cfg.AdminHandler.AgentsList)
		mux.HandleFunc("GET /admin/agents/new", cfg.AdminHandler.NewAgent)
		mux.HandleFunc("POST /admin/agents/new", cfg.AdminHandler.CreateAgent)
		mux.HandleFunc("GET /admin/agents/edit/{id}", cfg.AdminHandler.EditAgent)
		mux.HandleFunc("POST /admin/agents/edit/{id}", cfg.AdminHandler.UpdateAgent)
		mux.HandleFunc("POST /admin/agents/delete/{id}", cfg.AdminHandler.DeleteAgent)
	}
}
