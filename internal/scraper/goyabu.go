package scraper

import (
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
	GoyabuBase  = "https://goyabu.io"
	GoyabuAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type GoyabuClient struct {
	client  *http.Client
	baseURL string
}

func NewGoyabuClient() *GoyabuClient {
	return &GoyabuClient{
		client:  util.GetFastClient(),
		baseURL: GoyabuBase,
	}
}

func (c *GoyabuClient) SearchAnime(query string) ([]*models.Anime, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", GoyabuAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
	doc.Find(".anime-item, .card, .poster").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3, .title, .name").Text())
		href, _ := s.Find("a").First().Attr("href")
		img, _ := s.Find("img").First().Attr("src")

		if title != "" && href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			results = append(results, &models.Anime{
				ID:       generateGoyabuID(title),
				Title:    title,
				URL:      href,
				ImageURL: img,
				Source:   "Goyabu",
			})
		}
	})
	return results, nil
}

func (c *GoyabuClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *GoyabuClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	return "", nil, fmt.Errorf("not implemented")
}

func generateGoyabuID(title string) string {
	clean := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, strings.ToLower(title))
	return clean
}

// Adapter
type GoyabuAdapter struct {
	client *GoyabuClient
}

func NewGoyabuAdapter(client *GoyabuClient) *GoyabuAdapter {
	return &GoyabuAdapter{client: client}
}

func (a *GoyabuAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *GoyabuAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *GoyabuAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *GoyabuAdapter) GetType() ScraperType {
	return GoyabuType
}