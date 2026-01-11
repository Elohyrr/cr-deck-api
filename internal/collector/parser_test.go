package collector

import (
	"testing"
	"time"

	"github.com/leopoldhub/royal-api-personal/pkg/supercell"
)

func TestFilterPvPLadder(t *testing.T) {
	battles := []supercell.BattleRaw{
		{
			Type:     "PvP",
			GameMode: supercell.GameMode{Name: "Ladder"},
		},
		{
			Type:     "PvP",
			GameMode: supercell.GameMode{Name: "Ladder_GoldRush"},
		},
		{
			Type:     "riverRacePvP",
			GameMode: supercell.GameMode{Name: "ClanWar"},
		},
		{
			Type:     "challenge",
			GameMode: supercell.GameMode{Name: "Challenge"},
		},
		{
			Type:     "PvP",
			GameMode: supercell.GameMode{Name: "Tournament"},
		},
	}

	filtered := FilterPvPLadder(battles)

	if len(filtered) != 2 {
		t.Errorf("expected 2 PvP Ladder battles, got %d", len(filtered))
	}

	for _, battle := range filtered {
		if battle.Type != "PvP" {
			t.Errorf("expected type PvP, got %s", battle.Type)
		}
		if !contains(battle.GameMode.Name, "Ladder") {
			t.Errorf("expected gameMode to contain Ladder, got %s", battle.GameMode.Name)
		}
	}
}

func TestParseBattle(t *testing.T) {
	raw := supercell.BattleRaw{
		Type:       "PvP",
		BattleTime: "20240110T201530.000Z",
		GameMode:   supercell.GameMode{ID: 72000006, Name: "Ladder"},
		Team: []supercell.TeamMember{
			{
				Tag:    "#2PP",
				Name:   "Player1",
				Crowns: 3,
				Cards: []supercell.CardRaw{
					{ID: 26000000, Name: "Knight", Level: 14},
					{ID: 26000001, Name: "Archers", Level: 13},
					{ID: 26000003, Name: "Goblins", Level: 14},
					{ID: 26000010, Name: "Giant", Level: 12},
					{ID: 26000014, Name: "Prince", Level: 13},
					{ID: 26000027, Name: "Musketeer", Level: 14},
					{ID: 26000042, Name: "Hog Rider", Level: 13},
					{ID: 28000000, Name: "Fireball", Level: 14},
				},
			},
		},
		Opponent: []supercell.TeamMember{
			{
				Tag:    "#ABC",
				Name:   "Player2",
				Crowns: 1,
				Cards:  []supercell.CardRaw{},
			},
		},
	}

	battle, err := ParseBattle(raw)

	if err != nil {
		t.Fatalf("ParseBattle() error = %v", err)
	}

	if battle == nil {
		t.Fatal("ParseBattle() returned nil")
	}

	if battle.PlayerTag != "#2PP" {
		t.Errorf("PlayerTag = %v, want #2PP", battle.PlayerTag)
	}

	if battle.OpponentTag != "#ABC" {
		t.Errorf("OpponentTag = %v, want #ABC", battle.OpponentTag)
	}

	if battle.PlayerCrowns != 3 {
		t.Errorf("PlayerCrowns = %v, want 3", battle.PlayerCrowns)
	}

	if battle.OpponentCrowns != 1 {
		t.Errorf("OpponentCrowns = %v, want 1", battle.OpponentCrowns)
	}

	if !battle.IsVictory {
		t.Error("IsVictory should be true (3 crowns > 1 crown)")
	}

	if battle.DeckSignature == "" {
		t.Error("DeckSignature should not be empty")
	}

	if len(battle.DeckCards) != 8 {
		t.Errorf("expected 8 cards, got %d", len(battle.DeckCards))
	}
}

func TestParseBattle_InvalidData(t *testing.T) {
	tests := []struct {
		name string
		raw  supercell.BattleRaw
	}{
		{
			name: "empty team",
			raw: supercell.BattleRaw{
				Team:     []supercell.TeamMember{},
				Opponent: []supercell.TeamMember{{Tag: "#ABC"}},
			},
		},
		{
			name: "empty opponent",
			raw: supercell.BattleRaw{
				Team:     []supercell.TeamMember{{Tag: "#2PP"}},
				Opponent: []supercell.TeamMember{},
			},
		},
		{
			name: "invalid card count",
			raw: supercell.BattleRaw{
				Team: []supercell.TeamMember{
					{
						Tag:   "#2PP",
						Cards: []supercell.CardRaw{{ID: 26000000}},
					},
				},
				Opponent: []supercell.TeamMember{{Tag: "#ABC"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			battle, err := ParseBattle(tt.raw)
			if battle != nil {
				t.Errorf("expected nil battle for invalid data, got %v", battle)
			}
			if err != nil {
				t.Logf("got expected error: %v", err)
			}
		})
	}
}

func TestParseBattleTime(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
		wantYear int
	}{
		{
			input:    "20240110T201530.000Z",
			wantErr:  false,
			wantYear: 2024,
		},
		{
			input:   "invalid",
			wantErr: true,
		},
		{
			input:    "20231225T120000.000Z",
			wantErr:  false,
			wantYear: 2023,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseBattleTime(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBattleTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Year() != tt.wantYear {
					t.Errorf("Year = %v, want %v", result.Year(), tt.wantYear)
				}
				if result.Location() != time.UTC {
					t.Errorf("expected UTC timezone, got %v", result.Location())
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
