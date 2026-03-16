package playback

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charlesnobrega/STARDF-Anime/internal/watchlist"
	"github.com/charmbracelet/huh"
)

func GetUserInput() string {
	var choice string

	menu := huh.NewSelect[string]().
		Title("Playback Control").
		Description("What would you like to do next?").
		Options(
			huh.NewOption("Next episode", "n"),
			huh.NewOption("Previous episode", "p"),
			huh.NewOption("Select episode", "e"),
			huh.NewOption("⭐ Rate anime", "rate"),
			huh.NewOption("Change anime", "c"),
			huh.NewOption("← Back", "back"),
			huh.NewOption("Quit", "q"),
		).
		Value(&choice)

	if err := menu.Run(); err != nil {
		util.Errorf("Error showing menu: %v", err)
		return "n" // Default to next episode on error
	}

	return choice
}

// HandleRating prompts the user to rate the current anime and saves the score
func HandleRating(anime *models.Anime) {
	var input string
	err := huh.NewInput().
		Title("Rate Anime").
		Description(fmt.Sprintf("Enter a rating for %s (1-10):", anime.Name)).
		Value(&input).
		Run()
	if err != nil {
		util.Errorf("Error getting rating: %v", err)
		return
	}
	score, err := strconv.Atoi(input)
	if err != nil || score < 1 || score > 10 {
		util.Errorf("Invalid rating. Must be between 1 and 10.")
		return
	}
	tracker := tracking.GetGlobalTracker()
	if tracker != nil {
		media, err := tracker.GetFollowedMedia(anime.URL)
		if err == nil && media != nil {
			media.Score = score
			media.UpdatedAt = time.Now()
			tracker.UpdateFollowedMedia(*media)
			util.Infof("Rated %s: %d/10", anime.Name, score)
		} else {
			wm := watchlist.NewWatchlistManager(tracker)
			m := &models.Media{
				URL:           anime.URL,
				Name:          anime.Name,
				AnilistID:     anime.AnilistID,
				MediaType:     anime.MediaType,
				TotalEpisodes: 0,
				ImageURL:      anime.ImageURL,
				Episodes:      nil,
			}
			_ = wm.FollowMedia(m, "watching")
			// Fetch again to update score
			media, _ = tracker.GetFollowedMedia(anime.URL)
			if media != nil {
				media.Score = score
				media.UpdatedAt = time.Now()
				tracker.UpdateFollowedMedia(*media)
			}
			util.Infof("Rated %s: %d/10 (Added to watchlist)", anime.Name, score)
		}
	} else {
		util.Warnf("Tracker not initialized, unable to save rating.")
	}
}
