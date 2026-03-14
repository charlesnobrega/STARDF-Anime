package scraper

import (
	"fmt"

	"github.com/alvarorichard/Goanime/internal/models"
)

// AnimesOnlineCCClient placeholder
type AnimesOnlineCCClient struct{}

func NewAnimesOnlineCCClient() *AnimesOnlineCCClient { return &AnimesOnlineCCClient{} }

func (c *AnimesOnlineCCClient) SearchAnime(query string) ([]*models.Anime, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *AnimesOnlineCCClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *AnimesOnlineCCClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	return "", nil, fmt.Errorf("not implemented")
}

// Adapter
type AnimesOnlineCCAdapter struct {
	client *AnimesOnlineCCClient
}

func NewAnimesOnlineCCAdapter(client *AnimesOnlineCCClient) *AnimesOnlineCCAdapter {
	return &AnimesOnlineCCAdapter{client: client}
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