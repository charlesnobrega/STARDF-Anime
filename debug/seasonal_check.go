package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

func main() {
    util.IsDebug = true
	manager := scraper.NewScraperManager()
	
	targets := []string{
		"STEEL BALL RUN",
		"JUJUTSU KAISEN Season 3",
		"OSHI NO KO Season 3",
	}

	fmt.Println("--- VALIDAÇÃO DE LANÇAMENTOS NOS SCRAPERS ---")
	for _, target := range targets {
		fmt.Printf("\nBuscando: %s\n", target)
		results, err := manager.SearchAnime(target, nil)
		if err != nil {
			fmt.Printf("[❌] Erro geral: %v\n", err)
			continue
		}

		if len(results) == 0 {
			fmt.Println("[⚠️] NENHUM RESULTADO ENCONTRADO EM NENHUM SITE")
			continue
		}

		sourcesFound := make(map[string]int)
		for _, res := range results {
			sourcesFound[res.Source]++
		}

		for src, count := range sourcesFound {
			fmt.Printf("[✅] %s: %d resultados\n", src, count)
		}
	}
}
