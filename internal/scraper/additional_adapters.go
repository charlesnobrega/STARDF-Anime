package scraper

import (
	"fmt"
	"strings"

	"github.com/alvarorichard/Goanime/internal/models"
)

// CinebyAdapter adapts Cineby client to UnifiedScraper interface
type CinebyAdapter struct {
	client *CinebyClient
}

// NewCinebyAdapter creates a new Cineby adapter
func NewCinebyAdapter(client *CinebyClient) *CinebyAdapter {
	return &CinebyAdapter{client: client}
}

// SearchAnime searches for movies on Cineby
func (a *CinebyAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	movies, err := a.client.SearchMovies(query)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
	for _, movie := range movies {
		anime := &models.Anime{
			ID:    movie.ID,
			Title: movie.Title,
			Type:  models.TypeMovie,
			ImageURL: movie.ImageURL,
			Year:  movie.Year,
			URL:   movie.URL,
			Source: "Cineby",
		}
		results = append(results, anime)
	}

	return results, nil
}

// GetAnimeEpisodes returns empty for Cineby (films don't have episodes)
func (a *CinebyAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return nil, nil
}

// GetStreamURL returns streaming URLs for a movie
func (a *CinebyAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	streams, err := a.client.GetStreamURLs(episodeURL)
	if err != nil {
		return "", nil, err
	}

	if len(streams) == 0 {
		return "", nil, fmt.Errorf("no streams found")
	}

	stream := streams[0]
	metadata := map[string]string{
		"source":    "cineby",
		"quality":   stream.Quality,
		"subtitles": "",
	}

	return stream.VideoURL, metadata, nil
}

func (a *CinebyAdapter) GetType() ScraperType {
	return CinebyType
}

// AnimesOnlineCCAdapter adapts AnimesOnlineCC client
type AnimesOnlineCCAdapter struct {
	client *AnimesOnlineCCClient
}

func NewAnimesOnlineCCAdapter(client *AnimesOnlineCCClient) *AnimesOnlineCCAdapter {
	return &AnimesOnlineCCAdapter{client: client}
}

func (a *AnimesOnlineCCAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	results, err := a.client.SearchAnime(query)
	if err != nil {
		return nil, err
	}
	for _, anime := range results {
		anime.Source = "AnimesOnlineCC"
	}
	return results, nil
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

// GoyabuAdapter adapts Goyabu client
type GoyabuAdapter struct {
	client *GoyabuClient
}

func NewGoyabuAdapter(client *GoyabuClient) *GoyabuAdapter {
	return &GoyabuAdapter{client: client}
}

func (a *GoyabuAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	animes, err := a.client.SearchAnime(query)
	if err != nil {
		return nil, err
	}
	var results []*models.Anime
	for _, anime := range animes {
		results = append(results, &models.Anime{
			ID:       anime.ID,
			Title:    anime.Title,
			URL:      anime.URL,
			ImageURL: anime.ImageURL,
			Year:     anime.Year,
			Source:   "Goyabu",
		})
	}
	return results, nil
}

func (a *GoyabuAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *GoyabuAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	info, err := a.client.GetStreamURL(episodeURL)
	if err != nil {
		return "", nil, err
	}
	metadata := map[string]string{
		"source":  "goyabu",
		"quality": info.Quality,
	}
	return info.VideoURL, metadata, nil
}

func (a *GoyabuAdapter) GetType() ScraperType {
	return GoyabuType
}

// SuperAnimesAdapter adapts SuperAnimes client
type SuperAnimesAdapter struct {
	client *SuperAnimesClient
}

func NewSuperAnimesAdapter(client *SuperAnimesClient) *SuperAnimesAdapter {
	return &SuperAnimesAdapter{client: client}
}

func (a *SuperAnimesAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	results, err := a.client.SearchAnime(query)
	if err != nil {
		return nil, err
	}
	return results, nil
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