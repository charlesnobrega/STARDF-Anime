package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

type DynamicScraperConfig struct {
	Name           string            `json:"name"`
	BaseURL        string            `json:"baseURL"`
	SearchURL      string            `json:"searchURL"`
	ParentSelector string            `json:"parentSelector"`
	Selectors      map[string]string `json:"selectors"`
}

type DynamicScraper struct {
	client      *http.Client
	Config      DynamicScraperConfig
	scraperType ScraperType
}

func NewDynamicScraper(config DynamicScraperConfig, st ScraperType) *DynamicScraper {
	return &DynamicScraper{
		client:      util.GetScraperClient(),
		Config:      config,
		scraperType: st,
	}
}

func (s *DynamicScraper) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	searchURL := fmt.Sprintf(s.Config.SearchURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
	doc.Find(s.Config.ParentSelector).Each(func(i int, sel *goquery.Selection) {
		name := strings.TrimSpace(sel.Find(s.Config.Selectors["title"]).Text())
		link, _ := sel.Find(s.Config.Selectors["link"]).Attr("href")
		image, _ := sel.Find(s.Config.Selectors["image"]).Attr("src")

		if name == "" || link == "" {
			// Fallback: try different variant for link
			if link == "" {
				link, _ = sel.Attr("href")
			}
			if name == "" && link != "" {
				name = link
			}
		}

		if name != "" && link != "" {
			if !strings.HasPrefix(link, "http") {
				link = s.Config.BaseURL + link
			}
			if image != "" && !strings.HasPrefix(image, "http") {
				image = s.Config.BaseURL + image
			}

			results = append(results, &models.Anime{
				Name:     name,
				URL:      link,
				ImageURL: image,
				Source:   s.Config.Name,
			})
		}
	})

	return results, nil
}

func (s *DynamicScraper) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	// For now, let's keep it simple. In a real scenario, the JSON should also 
	// contain selectors for episodes.
	return nil, fmt.Errorf("episodes not yet implemented in dynamic scraper")
}

func (s *DynamicScraper) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return "", nil, fmt.Errorf("stream url not yet implemented in dynamic scraper")
}

func (s *DynamicScraper) GetType() ScraperType {
	return s.scraperType
}

type DynamicManifest struct {
	Scrapers []DynamicScraperConfig `json:"scrapers"`
}

func LoadDynamicScrapers(manifestURL string) ([]DynamicScraperConfig, error) {
	resp, err := http.Get(manifestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var manifest DynamicManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	return manifest.Scrapers, nil
}
