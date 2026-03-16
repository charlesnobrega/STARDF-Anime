package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

func main() {
	query := "one piece"
	baseURL := "https://goyabu.io"
	searchURL := fmt.Sprintf("https://goyabu.io/?s=%s", url.QueryEscape(query))
	
	client := util.GetScraperClient()
	req, _ := http.NewRequest("GET", searchURL, nil)
	ua := util.UserAgentList()
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", baseURL)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("goyabu_raw2.html", body, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Println("Saved to goyabu_raw2.html")
}
