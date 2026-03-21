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
	query ($search: String) {
	  Page(page: 1, perPage: 1) {
		media(search: $search, type: ANIME) {
		  id
		  title {
			romaji
			english
		  }
		}
	  }
	}`

	variables := map[string]interface{}{
		"search": "Naruto",
	}

	requestBody, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Printf("[❌] AniList Connection Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("[❌] AniList API returned Status %d\n", resp.StatusCode)
		return
	}

	fmt.Printf("[✅] AniList API is functional. Status: %d\n", resp.StatusCode)
}
