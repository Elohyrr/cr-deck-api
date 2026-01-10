package models

// Card represents a Clash Royale card
type Card struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Level      int    `json:"level,omitempty"`
	MaxLevel   int    `json:"max_level,omitempty"`
	Rarity     string `json:"rarity,omitempty"`
	ElixirCost int    `json:"elixir_cost,omitempty"`
}
