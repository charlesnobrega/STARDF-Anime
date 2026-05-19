package main

import (
	"encoding/json"
	"testing"
)

// TestGetAnimes tests the GetAnimes function
func TestGetAnimes(t *testing.T) {
	result := GetAnimes()

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	if response.Data == nil {
		t.Error("Expected data to be non-nil")
	}

	animes, ok := response.Data.([]interface{})
	if !ok {
		t.Errorf("Expected data to be a list, got %T", response.Data)
	}

	if len(animes) == 0 {
		t.Error("Expected at least one anime in response")
	}
}

// TestAddToWatchlist tests the AddToWatchlist function
func TestAddToWatchlist(t *testing.T) {
	tests := []struct {
		name      string
		animeID   string
		expectErr bool
	}{
		{
			name:      "Valid anime ID",
			animeID:   "1",
			expectErr: false,
		},
		{
			name:      "Empty anime ID",
			animeID:   "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddToWatchlist(tt.animeID)

			var response Response
			err := json.Unmarshal([]byte(result), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tt.expectErr {
				if response.Status != "error" {
					t.Errorf("Expected status 'error', got '%s'", response.Status)
				}
				if response.Error == nil {
					t.Error("Expected error to be non-nil")
				}
			} else {
				if response.Status != "success" {
					t.Errorf("Expected status 'success', got '%s'", response.Status)
				}
				if response.Data == nil {
					t.Error("Expected data to be non-nil")
				}
			}
		})
	}
}

// TestMarkAsWatched tests the MarkAsWatched function
func TestMarkAsWatched(t *testing.T) {
	tests := []struct {
		name      string
		animeID   string
		episodeID string
		expectErr bool
	}{
		{
			name:      "Valid parameters",
			animeID:   "1",
			episodeID: "ep1",
			expectErr: false,
		},
		{
			name:      "Empty anime ID",
			animeID:   "",
			episodeID: "ep1",
			expectErr: true,
		},
		{
			name:      "Empty episode ID",
			animeID:   "1",
			episodeID: "",
			expectErr: true,
		},
		{
			name:      "Both empty",
			animeID:   "",
			episodeID: "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MarkAsWatched(tt.animeID, tt.episodeID)

			var response Response
			err := json.Unmarshal([]byte(result), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tt.expectErr {
				if response.Status != "error" {
					t.Errorf("Expected status 'error', got '%s'", response.Status)
				}
			} else {
				if response.Status != "success" {
					t.Errorf("Expected status 'success', got '%s'", response.Status)
				}
			}
		})
	}
}

// TestGetSyncStatus tests the GetSyncStatus function
func TestGetSyncStatus(t *testing.T) {
	result := GetSyncStatus()

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	if response.Data == nil {
		t.Error("Expected data to be non-nil")
	}
}

// TestSyncWithDesktop tests the SyncWithDesktop function
func TestSyncWithDesktop(t *testing.T) {
	result := SyncWithDesktop()

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	if response.Data == nil {
		t.Error("Expected data to be non-nil")
	}
}

// TestResponseSerialization tests JSON serialization of Response
func TestResponseSerialization(t *testing.T) {
	response := Response{
		Status: "success",
		Data: map[string]interface{}{
			"id":    "1",
			"title": "Test",
		},
		Timestamp: "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled Response
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Status != response.Status {
		t.Errorf("Status mismatch: expected '%s', got '%s'", response.Status, unmarshaled.Status)
	}
}

// TestErrorResponseSerialization tests JSON serialization of error responses
func TestErrorResponseSerialization(t *testing.T) {
	response := Response{
		Status: "error",
		Error: &ErrorResponse{
			Code:    "TEST_ERROR",
			Message: "Test error message",
			Details: map[string]interface{}{
				"key": "value",
			},
		},
		Timestamp: "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled Response
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", unmarshaled.Status)
	}

	if unmarshaled.Error == nil {
		t.Error("Expected error to be non-nil")
	}

	if unmarshaled.Error.Code != "TEST_ERROR" {
		t.Errorf("Error code mismatch: expected 'TEST_ERROR', got '%s'", unmarshaled.Error.Code)
	}
}

// TestDataTypes tests serialization of different data types
func TestDataTypes(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "String",
			data: "test string",
		},
		{
			name: "Integer",
			data: 42,
		},
		{
			name: "Float",
			data: 3.14,
		},
		{
			name: "Boolean",
			data: true,
		},
		{
			name: "Array",
			data: []interface{}{1, 2, 3},
		},
		{
			name: "Object",
			data: map[string]interface{}{
				"key": "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := Response{
				Status:    "success",
				Data:      tt.data,
				Timestamp: "2024-01-01T00:00:00Z",
			}

			data, err := json.Marshal(response)
			if err != nil {
				t.Fatalf("Failed to marshal response: %v", err)
			}

			var unmarshaled Response
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if unmarshaled.Status != "success" {
				t.Errorf("Expected status 'success', got '%s'", unmarshaled.Status)
			}
		})
	}
}
