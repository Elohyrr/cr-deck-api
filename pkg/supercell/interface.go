package supercell

import (
	"context"
)

// Client defines operations to fetch data from Supercell API
type Client interface {
	// GetTopPlayers retrieves top N players from global rankings
	// NOTE 2026-01-11: Endpoint /locations/*/rankings/players retourne des listes vides
	// Bug côté Supercell API. Alternative: utiliser liste statique (voir collector/top_players.go)
	GetTopPlayers(ctx context.Context, limit int) ([]Player, error)

	// GetBattlelog retrieves last 25 battles for a player
	// ✅ Validé fonctionnel - utilisé pour la collecte principale
	GetBattlelog(ctx context.Context, tag string) ([]BattleRaw, error)
}

// Player represents a top player from rankings
type Player struct {
	Tag      string `json:"tag"`
	Name     string `json:"name"`
	Trophies int    `json:"trophies"`
	Rank     int    `json:"rank"`
}

// BattleRaw represents raw battle data from Supercell API
type BattleRaw struct {
	Type       string       `json:"type"`
	BattleTime string       `json:"battleTime"`
	GameMode   GameMode     `json:"gameMode"`
	Team       []TeamMember `json:"team"`
	Opponent   []TeamMember `json:"opponent"`
}

// GameMode represents the game mode information
type GameMode struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TeamMember represents a player in a battle
type TeamMember struct {
	Tag    string    `json:"tag"`
	Name   string    `json:"name"`
	Crowns int       `json:"crowns"`
	Cards  []CardRaw `json:"cards"`
}

// CardRaw represents raw card data from API
type CardRaw struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Level    int      `json:"level"`
	MaxLevel int      `json:"maxLevel"`
	IconUrls IconUrls `json:"iconUrls,omitempty"`
}

// IconUrls contains card icon URLs
type IconUrls struct {
	Medium string `json:"medium,omitempty"`
}
