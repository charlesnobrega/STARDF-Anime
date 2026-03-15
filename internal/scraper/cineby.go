package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
	searchURL := fmt.Sprintf("https://db.videasy.net/3/search/multi?language=en&page=1&query=%s", url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", CinebyAgent)
	req.Header.Set("Origin", "https://www.cineby.gd")
	req.Header.Set("Referer", "https://www.cineby.gd/")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cineby search failed: %s", resp.Status)
	}

	var apiResponse struct {
		Results []struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			Title        string `json:"title"`
			MediaType    string `json:"media_type"`
			PosterPath   string `json:"poster_path"`
			ReleaseDate  string `json:"release_date"`
			FirstAirDate string `json:"first_air_date"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode cineby response: %w", err)
	}

	var results []*models.Anime
	for _, res := range apiResponse.Results {
		name := res.Name
		if name == "" {
			name = res.Title
		}
		if name == "" {
			continue
		}

		mediaType := models.MediaTypeTV
		if res.MediaType == "movie" {
			mediaType = models.MediaTypeMovie
		}

		year := res.ReleaseDate
		if year == "" {
			year = res.FirstAirDate
		}
		if len(year) > 4 {
			year = year[:4]
		}

		results = append(results, &models.Anime{
			Name:      name,
			URL:       fmt.Sprintf("%s/%s/%d", CinebyBase, res.MediaType, res.ID),
			ImageURL:  fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", res.PosterPath),
			MediaType: mediaType,
			Source:    "Cineby",
			Year:      year,
		})
	}
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

