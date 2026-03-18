package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

func main() {
	util.IsDebug = true
	fmt.Println("Testing CineGratis Scraper...")
	
	client := scraper.NewCineGratisClient()
	
	query := "John Wick"
	fmt.Printf("Searching for: %s\n", query)
	results, err := client.Search(query)
	if err != nil {
		fmt.Printf("Search Error: %v\n", err)
	} else {
		fmt.Printf("Found %d results\n", len(results))
		for i, res := range results {
			if i >= 5 { break }
			fmt.Printf("%d: %s (%s) Type: %s\n", i+1, res.Name, res.URL, res.MediaType)
		}
		
		if len(results) > 0 {
			fmt.Printf("\nGetting episodes/stream for: %s\n", results[0].Name)
			if results[0].MediaType == "tv" {
				episodes, err := client.GetEpisodes(results[0].URL)
				if err != nil {
					fmt.Printf("Episodes Error: %v\n", err)
				} else {
					fmt.Printf("Found %d episodes\n", len(episodes))
				}
			} else {
				url, err := client.GetStreamURL(results[0].URL)
				if err != nil {
					fmt.Printf("Stream Error: %v\n", err)
				} else {
					fmt.Printf("Stream URL: %s\n", url)
				}
			}
		}
	}
}
