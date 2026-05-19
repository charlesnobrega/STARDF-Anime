package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Response represents the standard response structure for bridge calls
type Response struct {
	Status    string         `json:"status"`
	Data      interface{}    `json:"data,omitempty"`
	Error     *ErrorResponse `json:"error,omitempty"`
	Timestamp string         `json:"timestamp"`
}

// ErrorResponse represents error details
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Anime represents an anime entry
type Anime struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Episodes    int        `json:"episodes"`
	ImageURL    string     `json:"imageUrl"`
	AddedAt     time.Time  `json:"addedAt"`
	SyncedAt    *time.Time `json:"syncedAt,omitempty"`
}

// WatchlistItem represents a watchlist entry
type WatchlistItem struct {
	ID       string     `json:"id"`
	AnimeID  string     `json:"animeId"`
	AddedAt  time.Time  `json:"addedAt"`
	SyncedAt *time.Time `json:"syncedAt,omitempty"`
}

// HistoryEntry represents a history entry
type HistoryEntry struct {
	ID            string     `json:"id"`
	AnimeID       string     `json:"animeId"`
	EpisodeNumber int        `json:"episodeNumber"`
	WatchedAt     time.Time  `json:"watchedAt"`
	SyncedAt      *time.Time `json:"syncedAt,omitempty"`
}

// SyncQueueItem represents an item in the sync queue
type SyncQueueItem struct {
	ID         string                 `json:"id"`
	Operation  string                 `json:"operation"`
	EntityType string                 `json:"entityType"`
	Data       map[string]interface{} `json:"data"`
	CreatedAt  time.Time              `json:"createdAt"`
	Synced     bool                   `json:"synced"`
}

// SyncStatus represents the current sync status
type SyncStatus struct {
	IsSyncing     bool       `json:"isSyncing"`
	PendingItems  int        `json:"pendingItems"`
	LastSyncTime  *time.Time `json:"lastSyncTime,omitempty"`
	LastSyncError string     `json:"lastSyncError,omitempty"`
}

// GetAnimes returns a list of animes
// Exposed via GoMobile for Flutter bridge
// Validates: Requirements 3
func GetAnimes() string {
	log.Println("[Bridge] GetAnimes called")

	// TODO: Implement actual anime fetching from database
	// For now, return sample data
	animes := []Anime{
		{
			ID:          "1",
			Title:       "Attack on Titan",
			Description: "A dark fantasy anime about humanity's fight against giant humanoid creatures",
			Episodes:    139,
			ImageURL:    "https://example.com/aot.jpg",
			AddedAt:     time.Now(),
		},
		{
			ID:          "2",
			Title:       "Death Note",
			Description: "A psychological thriller about a notebook that can kill anyone",
			Episodes:    37,
			ImageURL:    "https://example.com/deathnote.jpg",
			AddedAt:     time.Now(),
		},
		{
			ID:          "3",
			Title:       "Demon Slayer",
			Description: "An action anime about demon slayers protecting humanity",
			Episodes:    55,
			ImageURL:    "https://example.com/demonslayer.jpg",
			AddedAt:     time.Now(),
		},
	}

	response := Response{
		Status:    "success",
		Data:      animes,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("[Bridge] Error marshaling response: %v", err)
		return errorResponse("BRIDGE_SERIALIZATION_ERROR", "Failed to serialize response", nil)
	}

	return string(data)
}

// AddToWatchlist adds an anime to the watchlist
// Exposed via GoMobile for Flutter bridge
func AddToWatchlist(animeID string) string {
	log.Printf("[Bridge] AddToWatchlist called with animeID: %s", animeID)

	if animeID == "" {
		return errorResponse("INVALID_PARAMETER", "animeID cannot be empty", nil)
	}

	// TODO: Implement actual watchlist addition to database
	item := WatchlistItem{
		ID:      fmt.Sprintf("wl_%d", time.Now().UnixNano()),
		AnimeID: animeID,
		AddedAt: time.Now(),
	}

	response := Response{
		Status:    "success",
		Data:      item,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("[Bridge] Error marshaling response: %v", err)
		return errorResponse("BRIDGE_SERIALIZATION_ERROR", "Failed to serialize response", nil)
	}

	return string(data)
}

// MarkAsWatched marks an episode as watched
// Exposed via GoMobile for Flutter bridge
func MarkAsWatched(animeID string, episodeID string) string {
	log.Printf("[Bridge] MarkAsWatched called with animeID: %s, episodeID: %s", animeID, episodeID)

	if animeID == "" || episodeID == "" {
		return errorResponse("INVALID_PARAMETER", "animeID and episodeID cannot be empty", nil)
	}

	// TODO: Implement actual history entry creation
	entry := HistoryEntry{
		ID:            fmt.Sprintf("hist_%d", time.Now().UnixNano()),
		AnimeID:       animeID,
		EpisodeNumber: 1, // TODO: Parse from episodeID
		WatchedAt:     time.Now(),
	}

	response := Response{
		Status:    "success",
		Data:      entry,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("[Bridge] Error marshaling response: %v", err)
		return errorResponse("BRIDGE_SERIALIZATION_ERROR", "Failed to serialize response", nil)
	}

	return string(data)
}

// GetSyncStatus returns the current sync status
// Exposed via GoMobile for Flutter bridge
func GetSyncStatus() string {
	log.Println("[Bridge] GetSyncStatus called")

	// TODO: Implement actual sync status retrieval
	status := SyncStatus{
		IsSyncing:    false,
		PendingItems: 0,
	}

	response := Response{
		Status:    "success",
		Data:      status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("[Bridge] Error marshaling response: %v", err)
		return errorResponse("BRIDGE_SERIALIZATION_ERROR", "Failed to serialize response", nil)
	}

	return string(data)
}

// SyncWithDesktop synchronizes with the desktop backend
// Exposed via GoMobile for Flutter bridge
func SyncWithDesktop() string {
	log.Println("[Bridge] SyncWithDesktop called")

	// TODO: Implement actual sync with desktop
	response := Response{
		Status:    "success",
		Data:      map[string]interface{}{"synced": true},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("[Bridge] Error marshaling response: %v", err)
		return errorResponse("BRIDGE_SERIALIZATION_ERROR", "Failed to serialize response", nil)
	}

	return string(data)
}

// errorResponse creates a standard error response
func errorResponse(code string, message string, details map[string]interface{}) string {
	response := Response{
		Status: "error",
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	data, _ := json.Marshal(response)
	return string(data)
}
