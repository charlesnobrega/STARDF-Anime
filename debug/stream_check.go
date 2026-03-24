//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
)

func main() {
	fmt.Println("🔬 --- DEBUG DE STREAM: SENTENCED TO BE A HERO @ ANIMEFIRE ---")

	sm := scraper.NewScraperManager()
	af, _ := sm.FindScraperByName("AnimeFire")

	// URL de um episódio válido do diagnóstico anterior
	epURL := "https://animefire.io/animes/yuusha-kei-ni-shosu-choubatsu-yuusha-9004-tai-keimu-kiroku-todos-os-episodios/1"

	fmt.Printf("[TESTE] Extraindo stream de: %s\n", epURL)

	stream, metadata, err := af.GetStreamURL(epURL)
	if err != nil {
		fmt.Printf("❌ ERRO: %v\n", err)
	} else {
		fmt.Printf("✅ SUCESSO! Stream URL: %s\n", stream)
		fmt.Printf("   Metadados: %v\n", metadata)
	}

	fmt.Println("\n🏁 --- FIM DO DIAGNÓSTICO ---")
}
