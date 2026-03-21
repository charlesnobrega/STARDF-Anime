package scraper

import (
	"encoding/base64"
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
	BetterAnimeBase      = "https://betteranime.io"
	BetterAnimeSearchURL = "https://betteranime.io/pesquisa?titulo=%s"
)

type BetterAnimeClient struct {
	client  *http.Client
	baseURL string
}

func NewBetterAnimeClient() *BetterAnimeClient {
	return &BetterAnimeClient{
		client:  util.GetScraperClient(),
		baseURL: BetterAnimeBase,
	}
}

func (c *BetterAnimeClient) SearchAnime(query string) ([]*models.Anime, error) {
	util.RandomDelay(0, 1)

	searchURL := fmt.Sprintf(BetterAnimeSearchURL, url.QueryEscape(query))
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
		return nil, fmt.Errorf("betteranime search failed with status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime

	// Parse search results based on typical DooPlay or similar WP themes 
	// Commonly found in article.item or div.box or similar
	doc.Find("article, .item, .anime-list-item").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a").First()
		href, exists := link.Attr("href")
		if !exists || href == "" || strings.Contains(href, "/generos/") || strings.Contains(href, "/ano/") {
			return
		}

		title := strings.TrimSpace(link.AttrOr("title", ""))
		if title == "" {
			title = strings.TrimSpace(s.Find("h3, h2, .title").Text())
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

		if title != "" && href != "" && !strings.Contains(href, "author/") {
			results = append(results, &models.Anime{
				Name:      title,
				URL:       href,
				ImageURL:  img,
				Source:    "BetterAnime",
				MediaType: models.MediaTypeAnime,
			})
		}
	})

	return results, nil
}

func (c *BetterAnimeClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
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
		return nil, fmt.Errorf("betteranime episode list failed with status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var episodes []models.Episode
	// Looking for episode links (commonly .episodios li a, .episodes a, etc)
	doc.Find(".episodes a, .episodios li a, .list-episodes a, ul.episodes.range li a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		title := strings.TrimSpace(s.Text())
		num := i + 1

		// Try to extract number from title
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

	// Sometimes episodes are listed newest first, so we might want to reverse them
	// but let's keep them as parsed for now.

	return episodes, nil
}

func (c *BetterAnimeClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
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

	htmlContent, _ := doc.Html()
	var finalStreamURL string

	// BetterAnime usually puts the stream source token in an iframe or a data attribute
	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if strings.Contains(src, "source=") || strings.Contains(src, "token=") {
			finalStreamURL = src
		}
	})

	if finalStreamURL == "" {
		// Try to find raw player scripts
		// Search for base64 encoded strings typically near player instantiation
		reSource := regexp.MustCompile(`(?:source|file|url)["']?\s*:\s*["']([^"']+)["']`)
		matches := reSource.FindAllStringSubmatch(htmlContent, -1)
		for _, m := range matches {
			val := m[1]
			// if it looks like mp4 or m3u8
			if strings.Contains(val, ".mp4") || strings.Contains(val, ".m3u8") {
				finalStreamURL = val
				break
			}
			// if it looks like base64
			if len(val) > 20 && !strings.Contains(val, " ") {
				decoded, err := base64.StdEncoding.DecodeString(val)
				if err == nil && (strings.Contains(string(decoded), "http") || strings.Contains(string(decoded), ".mp4")) {
					finalStreamURL = string(decoded)
					break
				}
			}
		}
	} else {
		// If we found an iframe url with a token/source
		u, err := url.Parse(finalStreamURL)
		if err == nil {
			sourceBase64 := u.Query().Get("source")
			if sourceBase64 == "" {
				sourceBase64 = u.Query().Get("token")
			}
			if sourceBase64 != "" {
				// The logic might be a simple base64 decode
				decoded, err := base64.StdEncoding.DecodeString(sourceBase64)
				if err == nil {
					finalStreamURL = string(decoded)
				}
			}
		}
	}

	if finalStreamURL == "" {
		// Even more aggressive fallback for any stream URL
		reGeneric := regexp.MustCompile(`https?://[^"']+\.(m3u8|mp4)[^"']*`)
		if matches := reGeneric.FindAllString(htmlContent, -1); len(matches) > 0 {
			for _, m := range matches {
				if !strings.Contains(m, "placeholder") && !strings.Contains(m, "ad") {
					finalStreamURL = m
					break
				}
			}
		}
	}

	if finalStreamURL == "" {
		return "", nil, fmt.Errorf("no stream URL found in BetterAnime")
	}

	metadata := map[string]string{
		"source":  "betteranime",
		"quality": "default",
	}

	return finalStreamURL, metadata, nil
}
