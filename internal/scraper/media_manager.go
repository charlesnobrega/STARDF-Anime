// Package scraper provides unified media handling for anime, movies, and TV shows
package scraper

import (
	"fmt"
	"strings"

	"github.com/charlesnobrega/STARDF-Anime/internal/models"
)

// MediaManager provides a unified interface for all media types
type MediaManager struct {
	scraperManager *ScraperManager
	flixhqClient   *FlixHQClient
}

// NewMediaManager creates a new MediaManager
func NewMediaManager() *MediaManager {
	sm := NewScraperManager()

	// Get the FlixHQ client from the adapter
	var flixhqClient *FlixHQClient
	if s, exists := sm.scrapers[FlixHQType]; exists {
		if adapter, ok := s.(*FlixHQAdapter); ok {
			flixhqClient = adapter.client
		}
	}
	if flixhqClient == nil {
		flixhqClient = NewFlixHQClient()
	}

	return &MediaManager{
		scraperManager: sm,
		flixhqClient:   flixhqClient,
	}
}

// SearchAll searches across all sources (anime + movies/TV)
func (mm *MediaManager) SearchAll(query string) ([]*models.Anime, error) {
	return mm.scraperManager.SearchAnime(query, nil)
}

// SearchAnimeOnly searches only anime sources (dynamically)
func (mm *MediaManager) SearchAnimeOnly(query string) ([]*models.Anime, error) {
	results, err := mm.scraperManager.SearchAnime(query, nil)
	if err != nil {
		return nil, err
	}

	var filtered []*models.Anime
	for _, a := range results {
		if !strings.Contains(a.Name, "[Movies/TV]") {
			filtered = append(filtered, a)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no results found for: %s", query)
	}

	return filtered, nil
}

// SearchMoviesAndTV searches only FlixHQ for movies and TV shows
func (mm *MediaManager) SearchMoviesAndTV(query string) ([]*FlixHQMedia, error) {
	return mm.flixhqClient.SearchMedia(query)
}

// GetTrendingMovies gets trending movies from FlixHQ
func (mm *MediaManager) GetTrendingMovies() ([]*FlixHQMedia, error) {
	return mm.flixhqClient.GetTrending()
}

// GetRecentMovies gets recent movies from FlixHQ
func (mm *MediaManager) GetRecentMovies() ([]*FlixHQMedia, error) {
	return mm.flixhqClient.GetRecentMovies()
}

// GetRecentTV gets recent TV shows from FlixHQ
func (mm *MediaManager) GetRecentTV() ([]*FlixHQMedia, error) {
	return mm.flixhqClient.GetRecentTV()
}

// GetTVSeasons gets all seasons for a TV show
func (mm *MediaManager) GetTVSeasons(mediaID string) ([]FlixHQSeason, error) {
	return mm.flixhqClient.GetSeasons(mediaID)
}

// GetTVEpisodes gets all episodes for a season
func (mm *MediaManager) GetTVEpisodes(seasonID string) ([]FlixHQEpisode, error) {
	return mm.flixhqClient.GetEpisodes(seasonID)
}

// GetMovieStreamInfo gets stream information for a movie
func (mm *MediaManager) GetMovieStreamInfo(mediaID, provider, quality, subsLanguage string) (*FlixHQStreamInfo, error) {
	if provider == "" { provider = "Vidcloud" }
	if quality == "" { quality = "1080" }
	if subsLanguage == "" { subsLanguage = "english" }

	episodeID, err := mm.flixhqClient.GetMovieServerID(mediaID, provider)
	if err != nil { return nil, err }

	embedLink, err := mm.flixhqClient.GetEmbedLink(episodeID)
	if err != nil { return nil, err }

	return mm.flixhqClient.ExtractStreamInfo(embedLink, quality, subsLanguage)
}

// GetTVEpisodeStreamInfo gets stream information for a TV episode
func (mm *MediaManager) GetTVEpisodeStreamInfo(dataID, provider, quality, subsLanguage string) (*FlixHQStreamInfo, error) {
	if provider == "" { provider = "Vidcloud" }
	if quality == "" { quality = "1080" }
	if subsLanguage == "" { subsLanguage = "english" }

	episodeID, err := mm.flixhqClient.GetEpisodeServerID(dataID, provider)
	if err != nil { return nil, err }

	embedLink, err := mm.flixhqClient.GetEmbedLink(episodeID)
	if err != nil { return nil, err }

	return mm.flixhqClient.ExtractStreamInfo(embedLink, quality, subsLanguage)
}

// GetAnimeStreamURL gets stream URL for anime episodes dynamically
func (mm *MediaManager) GetAnimeStreamURL(anime *models.Anime, episodeURL string, options ...interface{}) (string, map[string]string, error) {
	scraper, err := mm.scraperManager.FindScraperByName(anime.Source)
	if err != nil {
		return "", nil, err
	}
	return scraper.GetStreamURL(anime.URL, options...)
}

// ConvertFlixHQToAnime converts FlixHQ media list to Anime models
func ConvertFlixHQToAnime(media []*FlixHQMedia) []*models.Anime {
	var animes []*models.Anime
	for _, m := range media {
		anime := m.ToAnimeModel()
		if m.Type == MediaTypeMovie {
			anime.MediaType = models.MediaTypeMovie
		} else {
			anime.MediaType = models.MediaTypeTV
		}
		anime.Year = m.Year
		animes = append(animes, anime)
	}
	return animes
}

// ConvertFlixHQEpisodesToEpisodes converts FlixHQ episodes to Episode models
func ConvertFlixHQEpisodesToEpisodes(episodes []FlixHQEpisode) []models.Episode {
	var eps []models.Episode
	for _, e := range episodes {
		eps = append(eps, e.ToEpisodeModel())
	}
	return eps
}

// GetScraperManager returns the underlying scraper manager
func (mm *MediaManager) GetScraperManager() *ScraperManager {
	return mm.scraperManager
}

// GetFlixHQClient returns the FlixHQ client
func (mm *MediaManager) GetFlixHQClient() *FlixHQClient {
	return mm.flixhqClient
}
