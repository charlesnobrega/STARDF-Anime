package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alvarorichard/Goanime/internal/util"
)

func fetchAndParse(name, searchURL, containerSel, titleSel, linkSel, imgSel string) {
	fmt.Printf("\n=== %s ===\n", name)
	fmt.Printf("URL: %s\n", searchURL)

	ua := util.UserAgentList()
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("ERROR fetching: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s | Content-Encoding: %s\n", resp.Status, resp.Header.Get("Content-Encoding"))

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		fmt.Printf("Body snippet: %s\n", body)
		return
	}

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Printf("ERROR gzip: %v\n", err)
			return
		}
		defer gzReader.Close()
		reader = gzReader
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Printf("ERROR parsing: %v\n", err)
		return
	}

	total := doc.Find(containerSel).Length()
	fmt.Printf("Containers [%s]: %d\n", containerSel, total)

	count := 0
	doc.Find(containerSel).Each(func(i int, s *goquery.Selection) {
		var title, href string
		if strings.Contains(titleSel, ".tt") {
			titleEl := s.Find("a.tt")
			title = strings.TrimSpace(titleEl.Text())
			href, _ = titleEl.Attr("href")
		} else if titleSel == "h3 a" {
			titleEl := s.Find("h3 a")
			title = strings.TrimSpace(titleEl.Text())
			href, _ = titleEl.Attr("href")
		} else {
			title = strings.TrimSpace(s.Find(titleSel).Text())
			href, _ = s.Find(linkSel).Attr("href")
		}
		img, _ := s.Find(imgSel).Attr("src")
		if title != "" && href != "" {
			count++
			if count <= 3 {
				fmt.Printf("  [%d] %q -> %s\n", i, title, href)
				_ = img
			}
		}
	})
	fmt.Printf("Valid results: %d\n", count)
}

func main() {
	util.IsDebug = false
	util.InitLogger()

	fetchAndParse("SuperAnimes",
		fmt.Sprintf("https://superanimes.in/busca/?search_query=%s", url.QueryEscape("One Piece")),
		"div.box-anime", "a.tt", "a.tt", "img",
	)
	fetchAndParse("AnimesOnlineCC",
		fmt.Sprintf("https://animesonlinecc.to/?s=%s", url.QueryEscape("One Piece")),
		"article.item", "h3 a", "h3 a", "img",
	)
	fetchAndParse("Goyabu",
		fmt.Sprintf("https://goyabu.io/?s=%s", url.QueryEscape("One Piece")),
		"article.boxAN", ".title", "a", "img.cover",
	)
}
