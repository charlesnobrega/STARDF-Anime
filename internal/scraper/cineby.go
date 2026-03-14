// Package scraper provides web scraping functionality for Cineby movies
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
	CinebyBase      = "https://www.cineby.gd"
	CinebyAPI       = "https://www.cineby.gd/api" // ajustar se necessário
	CinebyUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0"
)

// CinebyClient handles interactions with Cineby
type CinebyClient struct {
	client     *http.Client
	baseURL    string
	apiURL     string
	userAgent  string
	maxRetries int
	retryDelay time.Duration
}

// CinebyMovie represents a movie from Cineby
type CinebyMovie struct {
	ID       string
	Title    string
	Year     string
	ImageURL string
	URL      string
	Quality  string
	Duration string
	Genre    []string
	Rating   string
}

// CinebyStreamInfo contains streaming information
type CinebyStreamInfo struct {
	VideoURL   string
	Quality    string
	Subtitles  []CinebySubtitle
	Referer    string
	SourceName string
}

// CinebySubtitle represents a subtitle track
type CinebySubtitle struct {
	URL      string
	Language string
	Label    string
}

// NewCinebyClient creates a new Cineby client
func NewCinebyClient() *CinebyClient {
	return &CinebyClient{
		client:     util.GetFastClient(),
		baseURL:    CinebyBase,
		apiURL:     CinebyAPI,
		userAgent:  CinebyUserAgent,
		maxRetries: 2,
		retryDelay: 300 * time.Millisecond,
	}
}

// SearchMovies searches for movies on Cineby
func (c *CinebyClient) SearchMovies(query string) ([]*CinebyMovie, error) {
	// Implementar busca na página /search ou API
	// Exemplo: GET https://www.cineby.gd/search?q=query
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
		return nil, fmt.Errorf("cineby search failed: %s", resp.Status)
	}

	// Parse HTML results
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var movies []*CinebyMovie
	// Selecionar elementos de movie (ajustar selector conforme site)
	doc.Find(".movie-item, .film-card, .poster").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2, .title, .name").Text())
		href, _ := s.Find("a").First().Attr("href")
		img, _ := s.Find("img").First().Attr("src")
		year := strings.TrimSpace(s.Find(".year, .release-date").Text())

		if title != "" && href != "" {
			movie := &CinebyMovie{
				Title:    title,
				URL:     c.baseURL + href,
				ImageURL: img,
				Year:    year,
				ID:      generateID(title, year),
			}
			movies = append(movies, movie)
		}
	})

	return movies, nil
}

// GetMovieDetails fetches detailed movie information
func (c *CinebyClient) GetMovieDetails(movieURL string) (*CinebyMovie, error) {
	req, err := http.NewRequest("GET", movieURL, nil)
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
		return nil, fmt.Errorf("failed to fetch movie details: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	movie := &CinebyMovie{URL: movieURL}
	movie.Title = strings.TrimSpace(doc.Find("h1.title, .movie-title").First().Text())
	movie.Duration = strings.TrimSpace(doc.Find(".duration, .runtime").First().Text())
	movie.Rating = strings.TrimSpace(doc.Find(".rating, .imdb-rating").First().Text())
	movie.Quality = strings.TrimSpace(doc.Find(".quality, .badge").First().Text())

	// Genre
	doc.Find(".genre, .categories").Each(func(i int, s *goquery.Selection) {
		movie.Genre = append(movie.Genre, strings.TrimSpace(s.Text()))
	})

	// Image
	if img, exists := doc.Find("img.poster, .poster img").First().Attr("src"); exists {
		movie.ImageURL = img
	}

	movie.ID = generateID(movie.Title, movie.Year)
	return movie, nil
}

// GetStreamURLs returns available streaming URLs for a movie
func (c *CinebyClient) GetStreamURLs(movieURL string) ([]CinebyStreamInfo, error) {
	// Acessar página do filme e extrair links de streaming
	// Pode ser que estejam em iframes, players, ou botões "Assistir"
	req, err := http.NewRequest("GET", movieURL, nil)
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

	var streams []CinebyStreamInfo
	// Procurar por links de stream (iframe src, data-src, etc.)
	doc.Find("iframe, .player, a.watch-button, .stream-link").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.HasPrefix(src, "http") {
			streams = append(streams, CinebyStreamInfo{
				VideoURL:   src,
				Quality:    "default",
				SourceName: "Cineby",
			})
		}
		if dataSrc, ok := s.Attr("data-src"); ok && strings.HasPrefix(dataSrc, "http") {
			streams = append(streams, CinebyStreamInfo{
				VideoURL:   dataSrc,
				Quality:    "default",
				SourceName: "Cineby",
			})
		}
	})

	return streams, nil
}

// generateID creates a unique ID for a movie
func generateID(title, year string) string {
	cleanTitle := regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(strings.ToLower(title), "")
	return fmt.Sprintf("%s_%s", cleanTitle, year)
}