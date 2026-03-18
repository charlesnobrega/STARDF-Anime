package scraper

import (
	"fmt"
	"math/rand"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

const (
	GoyabuBase      = "https://goyabu.io"
	GoyabuSearchURL = "https://goyabu.io/?s=%s"
)

type GoyabuClient struct {
	client  *http.Client
	baseURL string
}

func NewGoyabuClient() *GoyabuClient {
	return &GoyabuClient{
		client:  util.GetScraperClient(),
		baseURL: GoyabuBase,
	}
}

func (c *GoyabuClient) visitHome() {
	req, _ := http.NewRequest("GET", c.baseURL, nil)
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", c.baseURL)
	resp, err := c.client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

func (c *GoyabuClient) SearchAnime(query string) ([]*models.Anime, error) {
	util.RandomDelay(0, 1)
	c.visitHome()

	searchURL := fmt.Sprintf(GoyabuSearchURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", c.baseURL)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("goyabu search failed: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
	doc.Find("article.boxAN").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".title").Text())
		href, _ := s.Find("a").First().Attr("href")
		img, _ := s.Find("img.cover").First().Attr("src")
		if img == "" {
			img, _ = s.Find("img").First().Attr("src")
		}

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
			Source:    "Goyabu",
			MediaType: models.MediaTypeAnime,
		})
	})

	return results, nil
}

func (c *GoyabuClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	util.RandomDelay(1, 3)

	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", c.baseURL)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := string(bodyBytes)

	// Try extracting from allEpisodes JSON in script tag
	re := regexp.MustCompile(`const allEpisodes = (\[.*?\]);`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		var rawEpisodes []map[string]interface{}
		if err := json.Unmarshal([]byte(match[1]), &rawEpisodes); err == nil {
			var episodes []models.Episode
			for _, ep := range rawEpisodes {
				episodioStr, _ := ep["episodio"].(string)
				num := 0
				fmt.Sscanf(episodioStr, "%d", &num)
				link, _ := ep["link"].(string)
				if !strings.HasPrefix(link, "http") {
					link = c.baseURL + link
				}
				title, _ := ep["episode_name"].(string)
				if title == "" {
					title = fmt.Sprintf("Episódio %s", episodioStr)
				}
				episodes = append(episodes, models.Episode{
					Number: episodioStr,
					Num:    num,
					Title:  models.TitleDetails{English: title},
					URL:    link,
				})
			}
			if len(episodes) > 0 {
				return episodes, nil
			}
		}
	}

	// Fallback to DOM parsing
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	var episodes []models.Episode
	doc.Find(".episodios-list a, .episode-list a, .episode-item a, article.boxEP a").Each(func(i int, s *goquery.Selection) {
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

func (c *GoyabuClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	util.RandomDelay(1, 3)

	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return "", nil, err
	}
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", c.baseURL)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	body := string(bodyBytes)

	// Try extracting from playersData JSON in script tag
	re := regexp.MustCompile(`var playersData = (\[.*?\]);`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		var players []map[string]interface{}
		if err := json.Unmarshal([]byte(match[1]), &players); err == nil && len(players) > 0 {
			videoURL, _ := players[0]["url"].(string)
			if videoURL != "" {
				metadata := map[string]string{
					"source":  "goyabu",
					"quality": "default",
				}
				return videoURL, metadata, nil
			}
		}
	}

	// Fallback to DOM parsing
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", nil, err
	}

	var videoURL string
	doc.Find("iframe, video, .player, .video-container, .embed, #player-preview-image").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.HasPrefix(src, "http") {
			videoURL = src
		}
		if dataSrc, ok := s.Attr("data-src"); ok && strings.HasPrefix(dataSrc, "http") {
			videoURL = dataSrc
		}
		if dataVideo, ok := s.Attr("data-video"); ok && strings.HasPrefix(dataVideo, "http") {
			videoURL = dataVideo
		}
	})

	if videoURL == "" {
		return "", nil, fmt.Errorf("no stream found")
	}

	metadata := map[string]string{
		"source":  "goyabu",
		"quality": "default",
	}
	return videoURL, metadata, nil
}

