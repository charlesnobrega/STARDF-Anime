package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/charlesnobrega/STARDF-Anime/internal/anilist"
	"github.com/charlesnobrega/STARDF-Anime/internal/player"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"regexp"
)

type UnifiedMedia struct {
	Name          string         `json:"Name"`
	ImageURL      string         `json:"ImageURL"`
	TotalEpisodes int            `json:"TotalEpisodes"`
	MediaType     string         `json:"MediaType"`
	Sources       []SourceDetail `json:"Sources"`
}

type SourceDetail struct {
	Name string `json:"Name"`
	URL  string `json:"URL"`
}


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
	mux.HandleFunc("/api/chat", handleChat)
	mux.HandleFunc("/api/trending", handleGetTrending)
	mux.HandleFunc("/api/play", handlePlay)

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

	util.GlobalMediaType = mediaType
	scraperManager := scraper.NewScraperManager()
	results, err := scraperManager.SearchAnime(query, nil)
	if err != nil {
		util.Errorf("Web API Search Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Grouping Logic
	grouped := make(map[string]*UnifiedMedia)
	reTag := regexp.MustCompile(`^\[.*?\]\s*`)

	for _, res := range results {
		cleanName := reTag.ReplaceAllString(res.Name, "")
		
		if _, exists := grouped[cleanName]; !exists {
			grouped[cleanName] = &UnifiedMedia{
				Name:          cleanName,
				ImageURL:      res.ImageURL,
				TotalEpisodes: res.TotalEpisodes,
				MediaType:     string(res.MediaType),
				Sources:       []SourceDetail{},
			}
		}
		
		grouped[cleanName].Sources = append(grouped[cleanName].Sources, SourceDetail{
			Name: res.Source,
			URL:  res.URL,
		})
		
		// Fill missing images/episodes if empty
		if grouped[cleanName].ImageURL == "" && res.ImageURL != "" {
			grouped[cleanName].ImageURL = res.ImageURL
		}
		if grouped[cleanName].TotalEpisodes == 0 && res.TotalEpisodes > 0 {
			grouped[cleanName].TotalEpisodes = res.TotalEpisodes
		}
	}

	var finalResults []*UnifiedMedia
	for _, v := range grouped {
		finalResults = append(finalResults, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalResults)
}

func handleGetTrending(w http.ResponseWriter, r *http.Request) {
	client := anilist.NewClient()
	results, err := client.GetTrendingSeason(1)
	if err != nil {
		util.Errorf("Web API Trending Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var finalResults []*UnifiedMedia
	for _, res := range results {
		name := res.Title.Romaji
		if res.Title.English != "" {
			name = res.Title.English
		}
		
		totalEps := 0
		if res.Episodes != nil {
			totalEps = *res.Episodes
		}

		finalResults = append(finalResults, &UnifiedMedia{
			Name:          name,
			ImageURL:      res.CoverImage.Large,
			TotalEpisodes: totalEps,
			MediaType:     "anime",
			Sources:       []SourceDetail{{Name: "AniList (Sync)", URL: name}},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalResults)
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

func handlePlay(w http.ResponseWriter, r *http.Request) {
	streamURL := r.URL.Query().Get("url")
	referer := r.URL.Query().Get("referer")
	title := r.URL.Query().Get("title")

	if streamURL == "" {
		http.Error(w, "Parameter 'url' is required", http.StatusBadRequest)
		return
	}

	args := []string{
		"--hwdec=auto",
		fmt.Sprintf("--title=%s", title),
	}

	if referer != "" {
		args = append(args, fmt.Sprintf("--http-header-fields=Referer: %s", referer))
	}

	args = append(args, fmt.Sprintf("--user-agent=%s", util.UserAgentList()))

	go func() {
		_, err := player.StartVideo(streamURL, args)
		if err != nil {
			util.Errorf("Web API: Play Error: %v", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"Playing launched"}`))
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var msg struct {
			User string `json:"user"`
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		util.Debug("Chat Message Received:", "user", msg.User, "text", msg.Text)
		w.WriteHeader(http.StatusCreated)
		return
	}

	mockMsgs := []map[string]string{
		{"user": "Admin", "text": "Bem-vindo ao Cloud Chat!"},
		{"user": "Hacker", "text": "Alguém viu o ep 5?"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockMsgs)
}

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
		util.Errorf("Failed to open browser: %v", err)
	}
}
