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
	SuperAnimesBase      = "https://superanimes.in"
	SuperAnimesUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0"
)

type SuperAnimesClient struct {
	client    *http.Client
	baseURL   string
	userAgent string
}

func NewSuperAnimesClient() *SuperAnimesClient {
	return &SuperAnimesClient{
		client:    util.GetFastClient(),
		baseURL:   SuperAnimesBase,
		userAgent: SuperAnimesUserAgent,
	}
}

func (c *SuperAnimesClient) SearchAnime(query string) ([]*models.Anime, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(query))
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	var results []*models.Anime
	doc.Find(".anime-card, .poster").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3, .title").Text())
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
	clean := regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(strings.ToLower(title), "")
	return clean
}