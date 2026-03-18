// Package scraper provides a unified interface for different anime sources
package scraper

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

// ScraperType represents different scraper types
type ScraperType int

// Timeout configurations - balanced for multiple sources
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
	SuperAnimesType
	CineGratisType
	BetterAnimeType
	TopAnimesType
	AnimesDigitalType
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

	manager.scrapers[AnimefireType] = &AnimefireAdapter{client: NewAnimefireClient()}
	manager.scrapers[BetterAnimeType] = &BetterAnimeAdapter{client: NewBetterAnimeClient()}
	manager.scrapers[TopAnimesType] = &TopAnimesAdapter{client: NewTopAnimesClient()}
	manager.scrapers[AnimesDigitalType] = &AnimesDigitalAdapter{client: NewAnimesDigitalClient()}
	manager.scrapers[CinebyType] = &CinebyAdapter{client: NewCinebyClient()}
	manager.scrapers[CineGratisType] = &CineGratisAdapter{client: NewCineGratisClient()}
	manager.scrapers[FlixHQType] = &FlixHQAdapter{client: NewFlixHQClient()}
	
	manager.scrapers[GoyabuType] = &GoyabuAdapter{client: NewGoyabuClient()}
	// manager.scrapers[SuperAnimesType] = &SuperAnimesAdapter{client: NewSuperAnimesClient()}
	manager.scrapers[AnimesOnlineCCTYPE] = &AnimesOnlineCCAdapter{client: NewAnimesOnlineCCClient()}

	return manager
}

// SearchAnime searches across all available scrapers with enhanced Portuguese messaging
// Uses optimized goroutines with early return for better performance
func (sm *ScraperManager) SearchAnime(query string, scraperType *ScraperType) ([]*models.Anime, error) {
	timer := util.StartTimer("SearchAnime:Total")
	defer timer.Stop()

	util.PerfCount("search_requests")

	if scraperType != nil {
		return sm.searchSpecificScraper(query, *scraperType)
	}

	return sm.searchAllScrapersConcurrent(query)
}

// searchSpecificScraper searches using a single specific scraper
func (sm *ScraperManager) searchSpecificScraper(query string, scraperType ScraperType) ([]*models.Anime, error) {
	scraper, exists := sm.scrapers[scraperType]
	if !exists {
		return nil, fmt.Errorf("scraper type %v not found", scraperType)
	}

	util.Debug("Searching specific scraper", "scraper", sm.getScraperDisplayName(scraperType))

	results, err := scraper.SearchAnime(query)
	if err != nil {
		return nil, fmt.Errorf("busca falhou em %s: %w", sm.getScraperDisplayName(scraperType), err)
	}

	// Add language tags
	sm.tagResults(results, scraperType)

	if len(results) > 0 {
		util.Debug("Search completed", "scraper", sm.getScraperDisplayName(scraperType), "results", len(results))
	}

	return results, nil
}

// searchResult holds the result from a single scraper goroutine
type searchResult struct {
	scraperType ScraperType
	results     []*models.Anime
	err         error
}

// searchAllScrapersConcurrent searches all scrapers in parallel with early return optimization
func (sm *ScraperManager) searchAllScrapersConcurrent(query string) ([]*models.Anime, error) {
	util.Debug("Starting concurrent search across all sources", "query", query)

	ctx, cancel := context.WithTimeout(context.Background(), searchTimeout)
	defer cancel()

	// Thread-safe result collection
	var (
		allResults   []*models.Anime
		resultsMutex sync.Mutex
		searchErrors []string
		errorsMutex  sync.Mutex
	)

	// Track completion for early return
	var (
		completedCount  int32
		totalScrapers   = int32(len(sm.scrapers)) // #nosec G115 - scrapers count is always small (<10)
		firstResultTime time.Time
		firstResultOnce sync.Once
	)

	resultChan := make(chan searchResult, len(sm.scrapers))
	var wg sync.WaitGroup

	// Launch relevant scrapers concurrently based on category
	for sType, s := range sm.scrapers {
		// Filter by category if GlobalMediaType is set
		if util.GlobalMediaType == "anime" {
			if sType == CinebyType || sType == FlixHQType || sType == CineGratisType {
				atomic.AddInt32(&completedCount, 1) // Count as done skip
				continue
			}
		} else if util.GlobalMediaType == "movie" {
			if sType == AnimefireType || sType == AnimesOnlineCCTYPE || sType == GoyabuType || sType == SuperAnimesType || sType == BetterAnimeType || sType == TopAnimesType || sType == AnimesDigitalType {
				atomic.AddInt32(&completedCount, 1) // Count as done skip
				continue
			}
		}

		wg.Add(1)
		go func(st ScraperType, scr UnifiedScraper) {
			defer wg.Done()
			defer atomic.AddInt32(&completedCount, 1)

			result := sm.searchWithTimeout(ctx, st, scr, query)
			resultChan <- result
		}(sType, s)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Early return timer - starts when we get first good results
	var earlyReturnTimer <-chan time.Time

	// Collect results with early return logic
	for {
		select {
		case res, ok := <-resultChan:
			if !ok {
				// Channel closed, all scrapers done
				goto done
			}

			if res.err != nil {
				errorsMutex.Lock()
				sourceName := sm.getScraperDisplayName(res.scraperType)
				util.Debug("Search error", "source", sourceName, "error", res.err)
				searchErrors = append(searchErrors, fmt.Sprintf("%s: %v", sourceName, res.err))
				errorsMutex.Unlock()
				continue
			}

			if len(res.results) > 0 {
				// Tag and add results thread-safely
				sm.tagResults(res.results, res.scraperType)

				resultsMutex.Lock()
				allResults = append(allResults, res.results...)
				currentCount := len(allResults)
				resultsMutex.Unlock()

				util.Debug("Search results received", "source", sm.getScraperDisplayName(res.scraperType), "count", len(res.results), "total", currentCount)

				// Start early return timer on first results
				firstResultOnce.Do(func() {
					firstResultTime = time.Now()
					earlyReturnTimer = time.After(earlyReturnDelay)
				})
			}

		case <-earlyReturnTimer:
			// Early return: we have results and waited long enough for other sources
			resultsMutex.Lock()
			currentCount := len(allResults)
			resultsMutex.Unlock()

			completed := atomic.LoadInt32(&completedCount)
			if currentCount >= minResultsForEarlyReturn && completed < totalScrapers {
				util.Debug("Early return triggered",
					"results", currentCount,
					"completed", completed,
					"total", totalScrapers,
					"waitTime", time.Since(firstResultTime))
				goto done
			}

		case <-ctx.Done():
			util.Debug("Search timeout reached")
			goto done
		}
	}

done:
	// Log warnings for failed sources
	errorsMutex.Lock()
	if len(searchErrors) > 0 {
		for _, errMsg := range searchErrors {
			util.Warn("Search source unavailable", "details", errMsg)
		}
	}
	errorsMutex.Unlock()

	resultsMutex.Lock()
	finalResults := allResults
	resultsMutex.Unlock()

	if len(finalResults) == 0 {
		util.Debug("No anime found", "query", query)
		errorsMutex.Lock()
		defer errorsMutex.Unlock()
		if len(searchErrors) > 0 {
			return nil, fmt.Errorf("no anime found with name: %s (some sources failed: %s)", query, strings.Join(searchErrors, "; "))
		}
		return nil, fmt.Errorf("no anime found with name: %s", query)
	}

	sm.logSearchSummary(finalResults)
	return finalResults, nil
}

// searchWithTimeout executes a single scraper search with timeout
func (sm *ScraperManager) searchWithTimeout(ctx context.Context, st ScraperType, s UnifiedScraper, query string) searchResult {
	sourceName := sm.getScraperDisplayName(st)
	timer := util.StartTimer("Search:" + sourceName)
	util.Debug("Searching in source", "source", sourceName)

	// Create individual timeout context
	scraperCtx, scraperCancel := context.WithTimeout(ctx, perScraperTimeout)
	defer scraperCancel()

	// Channel for search result
	done := make(chan searchResult, 1)

	go func() {
		results, err := s.SearchAnime(query)
		done <- searchResult{
			scraperType: st,
			results:     results,
			err:         err,
		}
	}()

	// Wait for result or timeout
	select {
	case result := <-done:
		timer.Stop()
		tracker := tracking.GetGlobalTracker()
		if result.err == nil {
			util.PerfCount("search_success:" + sourceName)
			if tracker != nil {
				_ = tracker.TrackScraperAction(sourceName, true, "")
			}
		} else {
			util.PerfCount("search_error:" + sourceName)
			if tracker != nil {
				_ = tracker.TrackScraperAction(sourceName, false, result.err.Error())
			}
		}
		return result
	case <-scraperCtx.Done():
		timer.Stop()
		util.PerfCount("search_timeout:" + sourceName)
		util.Debug("Search timeout", "source", sourceName)
		tracker := tracking.GetGlobalTracker()
		if tracker != nil {
			_ = tracker.TrackScraperAction(sourceName, false, "timeout")
		}
		return searchResult{
			scraperType: st,
			results:     nil,
			err:         fmt.Errorf("search timed out after %v", perScraperTimeout),
		}
	}
}

// tagResults adds language tags and source metadata to results
func (sm *ScraperManager) tagResults(results []*models.Anime, scraperType ScraperType) {
	sourceName := sm.getScraperDisplayName(scraperType)
	languageTag := sm.getLanguageTag(scraperType)

	for _, anime := range results {
		// Check if the anime name already has any language tag
		hasLanguageTag := strings.Contains(anime.Name, "[English]") ||
			strings.Contains(anime.Name, "[Portuguese]") ||
			strings.Contains(anime.Name, "[Português]") ||
			strings.Contains(anime.Name, "[Movies/TV]")

		if !hasLanguageTag {
			anime.Name = fmt.Sprintf("%s %s", languageTag, anime.Name)
		}
		anime.Source = sourceName
	}
}

// logSearchSummary logs a summary of search results by source
func (sm *ScraperManager) logSearchSummary(results []*models.Anime) {
	if !util.IsDebug {
		return
	}

	counts := make(map[string]int)
	for _, anime := range results {
		counts[anime.Source]++
	}

	util.Debug("Search summary",
		"animeFire", counts["Animefire.io"],
		"animesOnlineCC", counts["AnimesOnlineCC"],
		"goyabu", counts["Goyabu"],
		"superAnimes", counts["SuperAnimes"],
		"cineby", counts["Cineby"],
		"flixHQ", counts["FlixHQ"],
		"total", len(results))
}

// GetScraper returns a specific scraper by type
func (sm *ScraperManager) GetScraper(scraperType ScraperType) (UnifiedScraper, error) {
	if scraper, exists := sm.scrapers[scraperType]; exists {
		return scraper, nil
	}
	return nil, fmt.Errorf("scraper type %v not found", scraperType)
}

// getScraperDisplayName returns a Portuguese display name for the scraper type
func (sm *ScraperManager) getScraperDisplayName(scraperType ScraperType) string {
	switch scraperType {
	case AnimesOnlineCCTYPE:
		return "AnimesOnlineCC"
	case AnimefireType:
		return "Animefire.io"
	case BetterAnimeType:
		return "BetterAnime.io"
	case TopAnimesType:
		return "TopAnimes.net"
	case AnimesDigitalType:
		return "AnimesDigital.org"
	case GoyabuType:
		return "Goyabu"
	case SuperAnimesType:
		return "SuperAnimes"
	case CinebyType:
		return "Cineby"
	case CineGratisType:
		return "CineGratis"
	case FlixHQType:
		return "FlixHQ"
	default:
		return "Desconhecido"
	}
}

// getLanguageTag returns a language tag for the source
func (sm *ScraperManager) getLanguageTag(scraperType ScraperType) string {
	switch scraperType {
	case AnimesOnlineCCTYPE, AnimefireType, GoyabuType, SuperAnimesType, BetterAnimeType, TopAnimesType, AnimesDigitalType:
		return "[Portuguese]"
	case CinebyType, CineGratisType, FlixHQType:
		return "[Movies/TV]"
	default:
		return "[Unknown]"
	}
}

// AnimefireAdapter adapts AnimefireClient to UnifiedScraper interface
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
	metadata := make(map[string]string)
	metadata["source"] = "animefire"
	return url, metadata, err
}

func (a *AnimefireAdapter) GetType() ScraperType {
	return AnimefireType
}

// GoyabuAdapter adapts GoyabuClient to UnifiedScraper interface
type GoyabuAdapter struct {
	client *GoyabuClient
}

func (a *GoyabuAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *GoyabuAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *GoyabuAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *GoyabuAdapter) GetType() ScraperType {
	return GoyabuType
}

// SuperAnimesAdapter adapts SuperAnimesClient to UnifiedScraper interface
type SuperAnimesAdapter struct {
	client *SuperAnimesClient
}

func (a *SuperAnimesAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *SuperAnimesAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *SuperAnimesAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *SuperAnimesAdapter) GetType() ScraperType {
	return SuperAnimesType
}

// AnimesOnlineCCAdapter adapts AnimesOnlineCCClient to UnifiedScraper interface
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

// CinebyAdapter adapts CinebyClient to UnifiedScraper interface
type CinebyAdapter struct {
	client *CinebyClient
}

func (a *CinebyAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchMedia(query)
}

func (a *CinebyAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *CinebyAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	streams, err := a.client.GetStreamURLs(episodeURL)
	if err != nil {
		return "", nil, err
	}
	if len(streams) == 0 {
		return "", nil, fmt.Errorf("no streams found")
	}
	metadata := map[string]string{
		"source":  "cineby",
		"quality": "default",
	}
	return streams[0], metadata, nil
}

func (a *CinebyAdapter) GetType() ScraperType {
	return CinebyType
}

// FlixHQAdapter adapts FlixHQClient to UnifiedScraper interface for movies and TV shows
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
		// Set the media type
		if m.Type == MediaTypeMovie {
			anime.MediaType = models.MediaTypeMovie
		} else {
			anime.MediaType = models.MediaTypeTV
		}
		anime.Year = m.Year
		animes = append(animes, anime)
	}

	return animes, nil
}

func (a *FlixHQAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	// For FlixHQ, animeURL contains the media ID
	// This needs to be called differently for movies vs TV shows
	// For movies, return a single "episode"
	// For TV shows, we need to get seasons first

	// This is a simplified implementation - in practice, you'd need to know if it's a movie or TV show
	return nil, fmt.Errorf("for FlixHQ, use GetSeasons and GetEpisodes directly on the client")
}

func (a *FlixHQAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	// Parse options
	provider := "Vidcloud"
	quality := "1080"
	subsLanguage := "english"

	for i, opt := range options {
		if s, ok := opt.(string); ok {
			switch i {
			case 0:
				provider = s
			case 1:
				quality = s
			case 2:
				subsLanguage = s
			}
		}
	}

	// Get embed link directly from episode ID
	embedLink, err := a.client.GetEmbedLink(episodeURL)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get embed link: %w", err)
	}

	streamInfo, err := a.client.ExtractStreamInfo(embedLink, quality, subsLanguage)
	if err != nil {
		return "", nil, fmt.Errorf("failed to extract stream info: %w", err)
	}

	metadata := make(map[string]string)
	metadata["source"] = "flixhq"
	metadata["provider"] = provider
	metadata["quality"] = quality

	// Include subtitle URLs in metadata
	if len(streamInfo.Subtitles) > 0 {
		var subURLs []string
		for _, sub := range streamInfo.Subtitles {
			subURLs = append(subURLs, sub.URL)
		}
		metadata["subtitles"] = strings.Join(subURLs, ",")
	}

	return streamInfo.VideoURL, metadata, nil
}

func (a *FlixHQAdapter) GetType() ScraperType {
	return FlixHQType
}

// GetClient returns the underlying FlixHQ client for direct access
func (a *FlixHQAdapter) GetClient() *FlixHQClient {
	return a.client
}
// CineGratisAdapter adapts CineGratisClient to UnifiedScraper interface
type CineGratisAdapter struct {
	client *CineGratisClient
}

func (a *CineGratisAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.Search(query)
}

func (a *CineGratisAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	if strings.Contains(animeURL, "/series/") {
		return a.client.GetEpisodes(animeURL)
	}
	// For movies, return the page itself as a single episode
	return []models.Episode{{Number: "Filme", Num: 1, URL: animeURL}}, nil
}

func (a *CineGratisAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	url, err := a.client.GetStreamURL(episodeURL)
	metadata := map[string]string{"source": "cinegratis"}
	return url, metadata, err
}

func (a *CineGratisAdapter) GetType() ScraperType {
	return CineGratisType
}

// BetterAnimeAdapter adapts BetterAnimeClient to UnifiedScraper interface
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

// TopAnimesAdapter adapts TopAnimesClient to UnifiedScraper interface
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

// AnimesDigitalAdapter adapts AnimesDigitalClient to UnifiedScraper interface
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
