package scraper_test

import (
	"testing"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
)

func TestNewScrapers_LiveSearch(t *testing.T) {
	query := "naruto"

	baClient := scraper.NewBetterAnimeClient()
	baRes, err := baClient.SearchAnime(query)
	t.Logf("BetterAnime: err=%v, results=%d", err, len(baRes))

	taClient := scraper.NewTopAnimesClient()
	taRes, err := taClient.SearchAnime(query)
	t.Logf("TopAnimes: err=%v, results=%d", err, len(taRes))

	adClient := scraper.NewAnimesDigitalClient()
	adRes, err := adClient.SearchAnime(query)
	t.Logf("AnimesDigital: err=%v, results=%d", err, len(adRes))
}
