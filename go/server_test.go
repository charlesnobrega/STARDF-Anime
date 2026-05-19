package main

import (
	"encoding/json"
	"testing"
	"time"
)

// TestStartWebServer tests the StartWebServer function
func TestStartWebServer(t *testing.T) {
	// Reset server state
	mu.Lock()
	isRunning = false
	isPaused = false
	server = nil
	mu.Unlock()

	result := StartWebServer(8080)

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	// Verify server is running
	mu.Lock()
	if !isRunning {
		t.Error("Expected server to be running")
	}
	mu.Unlock()

	// Cleanup
	ShutdownWebServer()
}

// TestStartWebServerAlreadyRunning tests starting server when already running
func TestStartWebServerAlreadyRunning(t *testing.T) {
	// Reset server state
	mu.Lock()
	isRunning = false
	isPaused = false
	server = nil
	mu.Unlock()

	// Start server first time
	StartWebServer(8080)

	// Try to start again
	result := StartWebServer(8080)

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	// Cleanup
	ShutdownWebServer()
}

// TestPauseWebServer tests the PauseWebServer function
func TestPauseWebServer(t *testing.T) {
	// Reset server state
	mu.Lock()
	isRunning = false
	isPaused = false
	server = nil
	mu.Unlock()

	// Start server first
	StartWebServer(8080)

	result := PauseWebServer()

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	// Verify server is paused
	mu.Lock()
	if !isPaused {
		t.Error("Expected server to be paused")
	}
	mu.Unlock()

	// Cleanup
	ShutdownWebServer()
}

// TestPauseWebServerNotRunning tests pausing when server is not running
func TestPauseWebServerNotRunning(t *testing.T) {
	// Reset server state
	mu.Lock()
	isRunning = false
	isPaused = false
	server = nil
	mu.Unlock()

	result := PauseWebServer()

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}
}

// TestResumeWebServer tests the ResumeWebServer function
func TestResumeWebServer(t *testing.T) {
	// Reset server state
	mu.Lock()
	isRunning = false
	isPaused = false
	server = nil
	mu.Unlock()

	// Start and pause server
	StartWebServer(8080)
	PauseWebServer()

	result := ResumeWebServer()

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	// Verify server is resumed
	mu.Lock()
	if isPaused {
		t.Error("Expected server to be resumed (not paused)")
	}
	mu.Unlock()

	// Cleanup
	ShutdownWebServer()
}

// TestResumeWebServerNotRunning tests resuming when server is not running
func TestResumeWebServerNotRunning(t *testing.T) {
	// Reset server state
	mu.Lock()
	isRunning = false
	isPaused = false
	server = nil
	mu.Unlock()

	result := ResumeWebServer()

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", response.Status)
	}

	if response.Error == nil {
		t.Error("Expected error to be non-nil")
	}

	if response.Error.Code != "SERVER_NOT_RUNNING" {
		t.Errorf("Expected error code 'SERVER_NOT_RUNNING', got '%s'", response.Error.Code)
	}
}

// TestShutdownWebServer tests the ShutdownWebServer function
func TestShutdownWebServer(t *testing.T) {
	// Reset server state
	mu.Lock()
	isRunning = false
	isPaused = false
	server = nil
	mu.Unlock()

	// Start server first
	StartWebServer(8080)

	// Give server time to start
	time.Sleep(200 * time.Millisecond)

	result := ShutdownWebServer()

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	// Verify server is not running
	mu.Lock()
	if isRunning {
		t.Error("Expected server to be stopped")
	}
	mu.Unlock()
}

// TestShutdownWebServerNotRunning tests shutting down when server is not running
func TestShutdownWebServerNotRunning(t *testing.T) {
	// Reset server state
	mu.Lock()
	isRunning = false
	isPaused = false
	server = nil
	mu.Unlock()

	result := ShutdownWebServer()

	var response Response
	err := json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}
}

// TestServerLifecycle tests the complete server lifecycle
func TestServerLifecycle(t *testing.T) {
	// Reset server state
	mu.Lock()
	isRunning = false
	isPaused = false
	server = nil
	mu.Unlock()

	// Start
	result := StartWebServer(8080)
	var response Response
	json.Unmarshal([]byte(result), &response)
	if response.Status != "success" {
		t.Error("Failed to start server")
	}

	// Pause
	result = PauseWebServer()
	json.Unmarshal([]byte(result), &response)
	if response.Status != "success" {
		t.Error("Failed to pause server")
	}

	// Resume
	result = ResumeWebServer()
	json.Unmarshal([]byte(result), &response)
	if response.Status != "success" {
		t.Error("Failed to resume server")
	}

	// Shutdown
	result = ShutdownWebServer()
	json.Unmarshal([]byte(result), &response)
	if response.Status != "success" {
		t.Error("Failed to shutdown server")
	}

	// Verify final state
	mu.Lock()
	if isRunning {
		t.Error("Expected server to be stopped after shutdown")
	}
	mu.Unlock()
}
