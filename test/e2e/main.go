package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charlesnobrega/STARDF-Anime/internal/anilist"
	"github.com/charlesnobrega/STARDF-Anime/internal/api/movie"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

type TestItem struct {
	Name string
	Type string
}

func main() {
	util.IsDebug = true
	util.InitLogger()

	fmt.Println("🚀 Iniciando E2E Integrado (V2) - Lançamentos da Temporada")

	var testMatrix []TestItem

	// 1. Pegar lançamentos de Anime do AniList
	aniClient := anilist.NewClient()
	fmt.Println("Consultando AniList por animes da temporada atual...")
	trends, err := aniClient.GetTrendingSeason(1)
	if err != nil {
		fmt.Printf("❌ Falha no AniList: %v\n", err)
	} else {
		count := 0
		for _, t := range trends {
			if count >= 3 {
				break
			}
			name := t.Title.Romaji
			if t.Title.English != "" {
				name = t.Title.English
			}
			testMatrix = append(testMatrix, TestItem{Name: name, Type: "anime"})
			count++
		}
	}

	// 2. Pegar lançamentos de Filmes e Séries do TMDB (se configurado) ou Fallback
	tmdbClient := movie.NewTMDBClient()
	if tmdbClient.IsConfigured() {
		fmt.Println("Consultando TMDB por filmes e séries populares...")
		
		movies, _ := tmdbClient.GetPopular("movie")
		if movies != nil && len(movies.Results) >= 3 {
			for i := 0; i < 3; i++ {
				title := movies.Results[i].Title
				if title == "" {
					title = movies.Results[i].Name
				}
				testMatrix = append(testMatrix, TestItem{Name: title, Type: "movie"})
			}
		}

		tvs, _ := tmdbClient.GetPopular("tv")
		if tvs != nil && len(tvs.Results) >= 3 {
			for i := 0; i < 3; i++ {
				title := tvs.Results[i].Name
				if title == "" {
					title = tvs.Results[i].Title
				}
				testMatrix = append(testMatrix, TestItem{Name: title, Type: "series"})
			}
		}
	} else {
		fmt.Println("TMDB (Filmes/Séries Reais) não configurado, usando hardcoded...")
		testMatrix = append(testMatrix, 
			TestItem{Name: "The Matrix", Type: "movie"},
			TestItem{Name: "Inception", Type: "movie"},
			TestItem{Name: "Avatar", Type: "movie"},
			TestItem{Name: "Breaking Bad", Type: "series"},
			TestItem{Name: "The Boys", Type: "series"},
			TestItem{Name: "Game of Thrones", Type: "series"},
		)
	}

	report := "# StarDF-Anime E2E Integrity Report (Dynamic Trends)\n\n"
	report += "| Content | Type | Search | Episodes | Stream | Metadata | Status |\n"
	report += "| :--- | :--- | :--- | :--- | :--- | :--- | :--- |\n"

	manager := scraper.NewScraperManager()

	for _, item := range testMatrix {
		util.GlobalMediaType = item.Type
		fmt.Printf("\n--- Testando [%s] (%s) ---\n", item.Name, item.Type)

		searchStatus := "❌"
		epsStatus := "❌"
		streamStatus := "❌"
		metaStatus := "❌"

		// 1. Busca nos Scrapers Estáticos e Dinâmicos
		results, err := manager.SearchAnime(item.Name, nil)
		var targetAnime *models.Anime
		if err == nil && len(results) > 0 {
			searchStatus = "✅"
			fmt.Printf("✔ Busca: %d resultados (Melhor fonte: %s)\n", len(results), results[0].Source)
			targetAnime = results[0]

			s, errScraper := manager.FindScraperByName(targetAnime.Source)
			if errScraper == nil {
				// 2. Extrapolar Episódios
				eps, errEps := s.GetAnimeEpisodes(targetAnime.URL)
				
				// Fix para filmes na FlixHQ onde a página pode retornar err mas na vdd é single map
				if item.Type == "movie" && strings.Contains(targetAnime.Source, "FlixHQ") && len(eps) == 0 {
					eps = []models.Episode{{URL: targetAnime.URL}}
				}

				if len(eps) > 0 {
					epsStatus = "✅"
					fmt.Printf("✔ Episódios: Encontrados %d episódios\n", len(eps))

					// 3. Pegar URL de Stream
					streamURL, metadata, errStream := s.GetStreamURL(eps[0].URL)
					if errStream == nil && streamURL != "" {
						streamStatus = "✅"
						fmt.Printf("✔ Stream: Resolvido [%s] (%s)\n", metadata["source"], streamURL[:15]+"...")
						
						// 4. Testar Fluxo de Gerar Metadados (.nfo) - Simular Processo
						// Isso validaria util.SyncMetadata
						targetAnime.AnilistID = 12345
						tempDir := "temp_download_e2e"
						os.MkdirAll(tempDir, 0755)
						
						metaErr := util.SyncMetadata(tempDir, targetAnime)
						if metaErr == nil {
							metaStatus = "✅"
							fmt.Printf("✔ Metadados: Arquivo .nfo e JSON gerados em %s\n", tempDir)
						} else {
							fmt.Printf("❌ Metadados Falhou: %v\n", metaErr)
						}
					} else {
						fmt.Printf("❌ Stream Falhou: %v\n", errStream)
					}
				} else {
					fmt.Printf("❌ Episódios Falharam: %v\n", errEps)
				}
			} else {
				fmt.Printf("❌ Falha ao inicializar scraper %s: %v\n", targetAnime.Source, errScraper)
			}
		} else {
			fmt.Printf("❌ Busca Falhou: Nenhuma fonte retornou %s\n", item.Name)
		}

		overall := "FAIL"
		if searchStatus == "✅" && epsStatus == "✅" && streamStatus == "✅" {
			overall = "PASS"
		}

		report += fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |\n",
			item.Name, item.Type, searchStatus, epsStatus, streamStatus, metaStatus, overall)
	}

	reportPath := "test_integrity_report_v2.md"
	errWrite := os.WriteFile(reportPath, []byte(report), 0644)
	if errWrite != nil {
		log.Fatal(errWrite)
	}

	fmt.Printf("\n🎉 E2E V2 Test completed! Relatório salvo em %s\n", reportPath)
}
