package scraper

import (
	"fmt"

	"github.com/alvarorichard/Goanime/internal/models"
)

// SuperAnimesClient placeholder
type SuperAnimesClient struct{}

func NewSuperAnimesClient() *SuperAnimesClient { return &SuperAnimesClient{} }

func (c *SuperAnimesClient) SearchAnime(query string) ([]*models.Anime, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *SuperAnimesClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *SuperAnimesClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	return "", nil, fmt.Errorf("not implemented")
}

// Adapter
type SuperAnimesAdapter struct {
	client *SuperAnimesClient
}

func NewSuperAnimesAdapter(client *SuperAnimesClient) *SuperAnimesAdapter {
	return &SuperAnimesAdapter{client: client}
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