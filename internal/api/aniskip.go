package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"os"
	"path/filepath"
	"strings"
)

// GetAniSkipData fetches skip times data for a given anime ID and episode
func GetAniSkipData(animeMalId int, episode int) (string, error) {
	baseURL := "https://api.aniskip.com/v1/skip-times"

	url := fmt.Sprintf("%s/%d/%d?types=op&types=ed", baseURL, animeMalId, episode)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching data from AniSkip API: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AniSkip API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// RoundTime rounds a time value to the specified precision
func RoundTime(timeValue float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Floor(timeValue*multiplier+0.5) / multiplier
}

// ParseAniSkipResponse parses the response text from the AniSkip API and updates the Episode struct
func ParseAniSkipResponse(responseText string, episode *models.Episode, timePrecision int) error {
	if responseText == "" {
		return fmt.Errorf("response text is empty")
	}

	var data models.SkipTimesResponse
	err := json.Unmarshal([]byte(responseText), &data)
	if err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	if util.IsDebug {
		// Log the raw response for debugging
		fmt.Printf("AniSkip Raw Response: %s\n", responseText)
	}

	if !data.Found {
		return fmt.Errorf("no skip times found")
	}

	// Populate skip times for the episode
	for _, result := range data.Results {
		start := int(RoundTime(result.Interval.StartTime, timePrecision))
		end := int(RoundTime(result.Interval.EndTime, timePrecision))

		// Populate based on the type of skip (OP or ED)
		switch result.Type {
		case "op":
			episode.SkipTimes.Op = models.Skip{Start: start, End: end}
		case "ed":
			episode.SkipTimes.Ed = models.Skip{Start: start, End: end}
		default:
			fmt.Printf("Unknown skip type encountered: %s\n", result.Type)
		}
	}

	return nil
}

// GetAndParseAniSkipData fetches and parses skip times for a given anime ID and episode
func GetAndParseAniSkipData(animeMalId int, episodeNum int, episode *models.Episode) error {
	responseText, err := GetAniSkipData(animeMalId, episodeNum)
	if err != nil {
		return err
	}
	return ParseAniSkipResponse(responseText, episode, 0)
}

// WriteAniSkipSidecar writes a JSON file with OP/ED skip times next to the video
func WriteAniSkipSidecar(videoPath string, ep *models.Episode) error {
	if ep == nil {
		return nil
	}
	// Only write if we have at least one skip window
	if ep.SkipTimes.Op.Start == 0 && ep.SkipTimes.Op.End == 0 && ep.SkipTimes.Ed.Start == 0 && ep.SkipTimes.Ed.End == 0 {
		return nil
	}

	type skipFile struct {
		Format  string `json:"format"`
		OPStart int    `json:"op_start"`
		OPEnd   int    `json:"op_end"`
		EDStart int    `json:"ed_start"`
		EDEnd   int    `json:"ed_end"`
		Updated string `json:"updated"`
		Episode string `json:"episode"`
		Source  string `json:"source"`
	}

	payload := skipFile{
		Format:  "aniskip",
		OPStart: ep.SkipTimes.Op.Start,
		OPEnd:   ep.SkipTimes.Op.End,
		EDStart: ep.SkipTimes.Ed.Start,
		EDEnd:   ep.SkipTimes.Ed.End,
		Updated: time.Now().Format(time.RFC3339),
		Episode: ep.Number,
		Source:  "GoAnime",
	}

	b, _ := json.MarshalIndent(payload, "", "  ")
	sidecar := strings.TrimSuffix(videoPath, filepath.Ext(videoPath)) + ".skips.json"
	// Restrictive permissions: owner read/write only
	return os.WriteFile(sidecar, b, 0600)
}
