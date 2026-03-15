package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://superanimes.in/busca/?search_query=naruto"
	fmt.Println("Testing:", url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	
	client := &http.Client{}
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Println("Status:", resp.Status)
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Parse error", err)
		return
	}
	
	title := doc.Find("title").Text()
	fmt.Println("Page Title:", strings.TrimSpace(title))
	
	fmt.Println("\nAll box-anime tags:")
	doc.Find(".box-anime").Each(func(i int, s *goquery.Selection) {
		html, _ := s.Html()
		fmt.Printf("--- box-anime %d HTML ---\n%s\n", i, html)
	})
	
	fmt.Println("\nLooking for 'naruto' in ANY link text:")
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")
		if strings.Contains(strings.ToLower(text), "naruto") || strings.Contains(strings.ToLower(href), "naruto") {
			fmt.Printf("href=%s, text=%s\n", href, text)
			parentClass, _ := s.Parent().Attr("class")
			fmt.Printf("  parent class: %s\n", parentClass)
		}
	})
	
	fmt.Println("\nAll sections/articles:")
	doc.Find("article, section").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		id, _ := s.Attr("id")
		fmt.Printf("tag='%s', id='%s', class='%s'\n", s.Nodes[0].Data, id, class)
		s.Find("a").Slice(0, 3).Each(func(j int, a *goquery.Selection) {
			text := strings.TrimSpace(a.Text())
			fmt.Printf("  sample a: %s\n", text)
		})
	})
}
