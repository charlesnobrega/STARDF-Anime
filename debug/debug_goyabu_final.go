package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// Try different common search URL patterns for Goyabu
	urls := []string{
		"https://goyabu.io/?s=One+Piece",
		"https://goyabu.com/?s=One+Piece",
		"https://goyabu.io/search/One+Piece",
		"https://goyabu.io/resultado-busca?s=One+Piece",
	}

	for _, u := range urls {
		fmt.Printf("Testing %s...\n", u)
		req, _ := http.NewRequest("GET", u, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		fmt.Printf("Status: %s\n", resp.Status)
		if resp.StatusCode == 200 {
			body, _ := io.ReadAll(resp.Body)
			filename := fmt.Sprintf("goyabu_test_%d.html", len(u))
			os.WriteFile(filename, body, 0644)
			fmt.Printf("Saved to %s\n", filename)
		}
	}
}
