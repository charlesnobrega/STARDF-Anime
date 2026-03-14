package scraper

import (
	"fmt"

	"github.com/alvarorichard/Goanime/internal/models"
)

// CinebyClient placeholder for movie source
type CinebyClient struct{}

func NewCinebyClient() *CinebyClient { return &CinebyClient{} }

func (c *CinebyClient) SearchMovies(query string) ([]*models.Anime, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CinebyClient) GetStreamURLs(movieURL string) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

// Adapter
type CinebyAdapter struct {
	client *CinebyClient
}

func NewCinebyAdapter(client *CinebyClient) *CinebyAdapter {
	return &CinebyAdapter{client: client}
}

func (a *CinebyAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchMovies(query)
}

func (a *CinebyAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return nil, nil
}

func (a *CinebyAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return "", nil, fmt.Errorf("not implemented")
}

func (a *CinebyAdapter) GetType() ScraperType {
	return CinebyType
}