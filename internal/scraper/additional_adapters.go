func (a *CinebyAdapter) GetType() ScraperType {
	return CinebyType
}

// AnimesOnlineCCAdapter adapts AnimesOnlineCC client to UnifiedScraper interface
type AnimesOnlineCCAdapter struct {
	client *AnimesOnlineCCClient
}

// NewAnimesOnlineCCAdapter creates a new AnimesOnlineCC adapter
func NewAnimesOnlineCCAdapter(client *AnimesOnlineCCClient) *AnimesOnlineCCAdapter {
	return &AnimesOnlineCCAdapter{client: client}
}

// SearchAnime searches for anime on AnimesOnlineCC
func (a *AnimesOnlineCCAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	results, err := a.client.SearchAnime(query)
	if err != nil {
		return nil, err
	}
	// Tag as Brazilian source
	for _, anime := range results {
		anime.Source = "AnimesOnlineCC"
	}
	return results, nil
}

// GetAnimeEpisodes returns episodes for an anime
func (a *AnimesOnlineCCAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

// GetStreamURL returns streaming URL for an episode
func (a *AnimesOnlineCCAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *AnimesOnlineCCAdapter) GetType() ScraperType {
	return AnimesOnlineCCTYPE
}

// GoyabuAdapter adapts Goyabu client to UnifiedScraper interface
type GoyabuAdapter struct {
	client *GoyabuClient
}

// NewGoyabuAdapter creates a new Goyabu adapter
func NewGoyabuAdapter(client *GoyabuClient) *GoyabuAdapter {
	return &GoyabuAdapter{client: client}
}

// SearchAnime searches for anime on Goyabu
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

// GetAnimeEpisodes returns episodes for an anime
func (a *GoyabuAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

// GetStreamURL returns streaming URL for an episode
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

// SuperAnimesAdapter adapts SuperAnimes client to UnifiedScraper interface
type SuperAnimesAdapter struct {
	client *SuperAnimesClient
}

// NewSuperAnimesAdapter creates a new SuperAnimes adapter
func NewSuperAnimesAdapter(client *SuperAnimesClient) *SuperAnimesAdapter {
	return &SuperAnimesAdapter{client: client}
}

// SearchAnime searches for anime on SuperAnimes
func (a *SuperAnimesAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	results, err := a.client.SearchAnime(query)
	if err != nil {
		return nil, err
	}
	// Source already set in client
	return results, nil
}

// GetAnimeEpisodes returns episodes for an anime
func (a *SuperAnimesAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

// GetStreamURL returns streaming URL for an episode
func (a *SuperAnimesAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *SuperAnimesAdapter) GetType() ScraperType {
	return SuperAnimesType
}