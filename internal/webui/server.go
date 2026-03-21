package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"regexp"
	"strings"

	"github.com/charlesnobrega/STARDF-Anime/internal/anilist"
	"github.com/charlesnobrega/STARDF-Anime/internal/player"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"embed"
	"io/fs"
)

//go:embed all:static
var staticContent embed.FS

type UnifiedMedia struct {
	Name          string         `json:"Name"`
	ImageURL      string         `json:"ImageURL"`
	TotalEpisodes int            `json:"TotalEpisodes"`
	MediaType     string         `json:"MediaType"`
	Sources       []SourceDetail `json:"Sources"`
}

type SourceDetail struct {
	Name      string `json:"Name"`
	URL       string `json:"URL"`
	AnimeName string `json:"AnimeName,omitempty"`
}


// StartWebUI starts the local web server and opens the browser
func StartWebUI(port int) error {
	mux := http.NewServeMux()

	// Static files (HTML, CSS, JS) - Embedded for standalone portability
	staticFS, _ := fs.Sub(staticContent, "static")
	fileServer := http.FileServer(http.FS(staticFS))
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
	mediaType := r.URL.Query().Get("type")

	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	util.GlobalMediaType = mediaType
	
	// 1. Get Canonical Results from AniList First
	aniClient := anilist.NewClient()
	aniResults, _ := aniClient.SearchAnimes(query, 1, 15) // Use SearchAnimes (plural)

	// 2. Get Scraper Results
	scraperManager := scraper.NewScraperManager()
	scraperResults, err := scraperManager.SearchAnime(query, nil)
	if err != nil {
		util.Errorf("Web API Scraper Error: %v", err)
	}

	// 3. Grouping Logic (Canonical First)
	grouped := make(map[string]*UnifiedMedia)
	aniMap := make(map[string]anilist.MediaSearchResult)

	// Create entries for each AniList result
	for _, ani := range aniResults {
		name := ani.Title.English
		if name == "" {
			name = ani.Title.Romaji
		}
		
		grouped[name] = &UnifiedMedia{
			Name:          name,
			ImageURL:      ani.CoverImage.Large,
			TotalEpisodes: 0,
			MediaType:     "anime",
			Sources:       []SourceDetail{},
		}
		if ani.Episodes != nil {
			grouped[name].TotalEpisodes = *ani.Episodes
		}
		aniMap[name] = ani
	}

	// Match Scraper results to AniList entries
	reTag := regexp.MustCompile(`^\[.*?\]\s*`)
	for _, res := range scraperResults {
		cleanName := strings.ToLower(reTag.ReplaceAllString(res.Name, ""))
		
		matched := false
		for aniName, media := range grouped {
			aniInfo := aniMap[aniName]
			titles := []string{
				strings.ToLower(aniInfo.Title.English),
				strings.ToLower(aniInfo.Title.Romaji),
				strings.ToLower(aniInfo.Title.Native),
			}

			// 1. Precise check: any title contains or is contained in scrap name
			for _, t := range titles {
				if t == "" { continue }
				if strings.Contains(cleanName, t) || strings.Contains(t, cleanName) {
					// CHECK FOR DUPLICATE (Normalizing URLs)
					exists := false
					normResURL := strings.TrimRight(strings.TrimPrefix(strings.TrimPrefix(res.URL, "https://"), "http://"), "/")
					for _, s := range media.Sources {
						normSrcURL := strings.TrimRight(strings.TrimPrefix(strings.TrimPrefix(s.URL, "https://"), "http://"), "/")
						if s.Name == res.Source && normSrcURL == normResURL {
							exists = true
							break
						}
					}
					if !exists {
						media.Sources = append(media.Sources, SourceDetail{
							Name:      res.Source,
							URL:       res.URL,
							AnimeName: res.Name, // Preenchendo o nome real que o scraper achou
						})
					}
					matched = true
					break
				}
			}
			if matched { break }

			// 2. Fragment match: check for seasonal arcs (Shimetsu Kaiyuu, etc)
			// If two words like "Jujutsu" and "Kaisen" match, we accept it for JJK card
			tokens := strings.Fields(cleanName)
			matchesCount := 0
			for _, token := range tokens {
				if len(token) < 4 { continue } // Ignore short words
				for _, t := range titles {
					if strings.Contains(t, token) {
						matchesCount++
						break
					}
				}
			}
			if matchesCount >= 2 {
				media.Sources = append(media.Sources, SourceDetail{
					Name: res.Source,
					URL:  res.URL,
				})
				matched = true
				break
			}
		}

		// If no match found in AniList, only add it if we have no AniList results
		if !matched && len(aniResults) == 0 {
			formattedName := reTag.ReplaceAllString(res.Name, "")
			if _, exists := grouped[formattedName]; !exists {
				grouped[formattedName] = &UnifiedMedia{
					Name:          formattedName,
					ImageURL:      res.ImageURL,
					TotalEpisodes: res.TotalEpisodes,
					MediaType:     string(res.MediaType),
					Sources:       []SourceDetail{},
				}
			}
			grouped[formattedName].Sources = append(grouped[formattedName].Sources, SourceDetail{
				Name: res.Source,
				URL:  res.URL,
			})
		}
	}

	var finalResults []*UnifiedMedia
	for _, v := range grouped {
		// Only show cards that have at least one source, OR are from AniList (allows click-to-search)
		if len(v.Sources) > 0 {
			finalResults = append(finalResults, v)
		}
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

	// FIX: If source is from AniList Dashboard, we need to find a real scraper for it
	if strings.Contains(sourceName, "AniList") {
		util.Infof("Dashboard click detected: searching sources for '%s'", animeURL)
		results, err := scraperManager.SearchAnime(animeURL, nil)
		if err != nil || len(results) == 0 {
			util.Warn("No sources found for dashboard item", "anime", animeURL)
			http.Error(w, "No sources found for this anime", http.StatusNotFound)
			return
		}
		// Pick the first/best source (can be improved later)
		animeURL = results[0].URL
		sourceName = results[0].Source
		util.Infof("Selected source for playback: %s", sourceName)
	}

	s, err := scraperManager.FindScraperByName(sourceName)
	if err != nil {
		util.Errorf("Web API: Scraper not found: %s", sourceName)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	episodes, err := s.GetAnimeEpisodes(animeURL)
	if err != nil {
		util.Errorf("Web API: Get Episodes Error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error(), "url": animeURL, "source": sourceName})
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

	urlStr, metadata, err := s.GetStreamURL(episodeURL)
	if err != nil {
		util.Errorf("Web API: Get Stream Error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// FIX: Handle redirect pages in provider links (like TopAnimes /aviso/?url=)
	if strings.Contains(urlStr, "/aviso/") || strings.Contains(urlStr, "url=") {
		u, parseErr := url.Parse(urlStr)
		if parseErr == nil {
			realURL := u.Query().Get("url")
			if realURL != "" {
				if util.IsDebug {
					util.Infof("Redirect detected! Extracting real URL: %s", realURL)
				}
				urlStr = realURL
			}
		}
	}

	response := map[string]interface{}{
		"url":      urlStr,
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

	// WAIT for StartVideo to confirm the player is actually UP and CONNECTED
	socketPath, err := player.StartVideo(streamURL, args)
	if err != nil {
		util.Errorf("Web API: Play Error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("MPV falhou ao abrir: %v", err)})
		return
	}

	util.Infof("Playing launched successfully on socket: %s", socketPath)
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
