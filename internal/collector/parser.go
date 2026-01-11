package collector

import (
	"fmt"
	"strings"
	"time"

	"github.com/leopoldhub/royal-api-personal/internal/models"
	"github.com/leopoldhub/royal-api-personal/pkg/supercell"
	"github.com/leopoldhub/royal-api-personal/pkg/utils"
)

// FilterPvPLadder filters battles to keep only competitive 1v1 matches
// Includes: PvP Ladder, Path of Legends (Ranked)
// Excludes: Clan Wars, Challenges, 2v2, Party modes
func FilterPvPLadder(battles []supercell.BattleRaw) []supercell.BattleRaw {
	var filtered []supercell.BattleRaw
	for _, battle := range battles {
		// Accept classic Ladder (type: PvP + Ladder in name)
		isClassicLadder := battle.Type == "PvP" && strings.Contains(battle.GameMode.Name, "Ladder")

		// Accept Path of Legends / Ranked (type: pathOfLegend + Ranked in name)
		isRanked := battle.Type == "pathOfLegend" && strings.Contains(battle.GameMode.Name, "Ranked")

		if isClassicLadder || isRanked {
			filtered = append(filtered, battle)
		}
	}
	return filtered
}

// ParseBattle converts a raw battle from API to internal Battle model
func ParseBattle(raw supercell.BattleRaw) (*models.Battle, error) {
	if len(raw.Team) == 0 || len(raw.Opponent) == 0 {
		return nil, nil
	}

	player := raw.Team[0]
	opponent := raw.Opponent[0]

	if len(player.Cards) != 8 {
		return nil, nil
	}

	battleTime, err := ParseBattleTime(raw.BattleTime)
	if err != nil {
		return nil, err
	}

	var deckCards [8]models.Card
	for i, card := range player.Cards {
		deckCards[i] = models.Card{
			ID:    fmt.Sprintf("%d", card.ID), // Convertir int → string
			Name:  card.Name,
			Level: card.Level,
		}
	}

	signature := utils.CalculateDeckSignature(deckCards)
	isVictory := player.Crowns > opponent.Crowns

	battle := &models.Battle{
		BattleTime:     battleTime,
		PlayerTag:      player.Tag,
		OpponentTag:    opponent.Tag,
		GameMode:       raw.GameMode.Name,
		PlayerCrowns:   player.Crowns,
		OpponentCrowns: opponent.Crowns,
		DeckSignature:  signature,
		DeckCards:      deckCards,
		IsVictory:      isVictory,
	}

	return battle, nil
}

// ParseBattleTime parses Supercell API timestamp format (20240110T201530.000Z)
func ParseBattleTime(timestamp string) (time.Time, error) {
	layout := "20060102T150405.000Z"
	return time.Parse(layout, timestamp)
}
