//go:build ignore
// +build ignore

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type TestResult struct {
	URL        string
	Status     string
	StatusCode int
	Latency    time.Duration
	Error      string
	Headers    http.Header
}

func testConnection(url string) TestResult {
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	start := time.Now()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return TestResult{URL: url, Error: err.Error(), Latency: latency}
	}
	defer resp.Body.Close()

	return TestResult{
		URL:        url,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Latency:    latency,
		Headers:    resp.Header,
	}
}

func main() {
	urls := []string{
		"https://animefire.io",
		"https://betteranime.io",
		"https://betteranime.net",
		"https://animesonlinecc.to",
		"https://goyabu.io",
		"https://cinegratis.tv",
		"https://flixhq.to",
		"https://cineby.gd",
		"https://topanimes.tv",
		"https://animesdigital.org",
	}

	fmt.Println("--- INICIANDO RELATÓRIO DE CONECTIVIDADE ---")
	for _, u := range urls {
		res := testConnection(u)
		fmt.Printf("URL: %s\nStatus: %s (%d)\nLatência: %v\nErro: %s\n", res.URL, res.Status, res.StatusCode, res.Latency, res.Error)
		if res.StatusCode == 301 || res.StatusCode == 302 {
			fmt.Printf("Redirecionamento para: %s\n", res.Headers.Get("Location"))
		}
		fmt.Println("-------------------------------------------")
	}
}
