package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("superanimes search failed: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
	doc.Find(".anime-card, .anime-item, .card").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3, .title, .name, a > img").AttrOr("alt", ""))
		if title == "" {
			title = strings.TrimSpace(s.Text())
		}
		href, _ := s.Find("a").First().Attr("href")
		img, _ := s.Find("img").First().Attr("src")

		if title != "" && href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			if img != "" && !strings.HasPrefix(img, "http") {
				img = c.baseURL + img
			}
			results = append(results, &models.Anime{
				Name:      title,
				URL:       href,
				ImageURL:  img,
				Source:    "SuperAnimes",
				MediaType: models.MediaTypeAnime,
			})
		}
	})

	return results, nil
}

func (c *SuperAnimesClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	req, err := http.NewRequest("GET", animeURL, nil)
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

	var episodes []models.Episode
	doc.Find(".episodes-list a, .episode-list a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := strings.TrimSpace(s.Text())
		num := i + 1

		if href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			td := models.TitleDetails{English: title}
			episodes = append(episodes, models.Episode{
				Number: fmt.Sprintf("%d", num),
				Num:    num,
				Title:  td,
				URL:    href,
			})
		}
	})

	return episodes, nil
}

func (c *SuperAnimesClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("User-Agent", SuperAnimesAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var videoURL string
	doc.Find("iframe, video, .player, .video-container, .embed").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.HasPrefix(src, "http") {
			videoURL = src
		}
		if dataSrc, ok := s.Attr("data-src"); ok && strings.HasPrefix(dataSrc, "http") {
			videoURL = dataSrc
		}
	})

	if videoURL == "" {
		return "", nil, fmt.Errorf("no stream found")
	}

	metadata := map[string]string{
		"source":  "superanimes",
		"quality": "default",
	}
	return videoURL, metadata, nil
}

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
