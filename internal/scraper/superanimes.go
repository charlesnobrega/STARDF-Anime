// Package scraper provides web scraping functionality for SuperAnimes
package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alvarorichard/Goanime/internal/models"
	"github.com/alvarorichard/Goanime/internal/util"
)

const (
	SuperAnimesBase      = "https://superanimes.in"
	SuperAnimesUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0"
)

// SuperAnimesClient handles interactions with SuperAnimes
type SuperAnimesClient struct {
	client     *http.Client
	baseURL    string
	userAgent  string
	maxRetries int
	retryDelay time.Duration
}

// NewSuperAnimesClient creates a new SuperAnimes client
func NewSuperAnimesClient() *SuperAnimesClient {
	return &SuperAnimesClient{
		client:     util.GetFastClient(),
		baseURL:    SuperAnimesBase,
		userAgent:  SuperAnimesUserAgent,
		maxRetries: 2,
		retryDelay: 300 * time.Millisecond,
	}
}

// SearchAnime searches for anime on SuperAnimes
func (c *SuperAnimesClient) SearchAnime(query string) ([]*models.Anime, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("superanimes search failed: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
	doc.Find(".anime-card, .poster, .item").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3, .title, .name").Text())
		href, _ := s.Find("a").First().Attr("href")
		img, _ := s.Find("img").First().Attr("src")

		if title != "" && href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			anime := &models.Anime{
				ID:       generateID(title),
				Title:    title,
				URL:      href,
				ImageURL: img,
				Source:   "SuperAnimes",
			}
			results = append(results, anime)
		}
	})

	return results, nil
}

// GetAnimeDetails fetches detailed anime info (placeholder)
func (c *SuperAnimesClient) GetAnimeDetails(animeURL string) (*models.Anime, error) {
	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	anime := &models.Anime{URL: animeURL}
	anime.Title = strings.TrimSpace(doc.Find("h1.title").First().Text())
	anime.Description = strings.TrimSpace(doc.Find(".description, .synopsis").First().Text())

	return anime, nil
}

// GetEpisodes returns episode list
func (c *SuperAnimesClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

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
	doc.Find(".episodes-list a, .episode a, .list-episodes a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := strings.TrimSpace(s.Text())
		num := i + 1

		if href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			episodes = append(episodes, models.Episode{
				ID:    fmt.Sprintf("%s_ep%d", generateID(""), num),
				Number: num,
				Title:  title,
				URL:    href,
			})
		}
	})

	return episodes, nil
}

// GetStreamURL returns video URL for an episode
func (c *SuperAnimesClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var videoURL string
	doc.Find("iframe, video, .player").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.HasPrefix(src, "http") {
			videoURL = src
		}
		if dataSrc, ok := s.Attr("data-src"); ok && strings.HasPrefix(dataSrc, "http") {
			videoURL = dataSrc
		}
	})

	if videoURL == "" {
		return "", nil, errors.New("no video found")
	}

	metadata := map[string]string{
		"source":  "superanimes",
		"quality": "default",
	}
	return videoURL, metadata, nil
}

// generateID creates a unique ID
func generateID(title string) string {
	clean := regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(strings.ToLower(title), "")
	return clean
}