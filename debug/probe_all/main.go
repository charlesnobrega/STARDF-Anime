package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

func main() {
	util.IsDebug = true
	fmt.Println("Probing all scrapers...")
	
	manager := scraper.NewScraperManager()

	query := "Vingadores"
	fmt.Printf("Searching for: %s\n\n", query)
	
	results, err := manager.SearchAnime(query, nil)
	if err != nil {
		fmt.Printf("Manager Search Error: %v\n", err)
	}

	fmt.Printf("\nTotal results found: %d\n", len(results))
	
	// Report results by source
	sources := make(map[string]int)
	for _, res := range results {
		sources[res.Source]++
	}
	
	fmt.Println("\nScraper Status Report:")
	for name, count := range sources {
		fmt.Printf("- %s: %d results\n", name, count)
	}
	
	// We also need to check which ones returned 0 results
	// We can't easily see from 'results' alone if they timed out or just found nothing.
}
