package utils

import (
	"sort"
	"strings"

	"github.com/leopoldhub/royal-api-personal/internal/models"
)

// CalculateDeckSignature generates a unique signature for a deck by sorting card IDs
func CalculateDeckSignature(cards [8]models.Card) string {
	ids := make([]string, 8)
	for i, card := range cards {
		ids[i] = card.ID
	}
	sort.Strings(ids)
	return strings.Join(ids, "-")
}

// CalculateDeckSignatureFromSlice generates a signature from a slice of cards
func CalculateDeckSignatureFromSlice(cards []models.Card) string {
	if len(cards) != 8 {
		return ""
	}

	ids := make([]string, 8)
	for i, card := range cards {
		ids[i] = card.ID
	}
	sort.Strings(ids)
	return strings.Join(ids, "-")
}
