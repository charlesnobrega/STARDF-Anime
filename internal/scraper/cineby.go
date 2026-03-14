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
	CinebyBase  = "https://www.cineby.gd"
	CinebyAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type CinebyClient struct {
	client  *http.Client
	baseURL string
}

func NewCinebyClient() *CinebyClient {
	return &CinebyClient{
		client:  util.GetFastClient(),
		baseURL: CinebyBase,
	}
}

func (c *CinebyClient) SearchMovies(query string) ([]*models.Anime, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", CinebyAgent)

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
	doc.Find(".movie-item, .film-card, .poster").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2, .title, .name").Text())
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
				Type:     models.MediaTypeMovie,
				Source:   "Cineby",
			})
		}
	})
	return results, nil
}

func (c *CinebyClient) GetStreamURLs(movieURL string) ([]string, error) {
	req, err := http.NewRequest("GET", movieURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", CinebyAgent)

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
	doc.Find("iframe, .player, a.watch-button, .stream-link").Each(func(i int, s *goquery.Selection) {
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
	clean := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, strings.ToLower(title))
	return fmt.Sprintf("%s", clean)
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