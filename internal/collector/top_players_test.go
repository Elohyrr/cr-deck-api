package collector

import (
	"testing"
)

func TestGetTopPlayerTags(t *testing.T) {
	tags := GetTopPlayerTags()

	// Vérifier que la liste n'est pas vide
	if len(tags) == 0 {
		t.Error("GetTopPlayerTags should return a non-empty list")
	}

	// Vérifier que tous les tags commencent par #
	for i, tag := range tags {
		if len(tag) == 0 {
			t.Errorf("Tag at index %d is empty", i)
		}
		if tag[0] != '#' {
			t.Errorf("Tag at index %d (%s) should start with '#'", i, tag)
		}
	}

	// Vérifier que la fonction retourne une copie (pas la référence)
	tags[0] = "#MODIFIED"
	tags2 := GetTopPlayerTags()
	if tags2[0] == "#MODIFIED" {
		t.Error("GetTopPlayerTags should return a copy, not the original slice")
	}
}

func TestGetTopPlayerCount(t *testing.T) {
	count := GetTopPlayerCount()

	// Vérifier que le count est cohérent avec GetTopPlayerTags
	tags := GetTopPlayerTags()
	if count != len(tags) {
		t.Errorf("GetTopPlayerCount() = %d, but GetTopPlayerTags() has %d items", count, len(tags))
	}

	// Vérifier qu'on a au moins quelques tags (MVP minimum)
	if count < 5 {
		t.Errorf("Expected at least 5 player tags, got %d", count)
	}
}

func TestTopPlayerTagsFormat(t *testing.T) {
	// Vérifier que les tags respectent le format Clash Royale
	for i, tag := range TopPlayerTags {
		if len(tag) < 3 {
			t.Errorf("Tag at index %d (%s) is too short (minimum 3 chars)", i, tag)
		}
		if len(tag) > 15 {
			t.Errorf("Tag at index %d (%s) is too long (maximum 15 chars)", i, tag)
		}
	}
}

func TestTopPlayerTagsNoDuplicates(t *testing.T) {
	// Vérifier qu'il n'y a pas de doublons
	seen := make(map[string]bool)
	for i, tag := range TopPlayerTags {
		if seen[tag] {
			t.Errorf("Duplicate tag found at index %d: %s", i, tag)
		}
		seen[tag] = true
	}
}
