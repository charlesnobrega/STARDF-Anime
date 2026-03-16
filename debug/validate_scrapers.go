package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
)

func main() {
	// Initialize utilities
	util.IsDebug = true
	
	query := "One Piece"
	fmt.Printf("Validating scrapers for query: %s\n", query)
	fmt.Println("--------------------------------------------------")

	scrapers := []struct {
		name    string
		search  func(string) ([]*models.Anime, error)
		episodes func(string) ([]models.Episode, error)
	}{
		{"AnimeFire", func(q string) ([]*models.Anime, error) { return scraper.NewAnimefireClient().SearchAnime(q) }, nil},
		{"BetterAnime", func(q string) ([]*models.Anime, error) { return scraper.NewBetterAnimeClient().SearchAnime(q) }, func(u string) ([]models.Episode, error) { return scraper.NewBetterAnimeClient().GetEpisodes(u) }},
		{"TopAnimes", func(q string) ([]*models.Anime, error) { return scraper.NewTopAnimesClient().SearchAnime(q) }, func(u string) ([]models.Episode, error) { return scraper.NewTopAnimesClient().GetEpisodes(u) }},
		{"AnimesDigital", func(q string) ([]*models.Anime, error) { return scraper.NewAnimesDigitalClient().SearchAnime(q) }, func(u string) ([]models.Episode, error) { return scraper.NewAnimesDigitalClient().GetEpisodes(u) }},
		{"CineGratis", func(q string) ([]*models.Anime, error) { return scraper.NewCineGratisClient().Search(q) }, func(u string) ([]models.Episode, error) { return scraper.NewCineGratisClient().GetEpisodes(u) }},
		{"FlixHQ", func(q string) ([]*models.Anime, error) { return scraper.NewFlixHQClient().SearchMedia(q) }, nil},
	}

	for _, s := range scrapers {
		fmt.Printf("Testing %s... ", s.name)
		results, err := s.search(query)
		if err != nil {
			fmt.Printf("[ERROR SEARCH] %v\n", err)
			continue
		}
		if len(results) == 0 {
			fmt.Printf("[EMPTY SEARCH]\n")
			continue
		}
		fmt.Printf("[OK SEARCH: %d results] ", len(results))

		if s.episodes != nil {
			eps, err := s.episodes(results[0].URL)
			if err != nil {
				fmt.Printf("[ERROR EPS] %v\n", err)
			} else if len(eps) == 0 {
				fmt.Printf("[EMPTY EPS]\n")
			} else {
				fmt.Printf("[OK EPS: %d episodes]\n", len(eps))
			}
		} else {
			fmt.Println()
		}
	}
}
