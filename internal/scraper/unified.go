// Package scraper provides a unified interface for different anime sources
package scraper

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"encoding/json"

	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"os"
	"regexp"
)

// ScraperType represents different scraper types
type ScraperType int

const (
	// searchTimeout is the maximum time to wait for all scrapers
	searchTimeout = 12 * time.Second
	// perScraperTimeout is the timeout for individual scrapers
	perScraperTimeout = 10 * time.Second
	// earlyReturnDelay is the time to wait after first results before returning
	earlyReturnDelay = 3000 * time.Millisecond
	// minResultsForEarlyReturn is the minimum results needed to trigger early return
	minResultsForEarlyReturn = 10
)

// ErrBackRequested is returned when the user requests to go back to the previous menu
var ErrBackRequested = fmt.Errorf("back requested")

const (
	AnimefireType ScraperType = iota
	FlixHQType
	CinebyType
	AnimesOnlineCCTYPE
	GoyabuType
	CineGratisType
	BetterAnimeType
	TopAnimesType
	AnimesDigitalType
	Dynamic1Type
	Dynamic2Type
	Dynamic3Type
	Dynamic4Type
	Dynamic5Type
	Dynamic6Type
	Dynamic7Type
	Dynamic8Type
	Dynamic9Type
	Dynamic10Type
)

// UnifiedScraper provides a common interface for all scrapers
type UnifiedScraper interface {
	SearchAnime(query string, options ...interface{}) ([]*models.Anime, error)
	GetAnimeEpisodes(animeURL string) ([]models.Episode, error)
	GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error)
	GetType() ScraperType
}

// ScraperManager manages multiple scrapers
type ScraperManager struct {
	scrapers map[ScraperType]UnifiedScraper
}

// NewScraperManager creates a new scraper manager
func NewScraperManager() *ScraperManager {
	manager := &ScraperManager{
		scrapers: make(map[ScraperType]UnifiedScraper),
	}

	// Register Static Scrapers (Special Logic)
	manager.scrapers[FlixHQType] = &FlixHQAdapter{client: NewFlixHQClient()}
	manager.scrapers[AnimefireType] = &AnimefireAdapter{client: NewAnimefireClient()}
	manager.scrapers[AnimesOnlineCCTYPE] = &AnimesOnlineCCAdapter{client: NewAnimesOnlineCCClient()}
	manager.scrapers[BetterAnimeType] = &BetterAnimeAdapter{client: NewBetterAnimeClient()}
	manager.scrapers[TopAnimesType] = &TopAnimesAdapter{client: NewTopAnimesClient()}
	manager.scrapers[AnimesDigitalType] = &AnimesDigitalAdapter{client: NewAnimesDigitalClient()}

	// Load Dynamic Scrapers (JSON Manifest from Spider)
	manifestURL := "http://localhost:3000/scrapers.json"
	dynamicConfigs, err := LoadDynamicScrapers(manifestURL)
	if err != nil {
		util.Debug("Spider URL failed, trying local scrapers.json...", "error", err)
		// Try local file fallback
		localPath := "scrapers.json"
		if _, statErr := os.Stat(localPath); statErr == nil {
			f, _ := os.Open(localPath)
			defer f.Close()
			var manifest DynamicManifest
			if decErr := json.NewDecoder(f).Decode(&manifest); decErr == nil {
				dynamicConfigs = manifest.Scrapers
				err = nil // Cleared error
			}
		}
	}

	if err == nil {
		dynamicTypes := []ScraperType{
			Dynamic1Type, Dynamic2Type, Dynamic3Type, Dynamic4Type, Dynamic5Type,
			Dynamic6Type, Dynamic7Type, Dynamic8Type, Dynamic9Type, Dynamic10Type,
		}
		for i, cfg := range dynamicConfigs {
			if i >= len(dynamicTypes) {
				break
			}
			st := dynamicTypes[i]
			manager.scrapers[st] = NewDynamicScraper(cfg, st)
			util.Debug("Dynamic source registered", "name", cfg.Name, "type", st)
		}
	} else {
		util.Error("Failed to load any dynamic scrapers", "error", err)
	}

	return manager
}

// GetScraper returns a specific scraper by type
func (sm *ScraperManager) GetScraper(scraperType ScraperType) (UnifiedScraper, error) {
	if s, exists := sm.scrapers[scraperType]; exists {
		return s, nil
	}
	return nil, fmt.Errorf("scraper %v not found", scraperType)
}

// FindScraperByName finds a scraper by its display name
func (sm *ScraperManager) FindScraperByName(name string) (UnifiedScraper, error) {
	for st, s := range sm.scrapers {
		if sm.getScraperDisplayName(st) == name {
			return s, nil
		}
	}
	// Try fuzzy match
	for st, s := range sm.scrapers {
		if strings.Contains(strings.ToLower(sm.getScraperDisplayName(st)), strings.ToLower(name)) {
			return s, nil
		}
	}
	return nil, fmt.Errorf("scraper %s not found", name)
}

func (sm *ScraperManager) SearchAnime(query string, scraperType *ScraperType) ([]*models.Anime, error) {
	if scraperType != nil {
		return sm.searchSpecificScraper(query, *scraperType)
	}
	return sm.searchAllScrapersConcurrent(query)
}

func (sm *ScraperManager) searchSpecificScraper(query string, scraperType ScraperType) ([]*models.Anime, error) {
	scraper, exists := sm.scrapers[scraperType]
	if !exists {
		return nil, fmt.Errorf("scraper type %v not found", scraperType)
	}

	results, err := scraper.SearchAnime(query)
	if err != nil {
		return nil, err
	}

	sm.tagResults(results, scraperType)
	return results, nil
}

type searchResult struct {
	scraperType ScraperType
	results     []*models.Anime
	err         error
}

func (sm *ScraperManager) searchAllScrapersConcurrent(query string) ([]*models.Anime, error) {
	ctx, cancel := context.WithTimeout(context.Background(), searchTimeout)
	defer cancel()

	var (
		allResults   []*models.Anime
		resultsMutex sync.Mutex
		searchErrors []string
		errorsMutex  sync.Mutex
	)

	var (
		completedCount  int32
		totalScrapers   = int32(len(sm.scrapers))
		firstResultOnce sync.Once
	)

	resultChan := make(chan searchResult, len(sm.scrapers))
	var wg sync.WaitGroup

	for sType, s := range sm.scrapers {
		wg.Add(1)
		go func(st ScraperType, scr UnifiedScraper) {
			defer wg.Done()
			defer atomic.AddInt32(&completedCount, 1)

			result := sm.searchWithTimeout(ctx, st, scr, query)
			resultChan <- result
		}(sType, s)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var earlyReturnTimer <-chan time.Time

	for {
		select {
		case res, ok := <-resultChan:
			if !ok {
				goto done
			}

			if res.err != nil {
				errorsMutex.Lock()
				searchErrors = append(searchErrors, fmt.Sprintf("%s: %v", sm.getScraperDisplayName(res.scraperType), res.err))
				errorsMutex.Unlock()
				continue
			}

			if len(res.results) > 0 {
				sm.tagResults(res.results, res.scraperType)
				resultsMutex.Lock()
				allResults = append(allResults, res.results...)
				currentCount := len(allResults)
				resultsMutex.Unlock()

				firstResultOnce.Do(func() {
					earlyReturnTimer = time.After(earlyReturnDelay)
				})

				if currentCount >= minResultsForEarlyReturn && atomic.LoadInt32(&completedCount) < totalScrapers {
					// Check if timer already fired or wait for it
				}
			}

		case <-earlyReturnTimer:
			goto done
		case <-ctx.Done():
			goto done
		}
	}

done:
	if len(allResults) == 0 && len(searchErrors) > 0 {
		return nil, fmt.Errorf("no results found (errors: %s)", strings.Join(searchErrors, "; "))
	}
	return allResults, nil
}

func (sm *ScraperManager) searchWithTimeout(ctx context.Context, st ScraperType, s UnifiedScraper, query string) searchResult {
	sourceName := sm.getScraperDisplayName(st)
	scraperCtx, scraperCancel := context.WithTimeout(ctx, perScraperTimeout)
	defer scraperCancel()

	done := make(chan searchResult, 1)
	go func() {
		results, err := s.SearchAnime(query)
		done <- searchResult{scraperType: st, results: results, err: err}
	}()

	select {
	case result := <-done:
		tracker := tracking.GetGlobalTracker()
		if tracker != nil {
			_ = tracker.TrackScraperAction(sourceName, result.err == nil, "")
		}
		return result
	case <-scraperCtx.Done():
		return searchResult{scraperType: st, err: fmt.Errorf("timeout")}
	}
}

func (sm *ScraperManager) tagResults(results []*models.Anime, scraperType ScraperType) {
	sourceName := sm.getScraperDisplayName(scraperType)
	tag := sm.getLanguageTag(scraperType)
	for _, anime := range results {
		if !strings.HasPrefix(anime.Name, "[") {
			anime.Name = fmt.Sprintf("%s %s", tag, anime.Name)
		}
		anime.Source = sourceName
	}
}

func (sm *ScraperManager) getScraperDisplayName(scraperType ScraperType) string {
	if scraperType == FlixHQType {
		return "FlixHQ"
	}
	if scraperType == AnimefireType {
		return "AnimeFire"
	}
	if scraperType == AnimesOnlineCCTYPE {
		return "AnimesOnlineCC"
	}
	if scraperType == BetterAnimeType {
		return "BetterAnime"
	}
	if scraperType == TopAnimesType {
		return "TopAnimes"
	}
	if scraperType == AnimesDigitalType {
		return "AnimesDigital"
	}
	if s, exists := sm.scrapers[scraperType]; exists {
		if ds, ok := s.(*DynamicScraper); ok {
			return ds.Config.Name
		}
	}
	return "Desconhecido"
}

func (sm *ScraperManager) getLanguageTag(scraperType ScraperType) string {
	if scraperType == FlixHQType {
		return "[Movies/TV]"
	}
	return "[Source]"
}

// FlixHQAdapter adapts FlixHQClient
type FlixHQAdapter struct {
	client *FlixHQClient
}

func (a *FlixHQAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	media, err := a.client.SearchMedia(query)
	if err != nil {
		return nil, err
	}
	var animes []*models.Anime
	for _, m := range media {
		anime := m.ToAnimeModel()
		if m.Type == MediaTypeMovie {
			anime.MediaType = models.MediaTypeMovie
		} else {
			anime.MediaType = models.MediaTypeTV
		}
		animes = append(animes, anime)
	}
	return animes, nil
}

func (a *FlixHQAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	// Extract ID from URL
	re := regexp.MustCompile(`-(\d+)$`)
	matches := re.FindStringSubmatch(animeURL)
	if len(matches) < 2 {
		return nil, fmt.Errorf("could not extract media ID from URL")
	}
	mediaID := matches[1]

	if strings.Contains(animeURL, "/movie/") {
		// Movies are a single episode
		return []models.Episode{
			{Number: "1", Num: 1, URL: mediaID, Title: models.TitleDetails{English: "Movie"}},
		}, nil
	}

	// TV Shows: Get seasons then episodes for season 1
	seasons, err := a.client.GetSeasons(mediaID)
	if err != nil {
		return nil, err
	}
	if len(seasons) == 0 {
		return nil, fmt.Errorf("no seasons found")
	}

	eps, err := a.client.GetEpisodes(seasons[0].ID)
	if err != nil {
		return nil, err
	}

	var results []models.Episode
	for _, e := range eps {
		results = append(results, e.ToEpisodeModel())
	}
	return results, nil
}

func (a *FlixHQAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	var episodeID string
	var err error

	// Distinguish between movieID (passed as episodeURL for movies) and episodeDataID
	if len(episodeURL) > 3 && !strings.Contains(episodeURL, "ajax") {
		// If it's pure ID, it might be a movie or we need to find the server first
		// In adapter, we try to get server first
		episodeID, err = a.client.GetMovieServerID(episodeURL, "UpCloud")
		if err != nil {
			episodeID, err = a.client.GetEpisodeServerID(episodeURL, "Vidcloud")
		}
	} else {
		episodeID = episodeURL
	}

	if err != nil && episodeID == "" {
		episodeID = episodeURL // Fallback
	}

	embedLink, err := a.client.GetEmbedLink(episodeID)
	if err != nil {
		return "", nil, err
	}
	info, err := a.client.ExtractStreamInfo(embedLink, "1080", "english")
	if err != nil {
		return "", nil, err
	}
	return info.VideoURL, map[string]string{"source": "flixhq"}, nil
}

func (a *FlixHQAdapter) GetType() ScraperType {
	return FlixHQType
}

// GetClient returns the underlying FlixHQ client
func (a *FlixHQAdapter) GetClient() *FlixHQClient {
	return a.client
}

// AnimefireAdapter adapts AnimefireClient
type AnimefireAdapter struct {
	client *AnimefireClient
}

func (a *AnimefireAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *AnimefireAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetAnimeEpisodes(animeURL)
}

func (a *AnimefireAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	url, err := a.client.GetEpisodeStreamURL(episodeURL)
	return url, map[string]string{"source": "animefire"}, err
}

func (a *AnimefireAdapter) GetType() ScraperType {
	return AnimefireType
}

// AnimesOnlineCCAdapter adapts AnimesOnlineCCClient
type AnimesOnlineCCAdapter struct {
	client *AnimesOnlineCCClient
}

func (a *AnimesOnlineCCAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *AnimesOnlineCCAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *AnimesOnlineCCAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *AnimesOnlineCCAdapter) GetType() ScraperType {
	return AnimesOnlineCCTYPE
}

// BetterAnimeAdapter adapts BetterAnimeClient
type BetterAnimeAdapter struct {
	client *BetterAnimeClient
}

func (a *BetterAnimeAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *BetterAnimeAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *BetterAnimeAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *BetterAnimeAdapter) GetType() ScraperType {
	return BetterAnimeType
}

// TopAnimesAdapter adapts TopAnimesClient
type TopAnimesAdapter struct {
	client *TopAnimesClient
}

func (a *TopAnimesAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *TopAnimesAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *TopAnimesAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *TopAnimesAdapter) GetType() ScraperType {
	return TopAnimesType
}

// AnimesDigitalAdapter adapts AnimesDigitalClient
type AnimesDigitalAdapter struct {
	client *AnimesDigitalClient
}

func (a *AnimesDigitalAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *AnimesDigitalAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *AnimesDigitalAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *AnimesDigitalAdapter) GetType() ScraperType {
	return AnimesDigitalType
}

