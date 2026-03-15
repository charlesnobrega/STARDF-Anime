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
	TopAnimesBase      = "https://topanimes.net"
	TopAnimesSearchURL = "https://topanimes.net/?s=%s"
)

type TopAnimesClient struct {
	client  *http.Client
	baseURL string
}

func NewTopAnimesClient() *TopAnimesClient {
	return &TopAnimesClient{
		client:  util.GetScraperClient(),
		baseURL: TopAnimesBase,
	}
}

func (c *TopAnimesClient) SearchAnime(query string) ([]*models.Anime, error) {
	util.RandomDelay(0, 1)

	searchURL := fmt.Sprintf(TopAnimesSearchURL, url.QueryEscape(query))
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
		return nil, fmt.Errorf("topanimes search failed with status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime

	doc.Find("article, .item, .result-item").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a").First()
		href, exists := link.Attr("href")
		if !exists || href == "" || strings.Contains(href, "/category/") {
			return
		}

		title := strings.TrimSpace(link.AttrOr("title", ""))
		if title == "" {
			title = strings.TrimSpace(s.Find(".title, h3, h2").Text())
		}
		if title == "" {
			title = strings.TrimSpace(link.Text())
		}

		img := s.Find("img").AttrOr("src", "")
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
				Source:    "TopAnimes",
				MediaType: models.MediaTypeAnime,
			})
		}
	})

	return results, nil
}

func (c *TopAnimesClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
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
		return nil, fmt.Errorf("topanimes episode list failed: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var episodes []models.Episode
	doc.Find(".episodios li a, .episodes a, ul.episodes li a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		title := strings.TrimSpace(s.Text())
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

func (c *TopAnimesClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
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
	
	// Option 1: find players inside iframes directly on page
	doc.Find("iframe, .player iframe").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if src != "" && !strings.Contains(src, "youtube") {
			finalStreamURL = src
		}
	})

	// Option 2: Some Dooplay themes use player scripts with data
	if finalStreamURL == "" {
		doc.Find(".play-video, .player-video, .player-box").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("data-src")
			if src != "" {
				finalStreamURL = src
			}
		})
	}

	if finalStreamURL == "" {
		return "", nil, fmt.Errorf("no stream URL found in TopAnimes")
	}

	metadata := map[string]string{
		"source":  "topanimes",
		"quality": "default",
	}

	return finalStreamURL, metadata, nil
}
