package supercell

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/leopoldhub/royal-api-personal/internal/errors"
)

func TestHTTPClient_GetTopPlayers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/locations/global/rankings/players" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		limit := r.URL.Query().Get("limit")
		if limit != "10" {
			t.Errorf("expected limit=10, got %s", limit)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test_token" {
			t.Errorf("unexpected auth header: %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"items": [
				{"tag": "#2PP", "name": "Player1", "trophies": 9876, "rank": 1},
				{"tag": "#ABC", "name": "Player2", "trophies": 9800, "rank": 2}
			]
		}`))
	}))
	defer server.Close()

	client := &HTTPClient{
		apiKey:     "test_token",
		baseURL:    server.URL + "/v1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	players, err := client.GetTopPlayers(ctx, 10)

	if err != nil {
		t.Fatalf("GetTopPlayers() error = %v", err)
	}

	if len(players) != 2 {
		t.Errorf("expected 2 players, got %d", len(players))
	}

	if players[0].Tag != "#2PP" {
		t.Errorf("expected tag #2PP, got %s", players[0].Tag)
	}
}

func TestHTTPClient_GetBattlelog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v1/players/") || !strings.HasSuffix(r.URL.Path, "/battlelog") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Note: Le serveur HTTP de test reçoit le path décodé (/players/#2PP)
		// url.PathEscape encode bien # en %23 côté client avant l'envoi

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// L'API Supercell retourne un array direct, pas {"items": [...]}
		w.Write([]byte(`[
			{
				"type": "PvP",
				"battleTime": "20240110T201530.000Z",
				"gameMode": {"id": 72000006, "name": "Ladder"},
				"team": [{"tag": "#2PP", "crowns": 3, "cards": []}],
				"opponent": [{"tag": "#ABC", "crowns": 1, "cards": []}]
			}
		]`))
	}))
	defer server.Close()

	client := &HTTPClient{
		apiKey:     "test_token",
		baseURL:    server.URL + "/v1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	battles, err := client.GetBattlelog(ctx, "#2PP")

	if err != nil {
		t.Fatalf("GetBattlelog() error = %v", err)
	}

	if len(battles) != 1 {
		t.Errorf("expected 1 battle, got %d", len(battles))
	}

	if battles[0].Type != "PvP" {
		t.Errorf("expected type PvP, got %s", battles[0].Type)
	}
}

func TestHTTPClient_RateLimitRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"reason": "rate limit", "message": "too many requests"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items": []}`))
	}))
	defer server.Close()

	client := &HTTPClient{
		apiKey:     "test_token",
		baseURL:    server.URL + "/v1",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	ctx := context.Background()
	_, err := client.GetTopPlayers(ctx, 10)

	if err != nil {
		t.Fatalf("should succeed after retry, got error: %v", err)
	}

	if attempts < 2 {
		t.Errorf("expected at least 2 attempts, got %d", attempts)
	}
}

func TestHTTPClient_NotFoundError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"reason": "notFound", "message": "player not found"}`))
	}))
	defer server.Close()

	client := &HTTPClient{
		apiKey:     "test_token",
		baseURL:    server.URL + "/v1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	_, err := client.GetBattlelog(ctx, "#INVALID")

	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}

	apiErr, ok := err.(*errors.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}

	if !apiErr.IsNotFound() {
		t.Errorf("expected IsNotFound() to be true")
	}
}

func TestParseRetryAfter(t *testing.T) {
	client := &HTTPClient{}

	tests := []struct {
		input    string
		expected int
	}{
		{"10", 10},
		{"", 5},
		{"invalid", 5},
		{"  30  ", 30},
	}

	for _, tt := range tests {
		got := client.parseRetryAfter(tt.input)
		if got != tt.expected {
			t.Errorf("parseRetryAfter(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}
