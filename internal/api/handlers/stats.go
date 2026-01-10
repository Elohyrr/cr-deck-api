package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/leopoldhub/royal-api-personal/internal/database/repository"
)

// StatsHandler handles statistics requests
type StatsHandler struct {
	battleRepo repository.BattleRepository
	metaRepo   repository.MetaDeckRepository
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(battleRepo repository.BattleRepository, metaRepo repository.MetaDeckRepository) *StatsHandler {
	return &StatsHandler{
		battleRepo: battleRepo,
		metaRepo:   metaRepo,
	}
}

type summaryResponse struct {
	Collection collectionStats `json:"collection"`
	TopDeck    *deckSummary    `json:"top_deck,omitempty"`
	BestDeck   *deckSummary    `json:"best_deck,omitempty"`
	TopCards   []cardUsage     `json:"top_cards"`
}

type collectionStats struct {
	TotalBattles   int       `json:"total_battles"`
	TotalDecks     int       `json:"total_decks"`
	LastCollection time.Time `json:"last_collection"`
	PlayersTracked int       `json:"players_tracked"`
}

type deckSummary struct {
	Signature  string  `json:"signature"`
	TotalGames int     `json:"total_games,omitempty"`
	WinRate    float64 `json:"win_rate,omitempty"`
}

type cardUsage struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	UsageRate float64 `json:"usage_rate"`
}

// GetSummary handles GET /stats/summary
func (h *StatsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	response := summaryResponse{
		Collection: collectionStats{
			PlayersTracked: 1000,
		},
		TopCards: []cardUsage{},
	}

	if totalBattles, err := h.battleRepo.Count(ctx); err == nil {
		response.Collection.TotalBattles = totalBattles
	}

	if totalDecks, err := h.metaRepo.Count(ctx); err == nil {
		response.Collection.TotalDecks = totalDecks
	}

	topByFrequency, err := h.metaRepo.GetTop(ctx, 1, "frequency", 1)
	if err == nil && len(topByFrequency) > 0 {
		deck := topByFrequency[0]
		response.TopDeck = &deckSummary{
			Signature:  deck.Signature,
			TotalGames: deck.Stats.TotalGames,
		}
	}

	topByWinRate, err := h.metaRepo.GetTop(ctx, 1, "win_rate", 20)
	if err == nil && len(topByWinRate) > 0 {
		deck := topByWinRate[0]
		response.BestDeck = &deckSummary{
			Signature: deck.Signature,
			WinRate:   deck.Stats.WinRate,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
