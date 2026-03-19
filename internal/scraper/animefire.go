// Package scraper provides web scraping functionality for animefire.io
package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

const (
	AnimefireBase  = "https://animefire.io"
	AnimefireAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// AnimefireClient handles interactions with Animefire.io
type AnimefireClient struct {
	client     *http.Client
	baseURL    string
	userAgent  string
	maxRetries int
	retryDelay time.Duration
}

// NewAnimefireClient creates a new Animefire client
func NewAnimefireClient() *AnimefireClient {
	return &AnimefireClient{
		client:     util.GetFastClient(), // Use shared fast client
		baseURL:    AnimefireBase,
		userAgent:  AnimefireAgent,
		maxRetries: 2,
		retryDelay: 250 * time.Millisecond, // Reduced from 350ms
	}
}

// SearchAnime searches for anime on Animefire.io using the original logic
func (c *AnimefireClient) SearchAnime(query string) ([]*models.Anime, error) {
	// AnimeFire expects spaces as hyphens in the URL
	normalizedQuery := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(query)), " ", "-")
	searchURL := fmt.Sprintf("%s/pesquisar/%s", c.baseURL, normalizedQuery)

	util.Debug("AnimeFire search", "query", query, "normalized", normalizedQuery, "url", searchURL)

	var lastErr error
	attempts := c.maxRetries + 1

	for attempt := 0; attempt < attempts; attempt++ {
		req, err := http.NewRequest("GET", searchURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		c.decorateRequest(req)

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to make request: %w", err)
			if c.shouldRetry(attempt) {
				c.sleep()
				continue
			}
			return nil, lastErr
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = c.handleStatusError(resp)
			_ = resp.Body.Close()
			if c.shouldRetry(attempt) {
				c.sleep()
				continue
			}
			return nil, lastErr
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to parse HTML: %w", err)
			if c.shouldRetry(attempt) {
				c.sleep()
				continue
			}
			return nil, lastErr
		}

		if c.isChallengePage(doc) {
			lastErr = errors.New("animefire returned a challenge page (try VPN or wait)")
			if c.shouldRetry(attempt) {
				c.sleep()
				continue
			}
			return nil, lastErr
		}

		animes := c.extractSearchResults(doc)
		if len(animes) == 0 {
			// Legitimate empty result set – return without error
			return []*models.Anime{}, nil
		}

		return animes, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("failed to retrieve results from AnimeFire")
}

func (c *AnimefireClient) decorateRequest(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", c.baseURL+"/")
}

func (c *AnimefireClient) handleStatusError(resp *http.Response) error {
	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("access restricted: VPN may be required")
	}
	return fmt.Errorf("server returned: %s", resp.Status)
}

func (c *AnimefireClient) shouldRetry(attempt int) bool {
	return attempt < c.maxRetries
}

func (c *AnimefireClient) sleep() {
	if c.retryDelay <= 0 {
		return
	}
	time.Sleep(c.retryDelay)
}

func (c *AnimefireClient) isChallengePage(doc *goquery.Document) bool {
	title := strings.ToLower(strings.TrimSpace(doc.Find("title").First().Text()))
	if strings.Contains(title, "just a moment") {
		return true
	}

	if doc.Find("#cf-wrapper").Length() > 0 || doc.Find("#challenge-form").Length() > 0 {
		return true
	}

	body := strings.ToLower(doc.Text())
	return strings.Contains(body, "cf-error") || strings.Contains(body, "cloudflare")
}

func (c *AnimefireClient) extractSearchResults(doc *goquery.Document) []*models.Anime {
	var animes []*models.Anime

	doc.Find(".row.ml-1.mr-1 a").Each(func(i int, s *goquery.Selection) {
		if urlPath, exists := s.Attr("href"); exists {
			name := strings.TrimSpace(s.Text())
			if name != "" {
				animes = append(animes, &models.Anime{
					Name: name,
					URL:  c.resolveURL(c.baseURL, urlPath),
				})
			}
		}
	})

	if len(animes) > 0 {
		return animes
	}

	doc.Find(".card_ani").Each(func(i int, s *goquery.Selection) {
		titleElem := s.Find(".ani_name a")
		title := strings.TrimSpace(titleElem.Text())
		link, exists := titleElem.Attr("href")

		if exists && title != "" {
			imgElem := s.Find(".div_img img")
			imgURL, _ := imgElem.Attr("src")
			if imgURL != "" {
				imgURL = c.resolveURL(c.baseURL, imgURL)
			}

			animes = append(animes, &models.Anime{
				Name:     title,
				URL:      c.resolveURL(c.baseURL, link),
				ImageURL: imgURL,
			})
		}
	})

	return animes
}

// resolveURL resolves relative URLs to absolute URLs
func (c *AnimefireClient) resolveURL(base, ref string) string {
	if strings.HasPrefix(ref, "http") {
		return ref
	}
	if strings.HasPrefix(ref, "/") {
		return base + ref
	}
	return base + "/" + ref
}

func (c *AnimefireClient) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	c.decorateRequest(req)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var episodes []models.Episode
	doc.Find(".div_episodes a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		num := i + 1
		title := strings.TrimSpace(s.Text())
		
		if href != "" {
			episodes = append(episodes, models.Episode{
				Number: fmt.Sprintf("%d", num),
				Num:    num,
				URL:    c.resolveURL(c.baseURL, href),
				Title:  models.TitleDetails{English: title},
			})
		}
	})
	
	return episodes, nil
}

// GetEpisodeStreamURL gets the streaming URL for a specific episode
// This is specific to scraper functionality, not duplicated from API
func (c *AnimefireClient) GetEpisodeStreamURL(episodeURL string) (string, error) {
	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return "", err
	}
	c.decorateRequest(req)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	
	// Stream URL is often in a script or hidden in an iframe
	var videoURL string
	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.Contains(src, "player") {
			videoURL = src
		}
	})
	
	if videoURL == "" {
		// Try searching for .m3u8 or .mp4 in scripts
		html, _ := doc.Html()
		re := regexp.MustCompile(`(https?://[^"']+\.(m3u8|mp4)[^"']*)`)
		if matches := re.FindStringSubmatch(html); len(matches) > 1 {
			videoURL = matches[1]
		}
	}
	
	if videoURL == "" {
		return "", errors.New("no stream found on page")
	}
	
	return videoURL, nil
}

// GetAnimeDetails - placeholder method, details are fetched by API layer
func (c *AnimefireClient) GetAnimeDetails(animeURL string) (*models.Anime, error) {
	return nil, fmt.Errorf("anime details should be fetched using API layer, not scraper")
}
