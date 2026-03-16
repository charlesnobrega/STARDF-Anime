package anilist

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ErrNotAuthenticated is returned when an operation requires authentication
var ErrNotAuthenticated = errors.New("not authenticated with AniList — run with --anilist-login to connect")

// OAuthConfig holds the OAuth2 app credentials
// Users need to register an app at https://anilist.co/settings/developer
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// TokenStore stores and retrieves the access token locally
type TokenStore struct {
	path string
}

// NewTokenStore creates a token store in the platform config dir
func NewTokenStore() *TokenStore {
	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = filepath.Join(os.Getenv("APPDATA"), "stardf-anime")
	case "darwin":
		dir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "stardf-anime")
	default:
		dir = filepath.Join(os.Getenv("HOME"), ".config", "stardf-anime")
	}
	return &TokenStore{path: filepath.Join(dir, "anilist_token")}
}

// Save persists the token to disk
func (ts *TokenStore) Save(token string) error {
	if err := os.MkdirAll(filepath.Dir(ts.path), 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}
	return os.WriteFile(ts.path, []byte(token), 0600)
}

// Load reads the token from disk
func (ts *TokenStore) Load() (string, error) {
	data, err := os.ReadFile(ts.path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// Delete removes the stored token
func (ts *TokenStore) Delete() error {
	return os.Remove(ts.path)
}

// GetAuthorizationURL builds the OAuth2 authorization URL for the user to visit
func GetAuthorizationURL(clientID string, redirectURI string) string {
	return fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code",
		anilistAuthURL, clientID, redirectURI,
	)
}

// ExchangeCode exchanges an authorization code for an access token
// This method makes an HTTP POST to the AniList token endpoint (requires network)
func ExchangeCode(cfg OAuthConfig, code string) (string, error) {
	body := fmt.Sprintf(
		`{"grant_type":"authorization_code","client_id":"%s","client_secret":"%s","redirect_uri":"%s","code":"%s"}`,
		cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI, code,
	)

	req, err := http.NewRequest("POST", anilistTokenURL, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token exchange failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange returned status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := parseJSON(resp, &result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}
