//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	query := `
	query {
	  Page(page: 1, perPage: 10) {
		media(status: RELEASING, sort: TRENDING_DESC, type: ANIME) {
		  title {
			romaji
			english
		  }
		  nextAiringEpisode {
			episode
			airingAt
		  }
		}
	  }
	}`

	requestBody, _ := json.Marshal(map[string]interface{}{
		"query": query,
	})

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Printf("Erro: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Page struct {
				Media []struct {
					Title struct {
						Romaji  string `json:"romaji"`
						English string `json:"english"`
					} `json:"title"`
				} `json:"media"`
			} `json:"page"`
		} `json:"data"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("--- LANÇAMENTOS TRENDING (ALVOS REAIS) ---")
	for _, m := range result.Data.Page.Media {
		name := m.Title.English
		if name == "" {
			name = m.Title.Romaji
		}
		fmt.Printf("- %s\n", name)
	}
}
