package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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

func (c *CinebyClient) SearchMedia(query string) ([]*models.Anime, error) {
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
		if res.MediaType != "movie" && res.MediaType != "tv" {
			continue
		}

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

func (c *CinebyClient) GetEpisodes(mediaURL string) ([]models.Episode, error) {
	req, err := http.NewRequest("GET", mediaURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", CinebyAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract data from __NEXT_DATA__
	re := regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)
	match := re.FindSubmatch(body)
	if len(match) < 2 {
		return nil, fmt.Errorf("could not find __NEXT_DATA__ in cineby page")
	}

	var nextData struct {
		Props struct {
			PageProps struct {
				Data struct {
					Seasons []struct {
						Episodes []struct {
							ID    int    `json:"id"`
							Num   int    `json:"episode_number"`
							Title string `json:"name"`
						} `json:"episodes"`
						Number int `json:"season_number"`
					} `json:"seasons"`
					MediaType string `json:"media_type"`
					ID        int    `json:"id"`
				} `json:"data"`
			} `json:"pageProps"`
		} `json:"props"`
	}

	if err := json.Unmarshal(match[1], &nextData); err != nil || nextData.Props.PageProps.Data.ID == 0 {
		// Fallback: extract ID and type from URL
		// Example: https://www.cineby.gd/movie/157336
		parts := strings.Split(strings.Trim(mediaURL, "/"), "/")
		if len(parts) >= 2 {
			idStr := parts[len(parts)-1]
			id, _ := strconv.Atoi(idStr)
			mType := parts[len(parts)-2]
			
			if id > 0 {
				var episodes []models.Episode
				if mType == "movie" {
					episodes = append(episodes, models.Episode{
						Number: "Filme",
						Num:    1,
						URL:    fmt.Sprintf("cineby|movie|%d", id),
					})
				} else {
					// For TV, we really need the seasons from JSON, but let's at least return a placeholder
					return nil, fmt.Errorf("failed to extract TV seasons from JSON for %s", mediaURL)
				}
				return episodes, nil
			}
		}
		return nil, fmt.Errorf("failed to parse __NEXT_DATA__ and URL fallback failed for %s", mediaURL)
	}

	var episodes []models.Episode
	data := nextData.Props.PageProps.Data

	if data.MediaType == "movie" || data.Seasons == nil {
		episodes = append(episodes, models.Episode{
			Number: "Filme",
			Num:    1,
			URL:    fmt.Sprintf("cineby|movie|%d", data.ID),
		})
		return episodes, nil
	}

	for _, season := range data.Seasons {
		for _, ep := range season.Episodes {
			episodes = append(episodes, models.Episode{
				Number:   fmt.Sprintf("T%d:E%d - %s", season.Number, ep.Num, ep.Title),
				Num:      ep.Num,
				URL:      fmt.Sprintf("cineby|tv|%d|%d|%d", data.ID, season.Number, ep.Num),
				SeasonID: fmt.Sprintf("%d", season.Number),
			})
		}
	}

	return episodes, nil
}

func (c *CinebyClient) GetStreamURLs(episodeURL string) ([]string, error) {
	// Format: cineby|type|id[|season|episode]
	parts := strings.Split(episodeURL, "|")
	if len(parts) < 3 {
		// Fallback for old URLs
		return c.legacyGetStreamURLs(episodeURL)
	}

	mediaType := parts[1]
	id := parts[2]

	var streamURL string
	if mediaType == "movie" {
		streamURL = fmt.Sprintf("https://db.videasy.net/3/movie/%s", id)
	} else if mediaType == "tv" && len(parts) >= 5 {
		season := parts[3]
		episode := parts[4]
		streamURL = fmt.Sprintf("https://db.videasy.net/3/tv/%s/%s/%s", id, season, episode)
	} else {
		return nil, fmt.Errorf("invalid cineby episode URL: %s", episodeURL)
	}

	// We return the db.videasy.net URL which is an iframe source
	return []string{streamURL}, nil
}

func (c *CinebyClient) legacyGetStreamURLs(movieURL string) ([]string, error) {
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

