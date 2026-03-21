package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
)

func main() {
	fmt.Println("🔬 --- DIAGNÓSTICO CIRÚRGICO: SOLO LEVELING @ ANIMEFIRE ---")

	sm := scraper.NewScraperManager()
	af, _ := sm.FindScraperByName("AnimeFire")

	// Teste 1: Busca pura
	fmt.Println("\n[TESTE 1: BUSCA]")
	results, _ := af.SearchAnime("Solo Leveling")
	for i, r := range results {
		fmt.Printf("   %d. [%s] URL: %s\n", i+1, r.Name, r.URL)
		
		// Teste 2: Listagem de episódios para cada resultado
		fmt.Printf("      -> Listando episódios para %s...\n", r.Name)
		eps, err := af.GetAnimeEpisodes(r.URL)
		if err != nil {
			fmt.Printf("      ❌ ERRO: %v\n", err)
		} else {
			fmt.Printf("      ✅ SUCESSO: %d episódios encontrados.\n", len(eps))
			if len(eps) > 0 {
				fmt.Printf("         EP 1 URL: %s\n", eps[0].URL)
			}
		}
	}

	fmt.Println("\n🏁 --- FIM DO DIAGNÓSTICO ---")
}
