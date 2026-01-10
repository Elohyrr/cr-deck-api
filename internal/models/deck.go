package models

import "time"

// Deck represents a deck of 8 cards with optional statistics
type Deck struct {
	Signature string     `json:"signature"`
	Cards     [8]Card    `json:"cards"`
	Stats     *DeckStats `json:"stats,omitempty"`
}

// DeckStats contains aggregated statistics for a deck
type DeckStats struct {
	TotalGames int       `json:"total_games"`
	Wins       int       `json:"wins"`
	Losses     int       `json:"losses"`
	WinRate    float64   `json:"win_rate"`
	FirstSeen  time.Time `json:"first_seen,omitempty"`
	LastSeen   time.Time `json:"last_seen"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
