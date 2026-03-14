package scraper

import (
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
	CinebyUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0"
)

type CinebyClient struct {
	client    *http.Client
	baseURL   string
	userAgent string
}

func NewCinebyClient() *CinebyClient {
	return &CinebyClient{
		client:    util.GetFastClient(),
		baseURL:   CinebyBase,
		userAgent: CinebyUserAgent,
	}
}

func (c *CinebyClient) SearchMovies(query string) ([]*models.Anime, error) {
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
	doc.Find(".movie-item, .film-card").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2, .title").Text())
		href, _ := s.Find("a").First().Attr("href")
		img, _ := s.Find("img").First().Attr("src")
		if title != "" && href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			results = append(results, &models.Anime{
				ID:       generateCinebyID(title),
				Title:    title,
				URL:      href,
				ImageURL: img,
				Type:     models.TypeMovie,
				Source:   "Cineby",
			})
		}
	})
	return results, nil
}

func (c *CinebyClient) GetStreamURLs(movieURL string) ([]string, error) {
	req, _ := http.NewRequest("GET", movieURL, nil)
	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	var streams []string
	doc.Find("iframe, .player, a.watch-button").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.HasPrefix(src, "http") {
			streams = append(streams, src)
		}
		if dataSrc, ok := s.Attr("data-src"); ok && strings.HasPrefix(dataSrc, "http") {
			streams = append(streams, dataSrc)
		}
	})
	return streams, nil
}

func generateCinebyID(title string) string {
	clean := regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(strings.ToLower(title), "")
	return clean
}