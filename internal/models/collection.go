package models

import "time"

// CollectionStats tracks statistics for a collection run
type CollectionStats struct {
	ID               int       `json:"id,omitempty"`
	StartedAt        time.Time `json:"started_at"`
	CompletedAt      time.Time `json:"completed_at,omitempty"`
	PlayersProcessed int       `json:"players_processed"`
	BattlesCollected int       `json:"battles_collected"`
	BattlesStored    int       `json:"battles_stored"`
	Errors           int       `json:"errors"`
	Status           string    `json:"status"`
	ErrorMessage     string    `json:"error_message,omitempty"`
}
