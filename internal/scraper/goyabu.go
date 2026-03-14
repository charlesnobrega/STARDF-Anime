// Package scraper provides web scraping functionality for Goyabu anime
package scraper

import (
	"encoding/json"
	"errors"
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
	GoyabuBase      = "https://goyabu.io"
	GoyabuAPI       = "https://api.goyabu.io" // verificar se há API
	GoyabuUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0"
)

// GoyabuClient handles interactions with Goyabu
type GoyabuClient struct {
	client     *http.Client
	baseURL    string
	apiURL     string
	userAgent  string
	maxRetries int
	retryDelay time.Duration
}

// GoyabuAnime represents an anime from Goyabu
type GoyabuAnime struct {
	ID       string
	Title    string
	URL      string
	ImageURL string
	Year     string
	Type     string
	Genres   []string
}

// GoyabuEpisode represents an episode
type GoyabuEpisode struct {
	ID    string
	Number int
	Title  string
	URL    string
}

// GoyabuStreamInfo contains streaming information
type GoyabuStreamInfo struct {
	VideoURL  string
	Quality   string
	Subtitles []string
	Referer   string
}

// NewGoyabuClient creates a new Goyabu client
func NewGoyabuClient() *GoyabuClient {
	return &GoyabuClient{
		client:     util.GetFastClient(),
		baseURL:    GoyabuBase,
		apiURL:     GoyabuAPI,
		userAgent:  GoyabuUserAgent,
		maxRetries: 2,
		retryDelay: 300 * time.Millisecond,
	}
}

// SearchAnime searches for anime on Goyabu
func (c *GoyabuClient) SearchAnime(query string) ([]*GoyabuAnime, error) {
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
		return nil, fmt.Errorf("goyabu search failed: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*GoyabuAnime
	doc.Find(".anime-item, .card, .poster").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3, .title").Text())
		href, _ := s.Find("a").First().Attr("href")
		img, _ := s.Find("img").First().Attr("src")
		year := strings.TrimSpace(s.Find(".year, .date").Text())

		if title != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			anime := &GoyabuAnime{
				Title:    title,
				URL:      href,
				ImageURL: img,
				Year:    year,
				ID:      generateID(title),
			}
			results = append(results, anime)
		}
	})

	return results, nil
}

// GetAnimeDetails fetches anime details
func (c *GoyabuClient) GetAnimeDetails(animeURL string) (*GoyabuAnime, error) {
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

	anime := &GoyabuAnime{URL: animeURL}
	anime.Title = strings.TrimSpace(doc.Find("h1.title").First().Text())

	// Genres
	doc.Find(".genres a, .category").Each(func(i int, s *goquery.Selection) {
		anime.Genres = append(anime.Genres, strings.TrimSpace(s.Text()))
	})

	// Image
	if img, ok := doc.Find("img.poster").First().Attr("src"); ok {
		anime.ImageURL = img
	}

	anime.ID = generateID(anime.Title)
	return anime, nil
}

// GetEpisodes returns episode list
func (c *GoyabuClient) GetEpisodes(animeURL string) ([]GoyabuEpisode, error) {
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

	var episodes []GoyabuEpisode
	doc.Find(".episodes-list a, .episode a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := strings.TrimSpace(s.Text())
		num := i + 1

		if href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			episodes = append(episodes, GoyabuEpisode{
				Number: num,
				Title:  title,
				URL:    href,
				ID:     fmt.Sprintf("%s_ep%d", generateID(""), num),
			})
		}
	})

	return episodes, nil
}

// GetStreamURL returns video URL for an episode
func (c *GoyabuClient) GetStreamURL(episodeURL string) (GoyabuStreamInfo, error) {
	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return GoyabuStreamInfo{}, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return GoyabuStreamInfo{}, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return GoyabuStreamInfo{}, err
	}

	var info GoyabuStreamInfo
	// Procurar video (iframe, video tag)
	doc.Find("iframe, video").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.HasPrefix(src, "http") {
			info.VideoURL = src
		}
		if src, ok := s.Attr("data-src"); ok && strings.HasPrefix(src, "http") {
			info.VideoURL = src
		}
	})

	info.Referer = episodeURL
	info.Quality = "default"
	return info, nil
}

// generateID creates a unique ID
func generateID(title string) string {
	clean := regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(strings.ToLower(title), "")
	return clean
}