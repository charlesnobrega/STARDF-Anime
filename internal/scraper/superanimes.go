package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alvarorichard/Goanime/internal/models"
	"github.com/alvarorichard/Goanime/internal/util"
)

const (
	SuperAnimesBase      = "https://superanimes.in"
	SuperAnimesAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	SuperAnimesSearchURL = "https://superanimes.in/?s=%s"
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
	searchURL := fmt.Sprintf(SuperAnimesSearchURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", SuperAnimesAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Cache-Control", "max-age=0")

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
	doc.Find("article, .post, .entry, .post-item, .anime-card").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2.entry-title, h3.title, .post-title, a").Text())
		href, _ := s.Find("a").First().Attr("href")
		img, _ := s.Find("img").First().Attr("src")

		if title == "" || href == "" {
			return
		}

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
	})

	return results, nil
}

func (c *SuperAnimesClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", SuperAnimesAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

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
	doc.Find(".episodios-list a, .episode-list a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := strings.TrimSpace(s.Text())
		num := i + 1

		if re := regexp.MustCompile(`[^\d]*(\d+)[^\d]*`); re.MatchString(title) {
			if match := re.FindStringSubmatch(title); len(match) > 1 {
				fmt.Sscanf(match[1], "%d", &num)
			}
		}

		if href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			episodes = append(episodes, models.Episode{
				Number: fmt.Sprintf("%d", num),
				Num:    num,
				Title:  models.TitleDetails{English: title},
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
