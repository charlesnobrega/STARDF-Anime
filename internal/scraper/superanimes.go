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
	SuperAnimesBase  = "https://superanimes.in"
	SuperAnimesAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type SuperAnimesClient struct {
	client  *http.Client
	baseURL string
}

func NewSuperAnimesClient() *SuperAnimesClient {
	return &SuperAnimesClient{
		client:  util.GetFastClient(),
		baseURL: SuperAnimesBase,
	}
}

func (c *SuperAnimesClient) SearchAnime(query string) ([]*models.Anime, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", SuperAnimesAgent)

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
	doc.Find(".anime-card, .poster, .item").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3, .title, .name").Text())
		href, _ := s.Find("a").First().Attr("href")
		img, _ := s.Find("img").First().Attr("src")

		if title != "" && href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			results = append(results, &models.Anime{
				ID:       generateSuperAnimesID(title),
				Title:    title,
				URL:      href,
				ImageURL: img,
				Source:   "SuperAnimes",
			})
		}
	})
	return results, nil
}

func (c *SuperAnimesClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *SuperAnimesClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	return "", nil, fmt.Errorf("not implemented")
}

func generateSuperAnimesID(title string) string {
	clean := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, strings.ToLower(title))
	return clean
}

// Adapter
type SuperAnimesAdapter struct {
	client *SuperAnimesClient
}

func NewSuperAnimesAdapter(client *SuperAnimesClient) *SuperAnimesAdapter {
	return &SuperAnimesAdapter{client: client}
}

func (a *SuperAnimesAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *SuperAnimesAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *SuperAnimesAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *SuperAnimesAdapter) GetType() ScraperType {
	return SuperAnimesType
}