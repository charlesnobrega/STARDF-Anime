package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alvarorichard/Goanime/internal/models"
	"github.com/alvarorichard/Goanime/internal/util"
)

const (
	CinebyBase      = "https://www.cineby.gd"
	CinebyAPI       = "https://www.cineby.gd/api"
	CinebyUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type CinebyClient struct {
	client     *http.Client
	baseURL    string
	apiURL     string
	userAgent  string
	maxRetries int
	retryDelay time.Duration
}

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

// SearchMovies busca filmes no Cineby
func (c *CinebyClient) SearchMovies(query string) ([]*models.Anime, error) {
	// Tentar API primeiro
	endpoint := fmt.Sprintf("%s/search?q=%s", c.apiURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			if results, ok := data["results"].([]interface{}); ok {
				var movies []*models.Anime
				for _, r := range results {
					if m, ok := r.(map[string]interface{}); ok {
						movies = append(movies, &models.Anime{
							Name:      fmt.Sprintf("%v", m["title"]),
							URL:       fmt.Sprintf("%v", m["url"]),
							ImageURL:  fmt.Sprintf("%v", m["poster"]),
							MediaType: models.MediaTypeMovie,
							Source:    "Cineby",
							Year:      fmt.Sprintf("%v", m["year"]),
						})
					}
				}
				if len(movies) > 0 {
					return movies, nil
				}
			}
		}
	}

	// Fallback: scraping HTML
	return c.searchHTML(query)
}

func (c *CinebyClient) searchHTML(query string) ([]*models.Anime, error) {
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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var movies []*models.Anime
	// Seletores comuns para filmes
	selectors := []string{".movie-card", ".film-item", ".poster", ".movie-item", "[data-movie]"}
	for _, sel := range selectors {
		doc.Find(sel).Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Find("h3, .title, .name, a").Text())
			href, _ := s.Find("a").First().Attr("href")
			img, _ := s.Find("img").First().Attr("src")
			year := strings.TrimSpace(s.Find(".year, .release-date, .date").Text())

			if title == "" {
				return
			}
			if href == "" {
				if a, ok := s.Attr("href"); ok {
					href = a
				} else if a, ok := s.Attr("data-url"); ok {
					href = a
				}
			}
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			if img != "" && !strings.HasPrefix(img, "http") {
				img = c.baseURL + img
			}

			movies = append(movies, &models.Anime{
				Name:      title,
				URL:       href,
				ImageURL:  img,
				MediaType: models.MediaTypeMovie,
				Source:    "Cineby",
				Year:      year,
			})
		})
		if len(movies) > 0 {
			return movies, nil
		}
	}

	return movies, nil
}

// GetStreamURLs retorna URLs de streaming para um filme
func (c *CinebyClient) GetStreamURLs(movieURL string) ([]string, error) {
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

	var streams []string
	// Procurar por players comuns
	doc.Find("iframe, .player, video, .video-container, .streaming, a.watch, a.btn-play").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.HasPrefix(src, "http") {
			streams = append(streams, src)
		}
		if dataSrc, ok := s.Attr("data-src"); ok && strings.HasPrefix(dataSrc, "http") {
			streams = append(streams, dataSrc)
		}
		if src, ok := s.Attr("data-video"); ok && strings.HasPrefix(src, "http") {
			streams = append(streams, src)
		}
		// Json-LD embedding
		if script, ok := s.Attr("data-player"); ok {
			// Pode conter URL embutida
		}
	})

	// Se não encontrou, tentar extrair de scripts
	if len(streams) == 0 {
		doc.Find("script").Each(func(i int, s *goquery.Selection) {
			scriptText := s.Text()
			// Procurar URLs http... em scripts
			if strings.Contains(scriptText, "http") && (strings.Contains(scriptText, "mp4") || strings.Contains(scriptText, "m3u8")) {
				// Extrair URL (simplificado)
				for _, part := range strings.Fields(scriptText) {
					if strings.HasPrefix(part, "http") && (strings.HasSuffix(part, ".mp4") || strings.HasSuffix(part, ".m3u8")) {
						streams = append(streams, part)
					}
				}
			}
		})
	}

	return streams, nil
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
	// Filmes não têm episódios
	return nil, nil
}

func (a *CinebyAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	streams, err := a.client.GetStreamURLs(episodeURL)
	if err != nil {
		return "", nil, err
	}
	if len(streams) == 0 {
		return "", nil, fmt.Errorf("no streams found")
	}
	metadata := map[string]string{
		"source":  "cineby",
		"quality": "default",
	}
	return streams[0], metadata, nil
}

func (a *CinebyAdapter) GetType() ScraperType {
	return CinebyType
}
