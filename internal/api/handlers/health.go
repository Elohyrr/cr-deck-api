package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/leopoldhub/royal-api-personal/internal/database/repository"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db         *sql.DB
	battleRepo repository.BattleRepository
	metaRepo   repository.MetaDeckRepository
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB, battleRepo repository.BattleRepository, metaRepo repository.MetaDeckRepository) *HealthHandler {
	return &HealthHandler{
		db:         db,
		battleRepo: battleRepo,
		metaRepo:   metaRepo,
	}
}

type healthResponse struct {
	Status         string    `json:"status"`
	Database       string    `json:"database"`
	LastCollection time.Time `json:"last_collection,omitempty"`
	TotalBattles   int       `json:"total_battles"`
	TotalDecks     int       `json:"total_decks"`
}

// Handle processes health check requests
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	response := healthResponse{
		Status:   "healthy",
		Database: "disconnected",
	}

	if err := h.db.Ping(); err == nil {
		response.Database = "connected"
	}

	if totalBattles, err := h.battleRepo.Count(ctx); err == nil {
		response.TotalBattles = totalBattles
	}

	if totalDecks, err := h.metaRepo.Count(ctx); err == nil {
		response.TotalDecks = totalDecks
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
