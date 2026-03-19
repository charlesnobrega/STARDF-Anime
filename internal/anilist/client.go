// Package anilist provides integration with the AniList GraphQL API.
// It supports OAuth2 authentication for syncing user watch progress.
package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	anilistAPI      = "https://graphql.anilist.co"
	anilistAuthURL  = "https://anilist.co/api/v2/oauth/authorize"
	anilistTokenURL = "https://anilist.co/api/v2/oauth/token"
)

// Client is the AniList API client
type Client struct {
	httpClient  *http.Client
	accessToken string
}

// NewClient creates a new unauthenticated AniList client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// NewAuthenticatedClient creates a client with an existing access token
func NewAuthenticatedClient(token string) *Client {
	return &Client{
		httpClient:  &http.Client{Timeout: 15 * time.Second},
		accessToken: token,
	}
}

// IsAuthenticated checks if the client has an access token
func (c *Client) IsAuthenticated() bool {
	return c.accessToken != ""
}

// SetToken sets the access token for authenticated requests
func (c *Client) SetToken(token string) {
	c.accessToken = token
}

// query sends a GraphQL query to the AniList API
func (c *Client) query(query string, variables map[string]interface{}, result interface{}) error {
	body := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", anilistAPI, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AniList API error %d: %s", resp.StatusCode, string(respBody))
	}

	return json.Unmarshal(respBody, result)
}

// GetTrendingSeason fetches trending anime for the current season.
func (c *Client) GetTrendingSeason(page int) ([]MediaSearchResult, error) {
	now := time.Now()
	year := now.Year()
	month := now.Month()

	var season string
	switch {
	case month >= 1 && month <= 3:
		season = "WINTER"
	case month >= 4 && month <= 6:
		season = "SPRING"
	case month >= 7 && month <= 9:
		season = "SUMMER"
	default:
		season = "FALL"
	}

	var resp struct {
		Data struct {
			Page struct {
				Media []MediaSearchResult `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	err := c.query(queryGetTrendingSeason, map[string]interface{}{
		"season":     season,
		"seasonYear": year,
		"page":       page,
	}, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Data.Page.Media, nil
}
