package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

const (
	AnimesDigitalBase      = "https://animesdigital.org"
	AnimesDigitalSearchURL = "https://animesdigital.org/?s=%s"
)

type AnimesDigitalClient struct {
	client  *http.Client
	baseURL string
}

func NewAnimesDigitalClient() *AnimesDigitalClient {
	return &AnimesDigitalClient{
		client:  util.GetScraperClient(),
		baseURL: AnimesDigitalBase,
	}
}

func (c *AnimesDigitalClient) SearchAnime(query string) ([]*models.Anime, error) {
	util.RandomDelay(0, 1)

	searchURL := fmt.Sprintf(AnimesDigitalSearchURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Referer", c.baseURL)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("animesdigital search failed with status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime

	doc.Find(".itemA").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a").First()
		href, exists := link.Attr("href")
		if !exists || href == "" {
			return
		}

		title := strings.TrimSpace(link.AttrOr("title", ""))
		if title == "" {
			title = strings.TrimSpace(s.Find(".tt").Text())
		}
		if title == "" {
			title = strings.TrimSpace(link.Text())
		}

		img := s.Find(".thumb img").AttrOr("src", "")
		if img == "" {
			img = s.Find("img").AttrOr("data-src", "")
		}

		if !strings.HasPrefix(href, "http") {
			href = c.baseURL + href
		}

		if title != "" && href != "" {
			results = append(results, &models.Anime{
				Name:      title,
				URL:       href,
				ImageURL:  img,
				Source:    "AnimesDigital",
				MediaType: models.MediaTypeAnime,
			})
		}
	})

	return results, nil
}

func (c *AnimesDigitalClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	util.RandomDelay(1, 2)

	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", util.UserAgentList())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("animesdigital episode list failed: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var episodes []models.Episode
	doc.Find(".item_ep, .episodios li, .episodes .item").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a.b_flex, a").First()
		href, exists := link.Attr("href")
		if !exists {
			return
		}

		title := strings.TrimSpace(link.Text())
		if title == "" {
			title = strings.TrimSpace(link.AttrOr("title", ""))
		}
		
		num := i + 1

		re := regexp.MustCompile(`(?i)(?:epis[oó]dio|ep)\s*(\d+)`)
		if match := re.FindStringSubmatch(title); len(match) > 1 {
			fmt.Sscanf(match[1], "%d", &num)
		}

		if !strings.HasPrefix(href, "http") {
			href = c.baseURL + href
		}

		episodes = append(episodes, models.Episode{
			Number: fmt.Sprintf("%d", num),
			Num:    num,
			Title:  models.TitleDetails{English: title},
			URL:    href,
		})
	})

	return episodes, nil
}

func (c *AnimesDigitalClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	util.RandomDelay(1, 2)

	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("User-Agent", util.UserAgentList())

	resp, err := c.client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var finalStreamURL string
	
	doc.Find(".player-video iframe, .play-video iframe, iframe").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if src != "" && !strings.Contains(src, "youtube") {
			finalStreamURL = src
		}
	})

	if finalStreamURL == "" {
		return "", nil, fmt.Errorf("no stream URL found in AnimesDigital")
	}

	metadata := map[string]string{
		"source":  "animesdigital",
		"quality": "default",
	}

	return finalStreamURL, metadata, nil
}
