package utils

import (
	"testing"

	"github.com/leopoldhub/royal-api-personal/internal/models"
)

func TestCalculateDeckSignature(t *testing.T) {
	cards := [8]models.Card{
		{ID: "26000003"},
		{ID: "26000000"},
		{ID: "26000010"},
		{ID: "26000001"},
		{ID: "26000042"},
		{ID: "26000027"},
		{ID: "28000000"},
		{ID: "26000014"},
	}

	expected := "26000000-26000001-26000003-26000010-26000014-26000027-26000042-28000000"
	got := CalculateDeckSignature(cards)

	if got != expected {
		t.Errorf("CalculateDeckSignature() = %v, want %v", got, expected)
	}
}

func TestCalculateDeckSignature_SameCardsOrder(t *testing.T) {
	cards1 := [8]models.Card{
		{ID: "26000000"}, {ID: "26000001"}, {ID: "26000003"}, {ID: "26000010"},
		{ID: "26000014"}, {ID: "26000027"}, {ID: "26000042"}, {ID: "28000000"},
	}

	cards2 := [8]models.Card{
		{ID: "28000000"}, {ID: "26000042"}, {ID: "26000027"}, {ID: "26000014"},
		{ID: "26000010"}, {ID: "26000003"}, {ID: "26000001"}, {ID: "26000000"},
	}

	sig1 := CalculateDeckSignature(cards1)
	sig2 := CalculateDeckSignature(cards2)

	if sig1 != sig2 {
		t.Errorf("Signatures should be identical regardless of order: %v != %v", sig1, sig2)
	}
}

func TestCalculateDeckSignatureFromSlice(t *testing.T) {
	cards := []models.Card{
		{ID: "26000003"},
		{ID: "26000000"},
		{ID: "26000010"},
		{ID: "26000001"},
		{ID: "26000042"},
		{ID: "26000027"},
		{ID: "28000000"},
		{ID: "26000014"},
	}

	expected := "26000000-26000001-26000003-26000010-26000014-26000027-26000042-28000000"
	got := CalculateDeckSignatureFromSlice(cards)

	if got != expected {
		t.Errorf("CalculateDeckSignatureFromSlice() = %v, want %v", got, expected)
	}
}

func TestCalculateDeckSignatureFromSlice_InvalidLength(t *testing.T) {
	cards := []models.Card{
		{ID: "26000000"},
		{ID: "26000001"},
	}

	got := CalculateDeckSignatureFromSlice(cards)
	if got != "" {
		t.Errorf("CalculateDeckSignatureFromSlice() with invalid length should return empty string, got %v", got)
	}
}
