package scraper

import (
	"fmt"

	"github.com/alvarorichard/Goanime/internal/models"
)

// GoyabuClient placeholder
type GoyabuClient struct{}

func NewGoyabuClient() *GoyabuClient { return &GoyabuClient{} }

func (c *GoyabuClient) SearchAnime(query string) ([]*models.Anime, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *GoyabuClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *GoyabuClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	return "", nil, fmt.Errorf("not implemented")
}

// Adapter
type GoyabuAdapter struct {
	client *GoyabuClient
}

func NewGoyabuAdapter(client *GoyabuClient) *GoyabuAdapter {
	return &GoyabuAdapter{client: client}
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