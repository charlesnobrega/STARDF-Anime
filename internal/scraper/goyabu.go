package scraper

import (
	"fmt"
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
	GoyabuBase      = "https://goyabu.io"
	GoyabuUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type GoyabuClient struct {
	client    *http.Client
	baseURL   string
	userAgent string
}

func NewGoyabuClient() *GoyabuClient {
	return &GoyabuClient{
		client:    util.GetFastClient(),
		baseURL:   GoyabuBase,
		userAgent: GoyabuUserAgent,
	}
}

// SearchAnime busca animes no Goyabu
func (c *GoyabuClient) SearchAnime(query string) ([]*models.Anime, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

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
	// Tentar múltiplos seletores
	doc.Find(".anime-card, .anime-item, .card, .poster, .item").Each(func(i int, s *goquery.Selection) {
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
				Source:    "Goyabu",
				MediaType: models.MediaTypeAnime,
			})
		}
	})

	return results, nil
}

// GetEpisodes retorna lista de episódios
func (c *GoyabuClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

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
	doc.Find(".episodes-list a, .episode-list a, .episodios a, .list-episodes a, .episode-item a, .ep-link").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := strings.TrimSpace(s.Text())
		num := i + 1

		// Extrair número do episódio do texto ou URL
		if title != "" {
			if re := regexp.MustCompile(`[^\d]*(\d+)[^\d]*`); re.MatchString(title) {
				if match := re.FindStringSubmatch(title); len(match) > 1 {
					fmt.Sscanf(match[1], "%d", &num)
				}
			}
		}

		if href != "" {
			if !strings.HasPrefix(href, "http") {
				href = c.baseURL + href
			}
			episodes = append(episodes, models.Episode{
				Number: fmt.Sprintf("%d", num),
				Num:    num,
				Title:  title,
				URL:    href,
			})
		}
	})

	return episodes, nil
}

// GetStreamURL retorna URL de streaming para um episódio
func (c *GoyabuClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

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
	quality := "default"

	// Procurar players comuns
	doc.Find("iframe, video, .player, .video-container, .embed, .stream, .watch-button").Each(func(i int, s *goquery.Selection) {
		if videoURL != "" {
			return
		}
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

	// Se não encontrou, buscar em scripts
	if videoURL == "" {
		doc.Find("script").Each(func(i int, s *goquery.Selection) {
			scriptText := s.Text()
			// Procurar URLs .mp4 ou .m3u8
			re := regexp.MustCompile(`https?://[^\s"']+\.(mp4|m3u8)[^\s"']*`)
			if match := re.FindString(scriptText); match != "" {
				videoURL = match
			}
		})
	}

	if videoURL == "" {
		return "", nil, fmt.Errorf("no stream found")
	}

	metadata := map[string]string{
		"source":  "goyabu",
		"quality": quality,
	}
	return videoURL, metadata, nil
}

// generateID cria ID único
func generateGoyabuID(title string) string {
	clean := regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(strings.ToLower(title), "")
	return clean
}

// Adapter
type GoyabuAdapter struct {
	client *GoyabuClient
}

func NewGoyabuAdapter(client *GoyabuClient) *GoyabuAdapter {
	return &GoyabuAdapter{client: client}
}

func (a *GoyabuAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *GoyabuAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *GoyabuAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *GoyabuAdapter) GetType() ScraperType {
	return GoyabuType
}
