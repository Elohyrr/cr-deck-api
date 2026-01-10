package models

import "time"

// Battle represents a single battle record
type Battle struct {
	ID             int       `json:"id,omitempty"`
	BattleTime     time.Time `json:"battle_time"`
	PlayerTag      string    `json:"player_tag"`
	OpponentTag    string    `json:"opponent_tag,omitempty"`
	GameMode       string    `json:"game_mode,omitempty"`
	PlayerCrowns   int       `json:"player_crowns"`
	OpponentCrowns int       `json:"opponent_crowns"`
	DeckSignature  string    `json:"deck_signature"`
	DeckCards      [8]Card   `json:"deck_cards"`
	IsVictory      bool      `json:"is_victory"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}
