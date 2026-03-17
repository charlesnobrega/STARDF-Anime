package mobile

import (
	"encoding/json"

	"github.com/charlesnobrega/STARDF-Anime/pkg/stardf"
	"github.com/charlesnobrega/STARDF-Anime/pkg/stardf/types"
)

var client *stardf.Client

func init() {
	client = stardf.NewClient()
}

// Search (query: string) -> JSON
func Search(query string) string {
	results, err := client.SearchAnime(query, nil)
	if err != nil {
		return errorToJSON(err)
	}
	return dataToJSON(results)
}

// GetEpisodes (animeURL: string, sourceName: string) -> JSON
func GetEpisodes(animeURL string, sourceName string) string {
	source, err := types.ParseSource(sourceName)
	if err != nil {
		return errorToJSON(err)
	}
	episodes, err := client.GetAnimeEpisodes(animeURL, source)
	if err != nil {
		return errorToJSON(err)
	}
	return dataToJSON(episodes)
}

// GetStream (animeJSON: string, episodeJSON: string) -> JSON
func GetStream(animeJSON string, episodeJSON string) string {
	var anime types.Anime
	if err := json.Unmarshal([]byte(animeJSON), &anime); err != nil {
		return errorToJSON(err)
	}

	var episode types.Episode
	if err := json.Unmarshal([]byte(episodeJSON), &episode); err != nil {
		return errorToJSON(err)
	}

	url, metadata, err := client.GetEpisodeStreamURL(&anime, &episode, nil)
	if err != nil {
		return errorToJSON(err)
	}

	response := map[string]interface{}{
		"url":      url,
		"metadata": metadata,
	}
	return dataToJSON(response)
}

func errorToJSON(err error) string {
	res := map[string]interface{}{
		"error": err.Error(),
	}
	b, _ := json.Marshal(res)
	return string(b)
}

func dataToJSON(data interface{}) string {
	res := map[string]interface{}{
		"data": data,
	}
	b, _ := json.Marshal(res)
	return string(b)
}
