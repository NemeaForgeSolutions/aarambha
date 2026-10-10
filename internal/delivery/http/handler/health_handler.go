package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"aarambha-agency/internal/domain"
)

// SystemStats contains runtime diagnostics for the health check.
type SystemStats struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
	AllocMB      uint64 `json:"alloc_mb"`
}

// HealthCheckResponse is the JSON structure returned by the /health endpoint.
type HealthCheckResponse struct {
	Status      string            `json:"status"`
	Timestamp   string            `json:"timestamp"`
	Uptime      string            `json:"uptime"`
	Environment string            `json:"environment,omitempty"`
	Version     string            `json:"version,omitempty"`
	Checks      map[string]string `json:"checks"`
	System      *SystemStats      `json:"system,omitempty"`
}

// HealthHandler manages system liveness, readiness, and diagnostic health checks.
type HealthHandler struct {
	startTime time.Time
	env       string
	version   string
	propRepo  domain.PropertyRepository
	leadRepo  domain.LeadRepository
}

// NewHealthHandler creates an instance of HealthHandler.
func NewHealthHandler(env, version string, propRepo domain.PropertyRepository, leadRepo domain.LeadRepository) *HealthHandler {
	if version == "" {
		version = "1.0.0"
	}
	return &HealthHandler{
		startTime: time.Now(),
		env:       env,
		version:   version,
		propRepo:  propRepo,
		leadRepo:  leadRepo,
	}
}

// Health handles GET /health, returning comprehensive system health diagnostics.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	checks := make(map[string]string)
	if h.propRepo != nil {
		count := h.propRepo.Count()
		checks["properties"] = fmt.Sprintf("healthy (%d properties loaded)", count)
	} else {
		checks["properties"] = "healthy (in-memory)"
	}

	if h.leadRepo != nil {
		checks["leads"] = "healthy"
	} else {
		checks["leads"] = "healthy (in-memory)"
	}

	resp := HealthCheckResponse{
		Status:      "ok",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Uptime:      time.Since(h.startTime).Truncate(time.Second).String(),
		Environment: h.env,
		Version:     h.version,
		Checks:      checks,
		System: &SystemStats{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			AllocMB:      memStats.Alloc / (1024 * 1024),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
