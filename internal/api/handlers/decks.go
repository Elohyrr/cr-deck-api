package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/leopoldhub/royal-api-personal/internal/database/repository"
	"github.com/leopoldhub/royal-api-personal/internal/models"
)

// DeckHandler handles deck-related requests
type DeckHandler struct {
	battleRepo repository.BattleRepository
	metaRepo   repository.MetaDeckRepository
}

// NewDeckHandler creates a new deck handler
func NewDeckHandler(battleRepo repository.BattleRepository, metaRepo repository.MetaDeckRepository) *DeckHandler {
	return &DeckHandler{
		battleRepo: battleRepo,
		metaRepo:   metaRepo,
	}
}

type metaDecksResponse struct {
	Decks    []*models.Deck `json:"decks"`
	Metadata metadata       `json:"metadata"`
}

type deckDetailResponse struct {
	Deck          *models.Deck     `json:"deck"`
	RecentBattles []*models.Battle `json:"recent_battles"`
}

type metadata struct {
	TotalDecks  int       `json:"total_decks"`
	LastUpdated time.Time `json:"last_updated"`
}

// GetMetaDecks handles GET /decks/meta
func (h *DeckHandler) GetMetaDecks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	limit := getQueryInt(r, "limit", 50)
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "win_rate"
	}
	minGames := getQueryInt(r, "min_games", 10)

	decks, err := h.metaRepo.GetTop(ctx, limit, sortBy, minGames)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch meta decks"}`, http.StatusInternalServerError)
		return
	}

	response := metaDecksResponse{
		Decks: decks,
		Metadata: metadata{
			TotalDecks:  len(decks),
			LastUpdated: time.Now(),
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetDeckBySignature handles GET /decks/{signature}
func (h *DeckHandler) GetDeckBySignature(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	signature := r.PathValue("signature")
	if signature == "" {
		http.Error(w, `{"error": "signature required"}`, http.StatusBadRequest)
		return
	}

	deck, err := h.metaRepo.GetBySignature(ctx, signature)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch deck"}`, http.StatusInternalServerError)
		return
	}

	if deck == nil {
		http.Error(w, `{"error": "deck not found"}`, http.StatusNotFound)
		return
	}

	recentBattles, err := h.battleRepo.GetRecent(ctx, signature, 5)
	if err != nil {
		recentBattles = []*models.Battle{}
	}

	response := deckDetailResponse{
		Deck:          deck,
		RecentBattles: recentBattles,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getQueryInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
