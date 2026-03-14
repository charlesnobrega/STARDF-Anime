package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alvarorichard/Goanime/internal/models"
	"github.com/alvarorichard/Goanime/internal/util"
)

const (
	AnimesOnlineCCBase  = "https://animesonlinecc.to"
	AnimesOnlineCCAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type AnimesOnlineCCClient struct {
	client  *http.Client
	baseURL string
}

func NewAnimesOnlineCCClient() *AnimesOnlineCCClient {
	return &AnimesOnlineCCClient{
		client:  util.GetFastClient(),
		baseURL: AnimesOnlineCCBase,
	}
}

// SearchAnime busca animes no AnimesOnlineCC
func (c *AnimesOnlineCCClient) SearchAnime(query string) ([]*models.Anime, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", AnimesOnlineCCAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("animesonlinecc search failed: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
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
				Source:    "AnimesOnlineCC",
				MediaType: models.MediaTypeAnime,
			})
		}
	})

	return results, nil
}

// GetEpisodes retorna lista de episódios
func (c *AnimesOnlineCCClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", AnimesOnlineCCAgent)

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
	doc.Find(".episodes-list a, .episode-list a, .episodios a, .list-episodes a, .episode-item a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := strings.TrimSpace(s.Text())
		num := i + 1

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

// GetStreamURL retorna URL de streaming
func (c *AnimesOnlineCCClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("User-Agent", AnimesOnlineCCAgent)

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
	doc.Find("iframe, video, .player, .video-container").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.HasPrefix(src, "http") {
			videoURL = src
		}
		if dataSrc, ok := s.Attr("data-src"); ok && strings.HasPrefix(dataSrc, "http") {
			videoURL = dataSrc
		}
	})

	if videoURL == "" {
		// Buscar em scripts
		doc.Find("script").Each(func(i int, s *goquery.Selection) {
			scriptText := s.Text()
			// Procurar URLs de vídeo
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
		"source":  "animesonlinecc",
		"quality": "default",
	}
	return videoURL, metadata, nil
}

// Adapter
type AnimesOnlineCCAdapter struct {
	client *AnimesOnlineCCClient
}

func NewAnimesOnlineCCAdapter(client *AnimesOnlineCCClient) *AnimesOnlineCCAdapter {
	return &AnimesOnlineCCAdapter{client: client}
}

func (a *AnimesOnlineCCAdapter) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	return a.client.SearchAnime(query)
}

func (a *AnimesOnlineCCAdapter) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	return a.client.GetEpisodes(animeURL)
}

func (a *AnimesOnlineCCAdapter) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	return a.client.GetStreamURL(episodeURL)
}

func (a *AnimesOnlineCCAdapter) GetType() ScraperType {
	return AnimesOnlineCCTYPE
}
