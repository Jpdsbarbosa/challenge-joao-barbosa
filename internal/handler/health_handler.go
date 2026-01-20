package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler gerencia o endpoint de health check
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler cria uma nova instância do handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// Handle processa requisições de health check
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "starkbank-challenge",
		"uptime":    time.Since(h.startTime).String(),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}
