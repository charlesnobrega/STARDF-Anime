package main

import (
	"fmt"
	"github.com/alvarorichard/Goanime/internal/player"
	"github.com/alvarorichard/Goanime/internal/scraper"
	"github.com/alvarorichard/Goanime/internal/util"
)

func main() {
	util.IsDebug = true
	util.InitLogger()

	fmt.Println("=== Real Player Test: Cineby ===")
	client := scraper.NewCinebyClient()
	
	// 1. Search for a medium (Interstellar)
	fmt.Println("Searching for Interstellar...")
	results, err := client.SearchMedia("Interstellar")
	if err != nil || len(results) == 0 {
		fmt.Printf("Search failed: %v (results: %d)\n", err, len(results))
		return
	}
	
	selected := results[0]
	fmt.Printf("Selected: %s - %s\n", selected.Name, selected.URL)
	
	// 2. Get episodes (should be 1 for movie)
	fmt.Println("Fetching episodes...")
	episodes, err := client.GetEpisodes(selected.URL)
	if err != nil || len(episodes) == 0 {
		fmt.Printf("Fetching episodes failed: %v\n", err)
		return
	}
	
	// 3. Get Stream URLs
	fmt.Println("Extracting stream URLs...")
	streamURLs, err := client.GetStreamURLs(episodes[0].URL)
	if err != nil || len(streamURLs) == 0 {
		fmt.Printf("Extraction failed: %v\n", err)
		return
	}
	
	videoURL := streamURLs[0]
	fmt.Printf("Video URL extracted: %s\n", videoURL)
	
	// 4. Attempt Playback
	fmt.Println("Launching Player (MPV should open)...")
	// Using a simplified call to playVideo if possible, or HandleDownloadAndPlay
	// Note: HandleDownloadAndPlay handles the TUI too, which might fail in non-interactive terminal,
	// but we want to see if it triggers the MPV launch command.
	
	err = player.HandleDownloadAndPlay(
		videoURL,
		episodes,
		1, // Episode 1
		selected.URL,
		"1", // Episode Number
		0, // MalID (not needed for movie)
		nil, // discord updater
	)
	
	if err != nil {
		fmt.Printf("Playback ended/failed: %v\n", err)
	} else {
		fmt.Println("Playback finished successfully.")
	}
}
