package scraper_test

import (
	"testing"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
)

func TestSuperAnimes_LiveSearch(t *testing.T) {
	client := scraper.NewSuperAnimesClient()
	results, err := client.SearchAnime("naruto")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) == 0 {
		t.Logf("Warning: no results found for 'naruto'")
	} else {
		for i, r := range results {
			t.Logf("Result %d: %s (%s)", i+1, r.Name, r.URL)
		}
	}
}
