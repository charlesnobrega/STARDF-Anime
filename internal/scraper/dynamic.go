package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
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
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", s.Config.BaseURL)

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
		selector = ".episodios li a, .episodes a, ul.episodes li a, .lEp, .episodiotitle a" // Robust Fallback
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

	// Enhanced extraction logic for modern players
	doc.Find(selector).Each(func(i int, sel *goquery.Selection) {
		if streamURL != "" { return }
		
		// 1. Check common attributes (src, data-src, link, etc)
		for _, attr := range []string{"src", "data-src", "data-video", "value", "data-l"} {
			if val, ok := sel.Attr(attr); ok && isValidStream(val) {
				streamURL = val
				return
			}
		}

		// 2. Check for nested video/source elements if selector was a container
		sel.Find("video, source, iframe").Each(func(j int, sub *goquery.Selection) {
			if streamURL != "" { return }
			for _, attr := range []string{"src", "data-src"} {
				if val, ok := sub.Attr(attr); ok && isValidStream(val) {
					streamURL = val
					return
				}
			}
		})
	})

	if streamURL == "" {
		// 3. Fallback: Search the entire HTML for "link=" or "file=" patterns in strings/scripts
		html, _ := doc.Html()
		// Regex for searching video URLs in various script/attr contexts
		re := regexp.MustCompile(`(?:link|file|src|url)\s*[:=]\s*["'](https?://[^"']+)["']`)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 && isValidStream(matches[1]) {
			streamURL = matches[1]
		}
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
