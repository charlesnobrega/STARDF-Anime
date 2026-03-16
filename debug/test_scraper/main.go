package main

import (
	"fmt"
	"log"

	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
)

func main() {
	c := scraper.NewGoyabuClient()
	animes, err := c.SearchAnime("one piece")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d animes\n", len(animes))
	for _, a := range animes {
		fmt.Printf("- %s: %s (Img: %s)\n", a.Name, a.URL, a.ImageURL)
	}
}
