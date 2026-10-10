package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aarambha-agency/configs"
	deliveryHttp "aarambha-agency/internal/delivery/http"
	"aarambha-agency/internal/delivery/http/handler"
	"aarambha-agency/internal/repository/postgres"
	"aarambha-agency/internal/service"
	"aarambha-agency/internal/view"
)

func main() {
	cfg := configs.LoadConfig()

	if cfg.DatabaseURL == "" {
		log.Fatalf("Fatal: DATABASE_URL environment variable is required.")
	}

	log.Printf("Initializing PostgreSQL connection...")
	db, err := postgres.OpenDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Fatal: failed to connect to PostgreSQL database: %v", err)
	}

	migrateCtx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	if err := postgres.ExecuteDDL(migrateCtx, db); err != nil {
		cancel()
		log.Fatalf("Fatal: database schema migration encountered an error: %v", err)
	}
	cancel()

	propRepo := postgres.NewPropertyRepo(db)
	catalogRepo := propRepo
	leadRepo := postgres.NewLeadRepo(db)
	userRepo := postgres.NewUserRepo(db)

	log.Printf("Connected to PostgreSQL DB")

	propertyService := service.NewPropertyService(propRepo, catalogRepo)
	leadService := service.NewLeadService(leadRepo)
	catalogService := service.NewCatalogService(catalogRepo)
	adminService := service.NewAdminService(propRepo, leadRepo)
	authService := service.NewAuthService(userRepo, cfg.SessionSecret)

	viewEngine, err := view.NewEngine()
	if err != nil {
		log.Fatalf("Fatal: failed to initialize view engine: %v", err)
	}

	pageHandler := handler.NewPageHandler(propertyService, catalogService, viewEngine, cfg.Company)
	propertyHandler := handler.NewPropertyHandler(propertyService, viewEngine, cfg.Company)
	leadHandler := handler.NewLeadHandler(leadService)
	adminHandler := handler.NewAdminHandler(propertyService, adminService, catalogService, viewEngine, cfg.Company)
	authHandler := handler.NewAuthHandler(authService, viewEngine, cfg.Company, cfg.IsProduction())
	seoHandler := handler.NewSEOHandler(propertyService, catalogService, "https://aarambha.in")
	healthHandler := handler.NewHealthHandler(cfg.Env, "1.0.0", propRepo, leadRepo)

	router := deliveryHttp.NewRouter(deliveryHttp.RouterConfig{
		PageHandler:     pageHandler,
		PropertyHandler: propertyHandler,
		LeadHandler:     leadHandler,
		AdminHandler:    adminHandler,
		AuthHandler:     authHandler,
		AuthService:     authService,
		SEOHandler:      seoHandler,
		HealthHandler:   healthHandler,
	})

	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Printf("Aarambha enterprise server running on http://%s (Env: %s)",
			cfg.Addr(), cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Received shutdown signal. Stopping server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server exited successfully.")
}
