package main

import (
	"fmt"
	"log"
	"github.com/charlesnobrega/STARDF-Anime/internal/anilist"
)

func main() {
	client := anilist.NewClient()
	trending, err := client.GetTrendingSeason(1)
	if err != nil {
		log.Fatalf("Error fetching trending: %v", err)
	}

	fmt.Println("Trending Animes for this Season:")
	for i, m := range trending {
		if i >= 3 { break }
		fmt.Printf("%d. %s (ID: %d)\n", i+1, m.Title.Romaji, m.ID)
	}
}
