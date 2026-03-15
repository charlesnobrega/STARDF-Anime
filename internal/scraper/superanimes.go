package scraper

import (
	"fmt"
	"math/rand"
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
	SuperAnimesBase      = "https://superanimes.in"
	SuperAnimesSearchURL = "https://superanimes.in/busca/?search_query=%s"
)

type SuperAnimesClient struct {
	client  *http.Client
	baseURL string
}

func NewSuperAnimesClient() *SuperAnimesClient {
	return &SuperAnimesClient{
		client:  util.GetScraperClient(),
		baseURL: SuperAnimesBase,
	}
}

func (c *SuperAnimesClient) visitHome() {
	req, _ := http.NewRequest("GET", c.baseURL, nil)
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")
	resp, err := c.client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

func (c *SuperAnimesClient) SearchAnime(query string) ([]*models.Anime, error) {
	util.RandomDelay(0, 1)
	c.visitHome()

	searchURL := fmt.Sprintf(SuperAnimesSearchURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
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
		return nil, fmt.Errorf("superanimes search failed: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
	doc.Find("div.box-anime").Each(func(i int, s *goquery.Selection) {
		titleEl := s.Find("a.tt")
		title := strings.TrimSpace(titleEl.Text())
		href, _ := titleEl.Attr("href")
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
	util.RandomDelay(1, 3)

	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", c.baseURL)

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
	util.RandomDelay(1, 3)

	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return "", nil, err
	}
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", c.baseURL)

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

