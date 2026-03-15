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
	CineGratisBase = "https://cinegratis.tv"
	CineGratisUA   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type CineGratisClient struct {
	client  *http.Client
	baseURL string
}

func NewCineGratisClient() *CineGratisClient {
	return &CineGratisClient{
		client:  util.GetSharedClient(),
		baseURL: CineGratisBase,
	}
}

func (c *CineGratisClient) Search(query string) ([]*models.Anime, error) {
	data := url.Values{}
	data.Set("do", "search")
	data.Set("subaction", "search")
	data.Set("story", query)

	req, err := http.NewRequest("POST", c.baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", CineGratisUA)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cinegratis search failed: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
	doc.Find(".tt").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")
		img, _ := s.Find("img").Attr("src")

		if title != "" && href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			
			mediaType := models.MediaTypeMovie
			if strings.Contains(href, "/series/") {
				mediaType = models.MediaTypeTV
			}

			results = append(results, &models.Anime{
				Name:      title,
				URL:       href,
				ImageURL:  img,
				MediaType: mediaType,
				Source:    "CineGratis",
			})
		}
	})

	return results, nil
}

func (c *CineGratisClient) GetStreamURL(pageURL string) (string, error) {
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", CineGratisUA)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Try to find embedUrl in JSON-LD
	var embedURL string
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		if strings.Contains(content, "embedUrl") {
			// Simple extraction since we don't want to bring in a heavy JSON parser here if avoidable
			// but for robustness we should use strings.Index
			start := strings.Index(content, `"embedUrl":"`)
			if start != -1 {
				start += 12
				end := strings.Index(content[start:], `"`)
				if end != -1 {
					embedURL = content[start : start+end]
					embedURL = strings.ReplaceAll(embedURL, "\\/", "/")
				}
			}
		}
	})

	if embedURL != "" {
		return embedURL, nil
	}

	// Fallback to iframe src
	iframe := doc.Find("iframe[src*='player'], iframe[src*='watch']").First()
	if src, exists := iframe.Attr("src"); exists {
		return src, nil
	}

	return "", fmt.Errorf("could not find stream URL on page")
}

func (c *CineGratisClient) GetEpisodes(seriesURL string) ([]models.Episode, error) {
	req, err := http.NewRequest("GET", seriesURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", CineGratisUA)

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
	// CineGratis usually lists episodes in a structured way
	doc.Find(".episodes-list a, .list-episodes a").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")
		if !strings.HasPrefix(href, "http") {
			href = c.baseURL + href
		}

		episodes = append(episodes, models.Episode{
			Number: title,
			Num:    i + 1,
			URL:    href,
		})
	})

	return episodes, nil
}
