//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
)

func main() {
	fmt.Println("🔍 --- SCAN BIT-A-BIT: SISTEMA STARDF-ANIME ---")

	sm := scraper.NewScraperManager()

	// BIT 1: DASHBOARD SEARCH MATCH (JJK S3)
	fmt.Println("\n[BIT 1: DASHBOARD MATCH]")
	query := "Jujutsu Kaisen Season 3"
	fmt.Printf("Searching unified sources for: %s\n", query)
	results, err := sm.SearchAnime(query, nil)
	if err != nil {
		fmt.Printf("❌ Falha na Busca: %v\n", err)
	} else if len(results) == 0 {
		fmt.Println("❌ Nenhuma fonte encontrada para JJK S3 via Scan.")
	} else {
		fmt.Printf("✅ Sucesso: %d fontes encontradas para JJK S3.\n", len(results))
		for _, res := range results {
			fmt.Printf("   - [%s] %s\n", res.Source, res.Name)
		}
	}

	// BIT 2: EPISODE LISTING (ANIMEFIRE)
	fmt.Println("\n[BIT 2: LISTING EPISODES - ANIMEFIRE]")
	af, _ := sm.FindScraperByName("AnimeFire")
	// Using a known URL for test
	afURL := "https://animefire.plus/animes/yuusha-kei-ni-shosu-choubatsu-yuusha-9004-tai-keimu-kiroku"
	eps, err := af.GetAnimeEpisodes(afURL)
	if err != nil {
		fmt.Printf("❌ AnimeFire Error: %v\n", err)
	} else if len(eps) == 0 {
		fmt.Println("❌ AnimeFire retornou 0 episódios (Falha no Seletor .lEp?)")
	} else {
		fmt.Printf("✅ AnimeFire Sucesso: %d episódios listados.\n", len(eps))
	}

	// BIT 3: STREAM EXTRACTION (ANIMEPLAYER)
	fmt.Println("\n[BIT 3: STREAM EXTRACTION - ANIMEPLAYER]")
	ap, _ := sm.FindScraperByName("AnimePlayer")
	// Search first to get a real episode URL
	apSearch, _ := ap.SearchAnime("Solo Leveling") // Using Solo Leveling for extraction speed test
	if len(apSearch) > 0 {
		apEps, _ := ap.GetAnimeEpisodes(apSearch[0].URL)
		if len(apEps) > 0 {
			stream, meta, err := ap.GetStreamURL(apEps[0].URL)
			if err != nil {
				fmt.Printf("❌ AnimePlayer Stream Error: %v\n", err)
			} else {
				fmt.Printf("✅ AnimePlayer Sucesso: Extrator capturou mídia!\n")
				fmt.Printf("   - Final URL: %s...\n", stream[:50])
				fmt.Printf("   - Source: %s\n", meta["source"])
			}
		} else {
			fmt.Println("❌ Falha ao listar EPs no AnimePlayer para teste de stream.")
		}
	} else {
		fmt.Println("❌ Falha na busca base para teste de extração.")
	}

	fmt.Println("\n🏁 --- FIM DO SCAN BIT-A-BIT ---")
}
