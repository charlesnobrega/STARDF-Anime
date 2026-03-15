package main

import (
	"fmt"
	"github.com/alvarorichard/Goanime/internal/models"
	"github.com/alvarorichard/Goanime/internal/scraper"
	"github.com/alvarorichard/Goanime/internal/util"
)

func main() {
	util.IsDebug = true
	util.InitLogger()
	util.PerfEnabled = true
	manager := scraper.NewScraperManager()
	
	// Test Anime Search
	fmt.Println("=== Testing Anime Search (Query: 'One Piece') ===")
	util.GlobalMediaType = "anime"
	animeResults, err := manager.SearchAnime("One Piece", nil)
	if err != nil {
		fmt.Printf("Error searching anime: %v\n", err)
	} else {
		fmt.Printf("Total Anime Results: %d\n", len(animeResults))
		displayTopResults(animeResults, 10)
	}

	fmt.Println("\n=== Testing Movie/TV Search (Query: 'Spider-Man') ===")
	util.GlobalMediaType = "movie"
	movieResults, err := manager.SearchAnime("Spider-Man", nil)
	if err != nil {
		fmt.Printf("Error searching movies: %v\n", err)
	} else {
		fmt.Printf("Total Movie Results: %d\n", len(movieResults))
		displayTopResults(movieResults, 10)
	}
}

func displayTopResults(results []*models.Anime, limit int) {
	for i, res := range results {
		if i >= limit {
			break
		}
		fmt.Printf("[%s] %s (%s) - URL: %s\n", res.Source, res.Name, res.Year, res.URL)
	}
}
