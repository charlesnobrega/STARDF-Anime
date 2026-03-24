//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"strings"
)

func main() {
	fmt.Println("🔬 --- DEBUG DE SUFIXO: SOLO LEVELING @ ANIMEFIRE ---")

	sm := scraper.NewScraperManager()
	af, _ := sm.FindScraperByName("AnimeFire")

	// URL problemática vinda da busca (conforme log anterior)
	testURL := "https://animefire.io/animes/ore-dake-level-up-na-ken-season-2-arise-from-the-shadow-dublado-todos-os-episodios"

	fmt.Printf("[PROVA] Testando URL inicial: %s\n", testURL)

	// Simulação manual da lógica de GetAnimeEpisodes (com fallback)
	cleanURL := strings.TrimRight(testURL, "/")
	hasSuffix := strings.HasSuffix(cleanURL, "-todos-os-episodios")

	fmt.Printf("   - Clean URL: %s\n", cleanURL)
	fmt.Printf("   - Has Suffix '-todos-os-episodios'? %v\n", hasSuffix)

	if !hasSuffix {
		fmt.Printf("   - [FALLBACK SERIA ACIONADO] Novo URL: %s\n", cleanURL+"-todos-os-episodios")
	} else {
		fmt.Println("   - [FALLBACK BLOQUEADO] Suffix já presente.")
	}

	// Execução real
	fmt.Println("\n[EXECUÇÃO REAL]")
	eps, err := af.GetAnimeEpisodes(testURL)
	if err != nil {
		fmt.Printf("❌ ERRO REAL: %v\n", err)
	} else {
		fmt.Printf("✅ SUCESSO: %d episódios.\n", len(eps))
	}

	fmt.Println("\n🏁 --- FIM DO DIAGNÓSTICO ---")
}
