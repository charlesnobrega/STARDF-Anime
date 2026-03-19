package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

// StartWebUI starts the local web server and opens the browser
func StartWebUI(port int) error {
	mux := http.NewServeMux()

	// Static files (HTML, CSS, JS)
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/", fileServer)

	// API Endpoints
	mux.HandleFunc("/api/search", handleSearch)
	mux.HandleFunc("/api/episodes", handleGetEpisodes)
	mux.HandleFunc("/api/stream", handleGetStream)

	addr := fmt.Sprintf("localhost:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	util.Infof("Iniciando StarDF-Anime Web UI em %s", url)
	
	// Open browser in a separate goroutine
	go func() {
		// Wait a small bit for server to be up
		util.Debug("Aguardando servidor para abrir navegador...")
		openBrowser(url)
	}()

	return http.ListenAndServe(addr, mux)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	mediaType := r.URL.Query().Get("type") // "anime" or "movie"

	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Set global context for scrapers
	util.GlobalMediaType = mediaType
	
	scraperManager := scraper.NewScraperManager()
	results, err := scraperManager.SearchAnime(query, nil)
	if err != nil {
		util.Errorf("Web API Search Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func handleGetEpisodes(w http.ResponseWriter, r *http.Request) {
	animeURL := r.URL.Query().Get("url")
	sourceName := r.URL.Query().Get("source")

	if animeURL == "" || sourceName == "" {
		http.Error(w, "Parameters 'url' and 'source' are required", http.StatusBadRequest)
		return
	}

	scraperManager := scraper.NewScraperManager()
	s, err := scraperManager.FindScraperByName(sourceName)
	if err != nil {
		util.Errorf("Web API: Scraper not found: %s", sourceName)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	episodes, err := s.GetAnimeEpisodes(animeURL)
	if err != nil {
		util.Errorf("Web API: Get Episodes Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
}

func handleGetStream(w http.ResponseWriter, r *http.Request) {
	episodeURL := r.URL.Query().Get("url")
	sourceName := r.URL.Query().Get("source")

	if episodeURL == "" || sourceName == "" {
		http.Error(w, "Parameters 'url' and 'source' are required", http.StatusBadRequest)
		return
	}

	scraperManager := scraper.NewScraperManager()
	s, err := scraperManager.FindScraperByName(sourceName)
	if err != nil {
		util.Errorf("Web API: Scraper not found: %s", sourceName)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	url, metadata, err := s.GetStreamURL(episodeURL)
	if err != nil {
		util.Errorf("Web API: Get Stream Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"url":      url,
		"metadata": metadata,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// openBrowser opens the specified URL in the default browser of the user.
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		util.Errorf("Erro ao abrir navegador: %v", err)
	}
}
