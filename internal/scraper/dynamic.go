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
	resp, err := s.client.Get(animeURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	selector := s.Config.Selectors["episodes"]
	if selector == "" {
		selector = ".episodios li a, .episodes a, ul.episodes li a" // Fallback
	}

	var episodes []models.Episode
	doc.Find(selector).Each(func(i int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		if link == "" {
			return
		}
		if !strings.HasPrefix(link, "http") {
			link = s.Config.BaseURL + link
		}

		title := strings.TrimSpace(sel.Text())
		if title == "" {
			title = fmt.Sprintf("Episódio %d", i+1)
		}

		episodes = append(episodes, models.Episode{
			Number: fmt.Sprintf("%d", i+1),
			Num:    i + 1,
			Title:  models.TitleDetails{English: title},
			URL:    link,
		})
	})

	return episodes, nil
}

func (s *DynamicScraper) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	resp, err := s.client.Get(episodeURL)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var streamURL string
	selector := s.Config.Selectors["stream"]
	if selector == "" {
		selector = "iframe, video source" // Broad Fallback
	}

	// Try multiple common attributes for the URL
	doc.Find(selector).Each(func(i int, sel *goquery.Selection) {
		if streamURL != "" { return }
		
		attributes := []string{"src", "data-src", "data-video", "data-l"}
		for _, attr := range attributes {
			val, _ := sel.Attr(attr)
			if val != "" && isValidStream(val) {
				streamURL = val
				break
			}
		}
	})

	if streamURL == "" {
		// Final attempt: check for matches in script tags (dangerous but sometimes needed)
		doc.Find("script").Each(func(i int, sel *goquery.Selection) {
			if streamURL != "" { return }
			content := sel.Text()
			if strings.Contains(content, "var player") || strings.Contains(content, "const v =") {
				// regex search for URL patterns starting with http could be done here if needed
			}
		})
	}

	if streamURL == "" {
		return "", nil, fmt.Errorf("could not find stream URL with selector: %s", selector)
	}

	// Handle relative URLs
	if !strings.HasPrefix(streamURL, "http") {
		if strings.HasPrefix(streamURL, "//") {
			streamURL = "https:" + streamURL
		} else {
			streamURL = strings.TrimSuffix(s.Config.BaseURL, "/") + "/" + strings.TrimPrefix(streamURL, "/")
		}
	}

	metadata := map[string]string{
		"source": s.Config.Name,
	}

	return streamURL, metadata, nil
}

// isValidStream checks if a URL is likely a valid stream and not an ad
func isValidStream(u string) bool {
	u = strings.ToLower(u)
	blocked := []string{"google", "ads", "advertising", "banner", "analytics", "doubleclick"}
	for _, b := range blocked {
		if strings.Contains(u, b) {
			return false
		}
	}
	return strings.HasPrefix(u, "http") || strings.HasPrefix(u, "//") || strings.HasPrefix(u, "/")
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
