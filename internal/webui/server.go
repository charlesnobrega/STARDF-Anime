package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"embed"
	"github.com/charlesnobrega/STARDF-Anime/internal/anilist"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/player"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
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

type AniListHistoryItem struct {
	AniListID     int            `json:"anilistId"`
	Name          string         `json:"name"`
	ImageURL      string         `json:"imageUrl"`
	Status        string         `json:"status"`
	Progress      int            `json:"progress"`
	TotalEpisodes int            `json:"totalEpisodes"`
	MediaType     string         `json:"mediaType"`
	UpdatedAt     string         `json:"updatedAt"`
	Sources       []SourceDetail `json:"sources"`
}

type aniListLoginRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURI  string `json:"redirectUri"`
	Code         string `json:"code"`
}

var (
	reLeadingTags   = regexp.MustCompile(`^\s*\[[^\]]+\]\s*`)
	reSeasonOrPart  = regexp.MustCompile(`(?i)\s*[-:–—]?\s*(?:season|temporada|part|parte)\s*\d+\s*$`)
	reNonAlnumMatch = regexp.MustCompile(`[^a-z0-9]+`)
	reMultiSpace    = regexp.MustCompile(`\s+`)
)

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func normalizeTitleForMatch(title string) string {
	cleaned := strings.ToLower(strings.TrimSpace(title))
	cleaned = reLeadingTags.ReplaceAllString(cleaned, "")
	cleaned = reSeasonOrPart.ReplaceAllString(cleaned, "")
	cleaned = strings.ReplaceAll(cleaned, "&", " and ")
	cleaned = strings.ReplaceAll(cleaned, "-", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")
	cleaned = reNonAlnumMatch.ReplaceAllString(cleaned, " ")
	cleaned = reMultiSpace.ReplaceAllString(cleaned, " ")
	return strings.TrimSpace(cleaned)
}

func titleMatchScore(query, candidate string) int {
	nq := normalizeTitleForMatch(query)
	nc := normalizeTitleForMatch(candidate)
	if nq == "" || nc == "" {
		return 0
	}
	if nq == nc {
		return 1200
	}

	score := 0
	if strings.HasPrefix(nc, nq+" ") {
		score += 280
	}
	if strings.Contains(" "+nc+" ", " "+nq+" ") {
		score += 220
	} else if strings.Contains(nc, nq) || strings.Contains(nq, nc) {
		score += 120
	}

	queryTokens := strings.Fields(nq)
	candTokens := strings.Fields(nc)
	if len(queryTokens) == 0 || len(candTokens) == 0 {
		return score
	}

	candidateSet := make(map[string]struct{}, len(candTokens))
	for _, tok := range candTokens {
		if len(tok) >= 3 {
			candidateSet[tok] = struct{}{}
		}
	}

	matches := 0
	validQueryTokens := 0
	for _, tok := range queryTokens {
		if len(tok) < 3 {
			continue
		}
		validQueryTokens++
		if _, ok := candidateSet[tok]; ok {
			matches++
			score += 130
		}
	}

	if validQueryTokens > 0 && matches == validQueryTokens {
		extraTokens := len(candTokens) - validQueryTokens
		if extraTokens < 0 {
			extraTokens = 0
		}
		score += 220
		score -= extraTokens * 50
	}

	score -= absInt(len(nc)-len(nq)) * 2
	return score
}

func pickBestScraperMatch(query string, results []*models.Anime) (*models.Anime, int, string) {
	if len(results) == 0 {
		return nil, 0, "nenhum resultado de busca"
	}

	candidates := make([]*models.Anime, 0, len(results))
	animeOnlyCandidates := make([]*models.Anime, 0, len(results))
	for _, candidate := range results {
		if candidate == nil {
			continue
		}
		candidates = append(candidates, candidate)
		if candidate.MediaType != models.MediaTypeMovie && candidate.MediaType != models.MediaTypeTV {
			animeOnlyCandidates = append(animeOnlyCandidates, candidate)
		}
	}
	if len(candidates) == 0 {
		return nil, 0, "resultados sem fonte utilizavel"
	}
	if len(animeOnlyCandidates) > 0 {
		candidates = animeOnlyCandidates
	}

	var best *models.Anime
	bestScore := -1
	secondBest := -1

	for _, candidate := range candidates {
		if candidate == nil || strings.TrimSpace(candidate.Name) == "" || strings.TrimSpace(candidate.URL) == "" || strings.TrimSpace(candidate.Source) == "" {
			continue
		}

		score := titleMatchScore(query, candidate.Name)
		if score > bestScore {
			secondBest = bestScore
			best = candidate
			bestScore = score
			continue
		}
		if score > secondBest {
			secondBest = score
		}
	}

	if best == nil {
		return nil, 0, "resultados sem fonte utilizavel"
	}

	tokenCount := len(strings.Fields(normalizeTitleForMatch(query)))
	minScore := 420
	margin := 50
	if tokenCount <= 1 {
		minScore = 260
		margin = 25
	}
	if bestScore < minScore {
		return nil, bestScore, "score insuficiente para match confiavel"
	}
	if bestScore < 900 && secondBest >= 0 && bestScore-secondBest < margin {
		return nil, bestScore, "resultado ambiguo (scores muito proximos)"
	}

	return best, bestScore, ""
}

func pickAniListTitle(english, romaji string, mediaID int) string {
	if strings.TrimSpace(english) != "" {
		return strings.TrimSpace(english)
	}
	if strings.TrimSpace(romaji) != "" {
		return strings.TrimSpace(romaji)
	}
	return fmt.Sprintf("AniList #%d", mediaID)
}

func ensureAniListViewer(session *anilist.AniListSession) (*anilist.User, error) {
	if !session.IsLoggedIn() {
		return nil, anilist.ErrNotAuthenticated
	}
	if session.CurrentUser != nil {
		return session.CurrentUser, nil
	}
	user, err := session.Client.GetViewer()
	if err != nil {
		return nil, err
	}
	session.CurrentUser = user
	return user, nil
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
	mux.HandleFunc("/api/anilist/status", handleAniListStatus)
	mux.HandleFunc("/api/anilist/sync", handleAniListSync)
	mux.HandleFunc("/api/anilist/logout", handleAniListLogout)
	mux.HandleFunc("/api/anilist/auth-url", handleAniListAuthURL)
	mux.HandleFunc("/api/anilist/login", handleAniListLogin)

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

func handleAniListAuthURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientID := strings.TrimSpace(r.URL.Query().Get("client_id"))
	redirectURI := strings.TrimSpace(r.URL.Query().Get("redirect_uri"))
	if redirectURI == "" {
		redirectURI = "https://anilist.co/api/v2/oauth/pin"
	}

	w.Header().Set("Content-Type", "application/json")
	if clientID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "client_id is required",
		})
		return
	}

	authURL := anilist.GetAuthorizationURL(clientID, redirectURI)
	json.NewEncoder(w).Encode(map[string]string{
		"authUrl":     authURL,
		"redirectUri": redirectURI,
	})
}

func handleAniListLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req aniListLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	req.ClientID = strings.TrimSpace(req.ClientID)
	req.ClientSecret = strings.TrimSpace(req.ClientSecret)
	req.Code = strings.TrimSpace(req.Code)
	req.RedirectURI = strings.TrimSpace(req.RedirectURI)
	if req.RedirectURI == "" {
		req.RedirectURI = "https://anilist.co/api/v2/oauth/pin"
	}

	if req.ClientID == "" || req.Code == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "clientId and code are required",
		})
		return
	}

	if err := anilist.GlobalSession.LoginWithCode(anilist.OAuthConfig{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		RedirectURI:  req.RedirectURI,
	}, req.Code); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Falha no login AniList: %v", err),
		})
		return
	}

	session := anilist.GlobalSession
	if session.CurrentUser == nil && session.IsLoggedIn() {
		if user, err := ensureAniListViewer(session); err == nil {
			session.CurrentUser = user
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":       true,
		"loggedIn": session.IsLoggedIn(),
		"user":     session.CurrentUser,
	})
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		query = strings.TrimSpace(r.URL.Query().Get("query"))
	}
	if query == "" {
		query = strings.TrimSpace(r.URL.Query().Get("name"))
	}
	mediaType := r.URL.Query().Get("type")

	if query == "" {
		http.Error(w, "Query parameter 'q' (or 'query') is required", http.StatusBadRequest)
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
				if t == "" {
					continue
				}
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
			if matched {
				break
			}

			// 2. Fragment match: check for seasonal arcs (Shimetsu Kaiyuu, etc)
			// If two words like "Jujutsu" and "Kaisen" match, we accept it for JJK card
			tokens := strings.Fields(cleanName)
			matchesCount := 0
			for _, token := range tokens {
				if len(token) < 4 {
					continue
				} // Ignore short words
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

func handleAniListStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	session := anilist.GlobalSession
	if !session.IsLoggedIn() {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"loggedIn": false,
		})
		return
	}

	user, err := ensureAniListViewer(session)
	if err != nil {
		util.Warnf("AniList status check failed: %v", err)
		_ = session.Logout()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"loggedIn": false,
			"error":    "Sua sessao AniList expirou. Faca login novamente no app principal.",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"loggedIn": true,
		"user":     user,
	})
}

func handleAniListLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := anilist.GlobalSession.Logout(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Falha ao desconectar AniList: %v", err),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"ok": true,
	})
}

func handleAniListSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	session := anilist.GlobalSession
	if !session.IsLoggedIn() {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"loggedIn": false,
			"error":    "AniList nao conectado",
		})
		return
	}

	user, err := ensureAniListViewer(session)
	if err != nil {
		util.Warnf("AniList sync viewer check failed: %v", err)
		_ = session.Logout()
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"loggedIn": false,
			"error":    "Sessao AniList expirada. Faca login novamente no app principal.",
		})
		return
	}

	list, err := session.Client.GetUserList(user.ID)
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "401") || strings.Contains(errMsg, "unauthorized") {
			_ = session.Logout()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"loggedIn": false,
				"error":    "Token AniList invalido/expirado. Faca login novamente no app principal.",
			})
			return
		}

		util.Errorf("AniList sync error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"loggedIn": true,
			"error":    fmt.Sprintf("Falha ao sincronizar AniList: %v", err),
		})
		return
	}

	type anilistEntry struct {
		aniListID     int
		name          string
		imageURL      string
		status        string
		progress      int
		totalEpisodes int
		updatedAt     int
	}

	latestByMedia := map[int]anilistEntry{}
	for _, listGroup := range list.Data.MediaListCollection.Lists {
		for _, entry := range listGroup.Entries {
			if entry.MediaID == 0 {
				continue
			}

			totalEps := 0
			if entry.Media.Episodes != nil {
				totalEps = *entry.Media.Episodes
			}

			current := anilistEntry{
				aniListID:     entry.MediaID,
				name:          pickAniListTitle(entry.Media.Title.English, entry.Media.Title.Romaji, entry.MediaID),
				imageURL:      entry.Media.CoverImage.Large,
				status:        strings.ToUpper(strings.TrimSpace(entry.Status)),
				progress:      entry.Progress,
				totalEpisodes: totalEps,
				updatedAt:     entry.UpdatedAt,
			}

			prev, exists := latestByMedia[entry.MediaID]
			if !exists || current.updatedAt >= prev.updatedAt {
				latestByMedia[entry.MediaID] = current
			}
		}
	}

	entries := make([]anilistEntry, 0, len(latestByMedia))
	for _, e := range latestByMedia {
		entries = append(entries, e)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].updatedAt > entries[j].updatedAt
	})

	counts := map[string]int{
		"total":     0,
		"CURRENT":   0,
		"COMPLETED": 0,
		"PLANNING":  0,
		"DROPPED":   0,
		"PAUSED":    0,
	}

	suggestions := make([]*UnifiedMedia, 0, 12)
	history := make([]AniListHistoryItem, 0, 24)

	for _, e := range entries {
		counts["total"]++
		if _, ok := counts[e.status]; ok {
			counts[e.status]++
		}

		defaultSource := []SourceDetail{{Name: "AniList (Sync)", URL: e.name}}

		if (e.status == "CURRENT" || e.status == "PLANNING" || e.status == "PAUSED") && len(suggestions) < 12 {
			suggestions = append(suggestions, &UnifiedMedia{
				Name:          e.name,
				ImageURL:      e.imageURL,
				TotalEpisodes: e.totalEpisodes,
				MediaType:     "anime",
				Sources:       defaultSource,
			})
		}

		if len(history) < 24 {
			updatedAtISO := ""
			if e.updatedAt > 0 {
				updatedAtISO = time.Unix(int64(e.updatedAt), 0).Format(time.RFC3339)
			}
			history = append(history, AniListHistoryItem{
				AniListID:     e.aniListID,
				Name:          e.name,
				ImageURL:      e.imageURL,
				Status:        e.status,
				Progress:      e.progress,
				TotalEpisodes: e.totalEpisodes,
				MediaType:     "anime",
				UpdatedAt:     updatedAtISO,
				Sources:       defaultSource,
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"loggedIn":    true,
		"user":        user,
		"counts":      counts,
		"suggestions": suggestions,
		"history":     history,
	})
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Nenhuma fonte encontrada para esta obra.",
				"query": animeURL,
			})
			return
		}

		best, bestScore, reason := pickBestScraperMatch(animeURL, results)
		if best == nil {
			util.Warn("AniList resolver could not find a reliable match", "query", animeURL, "reason", reason, "score", bestScore)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Nao foi possivel selecionar fonte confiavel para '%s' (%s).", animeURL, reason),
				"query": animeURL,
			})
			return
		}

		animeURL = best.URL
		sourceName = best.Source
		util.Infof("Selected source for playback: %s (score=%d, title=%s)", sourceName, bestScore, best.Name)
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
