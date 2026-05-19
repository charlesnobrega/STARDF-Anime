package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	server    *http.Server
	mu        sync.Mutex
	isRunning bool
	isPaused  bool
)

// StartWebServer starts the local HTTP server on the specified port
// Exposed via GoMobile for Flutter bridge
// Validates: Requirements 2, 3
func StartWebServer(port int) string {
	log.Printf("[Server] StartWebServer called with port: %d", port)

	mu.Lock()
	defer mu.Unlock()

	if isRunning {
		return successResponse(map[string]interface{}{
			"message": "Server already running",
			"port":    port,
		})
	}

	addr := fmt.Sprintf(":%d", port)
	server = &http.Server{
		Addr:         addr,
		Handler:      setupRoutes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("[Server] Starting HTTP server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[Server] Server error: %v", err)
		}
	}()

	// Give server time to start - ensure startup completes in <2 seconds
	time.Sleep(100 * time.Millisecond)
	isRunning = true
	isPaused = false

	return successResponse(map[string]interface{}{
		"message": "Server started successfully",
		"port":    port,
		"url":     fmt.Sprintf("http://localhost:%d", port),
	})
}

// PauseWebServer pauses the web server gracefully
// Exposed via GoMobile for Flutter bridge
// Validates: Requirements 2, 9
func PauseWebServer() string {
	log.Println("[Server] PauseWebServer called")

	mu.Lock()
	defer mu.Unlock()

	if !isRunning || server == nil {
		return successResponse(map[string]interface{}{
			"message": "Server not running",
		})
	}

	// Mark as paused - server remains running but can be paused at application level
	isPaused = true
	log.Println("[Server] Server paused")

	return successResponse(map[string]interface{}{
		"message": "Server paused successfully",
	})
}

// ResumeWebServer resumes the web server
// Exposed via GoMobile for Flutter bridge
// Validates: Requirements 2, 9
func ResumeWebServer() string {
	log.Println("[Server] ResumeWebServer called")

	mu.Lock()
	defer mu.Unlock()

	if !isRunning {
		return errorResponse("SERVER_NOT_RUNNING", "Server is not running", nil)
	}

	isPaused = false
	log.Println("[Server] Server resumed")

	return successResponse(map[string]interface{}{
		"message": "Server resumed successfully",
	})
}

// ShutdownWebServer shuts down the web server gracefully
// Exposed via GoMobile for Flutter bridge
// Validates: Requirements 2, 9
func ShutdownWebServer() string {
	log.Println("[Server] ShutdownWebServer called")

	mu.Lock()
	defer mu.Unlock()

	if !isRunning || server == nil {
		return successResponse(map[string]interface{}{
			"message": "Server not running",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[Server] Error shutting down server: %v", err)
		return errorResponse("SERVER_SHUTDOWN_ERROR", fmt.Sprintf("Failed to shutdown server: %v", err), nil)
	}

	isRunning = false
	isPaused = false
	server = nil

	log.Println("[Server] Server shut down successfully")

	return successResponse(map[string]interface{}{
		"message": "Server shut down successfully",
	})
}

// setupRoutes configures all HTTP routes
func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
	})

	// API endpoints
	mux.HandleFunc("/api/animes", handleGetAnimes)
	mux.HandleFunc("/api/watchlist", handleWatchlist)
	mux.HandleFunc("/api/history", handleHistory)
	mux.HandleFunc("/api/sync", handleSync)

	return mux
}

// handleGetAnimes handles GET /api/animes
func handleGetAnimes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, GetAnimes())
}

// handleWatchlist handles watchlist operations
func handleWatchlist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, successResponse(map[string]interface{}{"items": []interface{}{}}))
	case http.MethodPost:
		animeID := r.FormValue("animeId")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, AddToWatchlist(animeID))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleHistory handles history operations
func handleHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, successResponse(map[string]interface{}{"items": []interface{}{}}))
	case http.MethodPost:
		animeID := r.FormValue("animeId")
		episodeID := r.FormValue("episodeId")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, MarkAsWatched(animeID, episodeID))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSync handles sync operations
func handleSync(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, SyncWithDesktop())
}

// successResponse creates a standard success response
func successResponse(data interface{}) string {
	response := Response{
		Status:    "success",
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}
