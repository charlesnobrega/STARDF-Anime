//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"strings"
)

func main() {
	fmt.Println("🔬 --- DIAGNÓSTICO ESTRUTURAL V4: VALIDAÇÃO DE STREAM ---")

	sm := scraper.NewScraperManager()
	af, _ := sm.FindScraperByName("AnimeFire")

	query := "Solo Leveling"
	fmt.Printf("[TESTE 1: BUSCA] Query: %s\n", query)

	results, err := af.SearchAnime(query)
	if err != nil {
		fmt.Printf("❌ ERRO NA BUSCA: %v\n", err)
		return
	}

	for i, anime := range results {
		if i >= 3 {
			break
		} // Test only top 3 results
		fmt.Printf("\n   %d. [%s] URL: %s\n", i+1, anime.Name, anime.URL)

		// 1. List Episodes
		eps, err := af.GetAnimeEpisodes(anime.URL)
		if err != nil {
			fmt.Printf("      ❌ ERRO EPISÓDIOS: %v\n", err)
			continue
		}
		fmt.Printf("      ✅ EPISÓDIOS: %d encontrados.\n", len(eps))

		if len(eps) > 0 {
			// 2. Validate Stream Extraction (CRITICAL FIX)
			fmt.Printf("      🔍 Testando extração de stream (EP 1): %s\n", eps[0].URL)
			streamURL, metadata, err := af.GetStreamURL(eps[0].URL)
			if err != nil {
				fmt.Printf("         ❌ ERRO STREAM: %v\n", err)
			} else {
				fmt.Printf("         ✅ SUCESSO! Stream: %s\n", streamURL)
				fmt.Printf("         📦 Metadados: %v\n", metadata)

				// Final verification of domains
				if strings.Contains(streamURL, "lightspeedst.net") || strings.Contains(streamURL, ".mp4") || strings.Contains(streamURL, "m3u8") {
					fmt.Println("         📌 VEREDITO: Link de vídeo VÁLIDO e operacional.")
				} else {
					fmt.Println("         ⚠️  AVISO: Link de vídeo com formato não convencional.")
				}
			}
		}
	}

	fmt.Println("\n🏁 --- FIM DO DIAGNÓSTICO V4 ---")
}
