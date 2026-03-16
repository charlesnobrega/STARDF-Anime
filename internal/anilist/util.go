package anilist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// parseJSON is a helper to decode JSON from an http.Response
func parseJSON(resp *http.Response, v interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	return json.Unmarshal(body, v)
}
