package supercell

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/leopoldhub/royal-api-personal/internal/errors"
)

// HTTPClient implements the Client interface
type HTTPClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

var _ Client = (*HTTPClient)(nil)

// NewClient creates a new Supercell API client
func NewClient(apiKey string) Client {
	return &HTTPClient{
		apiKey:  apiKey,
		baseURL: "https://api.clashroyale.com/v1",
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetTopPlayers retrieves top N players from global rankings
func (c *HTTPClient) GetTopPlayers(ctx context.Context, limit int) ([]Player, error) {
	endpoint := fmt.Sprintf("/locations/global/rankings/players?limit=%d", limit)

	var response struct {
		Items []Player `json:"items"`
	}

	if err := c.doRequest(ctx, endpoint, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// GetBattlelog retrieves recent battles for a player
func (c *HTTPClient) GetBattlelog(ctx context.Context, tag string) ([]BattleRaw, error) {
	encodedTag := url.PathEscape(tag)
	endpoint := fmt.Sprintf("/players/%s/battlelog", encodedTag)

	var response struct {
		Items []BattleRaw `json:"items"`
	}

	if err := c.doRequest(ctx, endpoint, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// doRequest performs HTTP request with retry logic
func (c *HTTPClient) doRequest(ctx context.Context, endpoint string, result interface{}) error {
	const maxRetries = 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+endpoint, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Second)
				continue
			}
			return fmt.Errorf("request failed after %d attempts: %w", maxRetries, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := c.parseRetryAfter(resp.Header.Get("Retry-After"))
			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(retryAfter) * time.Second)
				continue
			}
			return &errors.APIError{
				StatusCode: resp.StatusCode,
				Message:    "rate limit exceeded",
				Endpoint:   endpoint,
				RetryAfter: retryAfter,
			}
		}

		if resp.StatusCode == http.StatusNotFound {
			return &errors.APIError{
				StatusCode: resp.StatusCode,
				Message:    "resource not found",
				Endpoint:   endpoint,
			}
		}

		if resp.StatusCode >= 500 {
			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Second)
				continue
			}
			return &errors.APIError{
				StatusCode: resp.StatusCode,
				Message:    "server error",
				Endpoint:   endpoint,
			}
		}

		if resp.StatusCode != http.StatusOK {
			var apiErr struct {
				Reason  string `json:"reason"`
				Message string `json:"message"`
			}
			if err := json.Unmarshal(body, &apiErr); err == nil {
				return &errors.APIError{
					StatusCode: resp.StatusCode,
					Message:    fmt.Sprintf("%s: %s", apiErr.Reason, apiErr.Message),
					Endpoint:   endpoint,
				}
			}
			return &errors.APIError{
				StatusCode: resp.StatusCode,
				Message:    string(body),
				Endpoint:   endpoint,
			}
		}

		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		return nil
	}

	return fmt.Errorf("max retries exceeded")
}

// parseRetryAfter parses the Retry-After header value
func (c *HTTPClient) parseRetryAfter(value string) int {
	if value == "" {
		return 5
	}

	if seconds, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
		return seconds
	}

	return 5
}
