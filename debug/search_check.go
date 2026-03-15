package main

import (
	"fmt"
	"github.com/alvarorichard/Goanime/internal/scraper"
	"github.com/alvarorichard/Goanime/internal/util"
)

func main() {
	util.IsDebug = true
	util.InitLogger()
	
	manager := scraper.NewScraperManager()
	
	// Test specifically Goyabu
	fmt.Println("\n--- Testing Goyabu specifically ---")
	gType := scraper.GoyabuType
	results, err := manager.SearchAnime("Naruto", &gType)
	if err != nil {
		fmt.Printf("Goyabu search failed: %v\n", err)
	} else {
		fmt.Printf("Goyabu found %d results\n", len(results))
	}

	// Test specifically AnimesOnlineCC
	fmt.Println("\n--- Testing AnimesOnlineCC specifically ---")
	aType := scraper.AnimesOnlineCCTYPE
	results, err = manager.SearchAnime("Naruto", &aType)
	if err != nil {
		fmt.Printf("AnimesOnlineCC search failed: %v\n", err)
	} else {
		fmt.Printf("AnimesOnlineCC found %d results\n", len(results))
	}
}
