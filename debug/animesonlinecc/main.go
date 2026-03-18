package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

func main() {
	util.IsDebug = true
	fmt.Println("Testing AnimesOnlineCC Scraper...")
	
	client := scraper.NewAnimesOnlineCCClient()
	
	query := "Naruto Shippuden"
	fmt.Printf("Searching for: %s\n", query)
	results, err := client.SearchAnime(query)
	if err != nil {
		fmt.Printf("Search Error: %v\n", err)
	} else {
		fmt.Printf("Found %d results\n", len(results))
		for i, res := range results {
			if i >= 5 { break }
			fmt.Printf("%d: %s (%s)\n", i+1, res.Name, res.URL)
		}
		
		if len(results) > 0 {
			fmt.Printf("\nGetting episodes for: %s\n", results[0].Name)
			episodes, err := client.GetEpisodes(results[0].URL)
			if err != nil {
				fmt.Printf("Episodes Error: %v\n", err)
			} else {
				fmt.Printf("Found %d episodes\n", len(episodes))
				for i, ep := range episodes {
					if i >= 3 { break }
					fmt.Printf("%d: %s (%s)\n", i+1, ep.Number, ep.URL)
				}
			}
		}
	}
}
