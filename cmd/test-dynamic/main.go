package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

func main() {
	// Enable Debug to see the Dynamic Loading logs
	util.IsDebug = true
	
	fmt.Println("=== Testando Scraper Manager Inteligente ===")
	sm := scraper.NewScraperManager()
	
	query := "Naruto"
	fmt.Printf("Buscando por: %s...\n", query)
	
	results, err := sm.SearchAnime(query, nil)
	if err != nil {
		fmt.Printf("Erro na busca: %v\n", err)
		return
	}
	
	fmt.Printf("\nResultados Encontrados (%d):\n", len(results))
	for i, anime := range results {
		if i >= 10 { // Limit to 10 for log brevity
			break
		}
		fmt.Printf("[%d] %s (%s) -> %s\n", i+1, anime.Name, anime.Source, anime.URL)
	}
}
