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

	// Load Dynamic Scrapers (JSON Manifest from Spider)
	manifestURL := "http://localhost:3000/scrapers.json"
	dynamicConfigs, err := LoadDynamicScrapers(manifestURL)
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
		util.Error("Failed to load dynamic scrapers", "error", err)
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
	return nil, fmt.Errorf("use GetSeasons on client")
}

func (a *FlixHQAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	embedLink, _ := a.client.GetEmbedLink(episodeURL)
	info, err := a.client.ExtractStreamInfo(embedLink, "1080", "english")
	return info.VideoURL, map[string]string{"source": "flixhq"}, err
}

func (a *FlixHQAdapter) GetType() ScraperType {
	return FlixHQType
}

// GetClient returns the underlying FlixHQ client
func (a *FlixHQAdapter) GetClient() *FlixHQClient {
	return a.client
}
