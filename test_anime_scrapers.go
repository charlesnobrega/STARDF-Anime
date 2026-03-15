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

	query := "One Piece"

	fmt.Println("--- Testing Animefire ---")
	af := scraper.NewAnimefireClient()
	results, err := af.SearchAnime(query)
	printResults("Animefire", results, err)

	fmt.Println("\n--- Testing Goyabu ---")
	gb := scraper.NewGoyabuClient()
	results, err = gb.SearchAnime(query)
	printResults("Goyabu", results, err)

	fmt.Println("\n--- Testing SuperAnimes ---")
	sa := scraper.NewSuperAnimesClient()
	results, err = sa.SearchAnime(query)
	printResults("SuperAnimes", results, err)

	fmt.Println("\n--- Testing AnimesOnlineCC ---")
	ac := scraper.NewAnimesOnlineCCClient()
	results, err = ac.SearchAnime(query)
	printResults("AnimesOnlineCC", results, err)
}

func printResults(source string, results []*models.Anime, err error) {
	if err != nil {
		fmt.Printf("[%s] ERROR: %v\n", source, err)
		return
	}
	fmt.Printf("[%s] Found %d results\n", source, len(results))
	for i, res := range results {
		if i >= 3 {
			break
		}
		fmt.Printf("  - %s: %s\n", res.Name, res.URL)
	}
}
// Note: models import is missing here, will fix in next step if needed or just use the same import block as auditor
