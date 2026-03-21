package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
)

func main() {
	fmt.Println("🔍 --- SCAN BIT-A-BIT V2: STARDF-ANIME (2026 RECOVERY) ---")

	sm := scraper.NewScraperManager()
	
	// BIT 1: DASHBOARD MAPPING
	query := "Steel Ball Run"
	fmt.Printf("\n[BIT 1: DASHBOARD] Searching unified sources for: %s\n", query)
	results, _ := sm.SearchAnime(query, nil)
	if len(results) > 0 {
		fmt.Printf("✅ Sucesso: %d fontes encontradas.\n", len(results))
		for _, res := range results {
			fmt.Printf("   - [%s] %s\n", res.Source, res.Name)
		}
	} else {
		fmt.Println("❌ Nenhuma fonte encontrada para Steel Ball Run.")
	}

	// BIT 2: EPISODE LISTING (ANIMEFIRE RECOVERY)
	fmt.Println("\n[BIT 2: LISTING - ANIMEFIRE (Adaptive)]")
	af, _ := sm.FindScraperByName("AnimeFire")
	// This URL needs fallback to -todos-os-episodios
	afURL := "https://animefire.io/animes/yuusha-kei-ni-shosu-choubatsu-yuusha-9004-tai-keimu-kiroku"
	eps, _ := af.GetAnimeEpisodes(afURL)
	if len(eps) > 0 {
		fmt.Printf("✅ AnimeFire Sucesso (Fallback Ativado): %d episódios listados.\n", len(eps))
	} else {
		fmt.Println("❌ AnimeFire ainda falha na listagem.")
	}

	// BIT 3: NEW SOURCES (ANROLL & ANIMESONLINE)
	fmt.Println("\n[BIT 3: NEW SOURCE CHECK]")
	sources := []string{"Anroll", "AnimesOnline", "TopAnimes"}
	for _, name := range sources {
		s, err := sm.FindScraperByName(name)
		if err == nil {
			fmt.Printf("✅ [%s] REGISTRADO E PRONTO.\n", name)
			// Quick test search
			res, _ := s.SearchAnime("Solo Leveling")
			if len(res) > 0 {
				fmt.Printf("   - [%s] Validado: Retornou resultados de busca.\n", name)
			} else {
				fmt.Printf("   - [%s] ⚠️  Sem resultados na busca (Bot Detect?).\n", name)
			}
		} else {
			fmt.Printf("❌ [%s] NÃO LOCALIZADO.\n", name)
		}
	}

	fmt.Println("\n🏁 --- FIM DO SCAN BIT-A-BIT V2 ---")
}
