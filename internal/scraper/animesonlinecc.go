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
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

const (
	AnimesOnlineCCBase      = "https://animesonlinecc.to"
	AnimesOnlineCCSearchURL = "https://animesonlinecc.to/?s=%s"
)

type AnimesOnlineCCClient struct {
	client  *http.Client
	baseURL string
}

func NewAnimesOnlineCCClient() *AnimesOnlineCCClient {
	return &AnimesOnlineCCClient{
		client:  util.GetScraperClient(),
		baseURL: AnimesOnlineCCBase,
	}
}

func (c *AnimesOnlineCCClient) visitHome() {
	req, _ := http.NewRequest("GET", c.baseURL, nil)
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	resp, err := c.client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

func (c *AnimesOnlineCCClient) SearchAnime(query string) ([]*models.Anime, error) {
	util.RandomDelay(0, 1)
	c.visitHome()

	searchURL := fmt.Sprintf(AnimesOnlineCCSearchURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Referer", c.baseURL+"/")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if util.IsDebug {
		fmt.Printf("AnimesOnlineCC Search Status: %d\n", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("animesonlinecc search failed: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []*models.Anime
	doc.Find("article.item").Each(func(i int, s *goquery.Selection) {
		titleEl := s.Find(".data h3 a")
		if titleEl.Length() == 0 {
			titleEl = s.Find(".poster a").Last()
		}
		title := strings.TrimSpace(titleEl.Text())
		href, _ := titleEl.Attr("href")
		
		imgEl := s.Find("img").First()
		img, _ := imgEl.Attr("src")
		if img == "" {
			img, _ = imgEl.Attr("data-src")
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
			Source:    "AnimesOnlineCC",
			MediaType: models.MediaTypeAnime,
		})
	})

	return results, nil
}

func (c *AnimesOnlineCCClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	util.RandomDelay(1, 3)

	req, err := http.NewRequest("GET", animeURL, nil)
	if err != nil {
		return nil, err
	}
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Referer", c.baseURL+"/")

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
	doc.Find("ul.episodios li, .episodios-list a, .episode-list a").Each(func(i int, s *goquery.Selection) {
		linkEl := s.Find(".episodiotitle a")
		if linkEl.Length() == 0 {
			linkEl = s.Find("a").First()
		}
		
		href, _ := linkEl.Attr("href")
		title := strings.TrimSpace(linkEl.Text())
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

func (c *AnimesOnlineCCClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	util.RandomDelay(1, 3)

	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return "", nil, err
	}
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Referer", c.baseURL+"/")

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
	doc.Find("iframe.metaframe, iframe, video, .player, .video-container, .embed").Each(func(i int, s *goquery.Selection) {
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
		"source":  "animesonlinecc",
		"quality": "default",
	}
	return videoURL, metadata, nil
}

