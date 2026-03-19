package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
)

// SyncMetadata saves anime/movie/series metadata to the download folder
func SyncMetadata(outputDir string, anime *models.Anime) error {
	if anime == nil {
		return nil
	}

	metadataPath := filepath.Join(outputDir, "metadata.json")
	
	// Create or truncate the file
	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(anime); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	// Also create a basic .nfo file for media players (Kodi/Plex style)
	nfoPath := filepath.Join(outputDir, "tvshow.nfo")
	if GlobalMediaType == "movie" {
		nfoPath = filepath.Join(outputDir, "movie.nfo")
	}

	id := anime.AnilistID
	if id == 0 {
		id = anime.TMDBID
	}

	nfoContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<movie>
    <title>%s</title>
    <plot>%s</plot>
    <thumb>%s</thumb>
    <id>%d</id>
</movie>`, anime.Name, "Downloaded via StarDF-Anime", anime.ImageURL, id)

	_ = os.WriteFile(nfoPath, []byte(nfoContent), 0644)

	Infof("Metadados sincronizados em: %s", outputDir)
	return nil
}
