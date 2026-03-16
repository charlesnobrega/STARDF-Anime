package playback

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/charlesnobrega/STARDF-Anime/internal/api"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/player"
	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charlesnobrega/STARDF-Anime/internal/watchlist"
	"github.com/charmbracelet/huh"
)

func HandleSeries(anime *models.Anime, episodes []models.Episode, totalEpisodes int, discordEnabled bool) error {
	fmt.Printf("The selected anime is a series with %d episodes.\n", totalEpisodes)
	animeMutex := sync.Mutex{}
	isPaused := false

	selectedEpisodeURL, episodeNumberStr, selectedEpisodeNum, err := SelectInitialEpisode(anime, episodes)
	if err != nil {
		// If user selected back at initial episode selection, return to anime selection
		if errors.Is(err, player.ErrBackRequested) {
			return player.ErrBackToAnimeSelection
		}
		log.Printf("Episode selection error: %v", util.ErrorHandler(err))
		return err
	}

	for {
		err := PlayEpisode(
			anime,
			episodes,
			selectedEpisodeNum,
			selectedEpisodeURL,
			episodeNumberStr,
			discordEnabled,
			&isPaused,
			&animeMutex,
		)

		// Check if user quit during video playback
		if errors.Is(err, player.ErrUserQuit) {
			log.Println("Quitting application as per user request.")
			break
		}

		// Check if user requested to go back to episode selection (from server selection)
		if errors.Is(err, player.ErrBackToEpisodeSelection) {
			selectedEpisodeURL, episodeNumberStr, selectedEpisodeNum, err = SelectInitialEpisode(anime, episodes)
			if err != nil {
				// If user selected back at episode selection, go back to anime selection
				if errors.Is(err, player.ErrBackRequested) {
					return player.ErrBackToAnimeSelection
				}
				log.Printf("Error selecting episode: %v", err)
			}
			continue
		}

		// Check if user requested to change anime during video playback
		if errors.Is(err, player.ErrChangeAnime) {
			newAnime, newEpisodes, err := ChangeAnimeLocal()
			if err != nil {
				log.Printf("Error changing anime: %v", err)
				continue // Stay with current anime if change fails
			}

			// Update anime and episodes
			anime = newAnime
			episodes = newEpisodes

			// Check if new anime is a series and get new total episodes
			series, newTotalEpisodes := CheckIfSeriesEnhanced(anime)
			totalEpisodes = newTotalEpisodes

			if !series {
				// If new anime is a movie, handle it differently
				log.Println("Switched to a movie/OVA, handling as single episode.")
				if err := HandleMovie(anime, episodes, discordEnabled); err != nil {
					if errors.Is(err, player.ErrBackToAnimeSelection) {
						return err
					}
				}
				break
			}

			// Select initial episode for the new anime
			selectedEpisodeURL, episodeNumberStr, selectedEpisodeNum, err = SelectInitialEpisode(anime, episodes)
			if err != nil {
				log.Printf("Error selecting episode for new anime: %v", err)
				continue
			}

			fmt.Printf("Switched to anime: %s with %d episodes.\n", anime.Name, totalEpisodes)
			continue // Skip normal navigation and start playing the new anime
		}

		// Handle other errors
		if err != nil {
			log.Printf("Error during episode playback: %v", err)
		}

		userInput := GetUserInput()
		if userInput == "q" || userInput == "quit" {
			log.Println("Quitting application as per user request.")
			break
		}

		if userInput == "rate" {
			HandleRating(anime)
			continue
		}

		// Handle back/change anime - both options allow searching for a new anime
		if userInput == "c" || userInput == "back" {
			newAnime, newEpisodes, err := ChangeAnimeLocal()
			if err != nil {
				log.Printf("Error changing anime: %v", err)
				continue // Stay with current anime if change fails
			}

			// Update anime and episodes
			anime = newAnime
			episodes = newEpisodes

			// Check if new anime is a series and get new total episodes
			series, newTotalEpisodes := CheckIfSeriesEnhanced(anime)
			totalEpisodes = newTotalEpisodes

			if !series {
				// If new anime is a movie, handle it differently
				log.Println("Switched to a movie/OVA, handling as single episode.")
				if err := HandleMovie(anime, episodes, discordEnabled); err != nil {
					if errors.Is(err, player.ErrBackToAnimeSelection) {
						return err
					}
				}
				break
			}

			// Select initial episode for the new anime
			selectedEpisodeURL, episodeNumberStr, selectedEpisodeNum, err = SelectInitialEpisode(anime, episodes)
			if err != nil {
				log.Printf("Error selecting episode for new anime: %v", err)
				continue
			}

			fmt.Printf("Switched to anime: %s with %d episodes.\n", anime.Name, totalEpisodes)
			continue // Skip normal navigation and start playing the new anime
		}

		// Handle episode selection
		if userInput == "e" {
			selectedEpisodeURL, episodeNumberStr, selectedEpisodeNum, err = SelectInitialEpisode(anime, episodes)
			if err != nil {
				// If user selected back, just continue without changing episode
				if errors.Is(err, player.ErrBackRequested) {
					continue
				}
				log.Printf("Error selecting episode: %v", err)
				continue
			}
			continue
		}

		selectedEpisodeURL, episodeNumberStr, selectedEpisodeNum = handleUserNavigation(
			userInput,
			episodes,
			selectedEpisodeNum,
			totalEpisodes,
		)
	}
	return nil
}

func SelectInitialEpisode(anime *models.Anime, episodes []models.Episode) (string, string, int, error) {
	for {
		url, numStr, err := player.SelectEpisodeWithFuzzyFinder(episodes)
		if err != nil {
			if errors.Is(err, player.ErrFollowRequested) {
				// Handle follow request
				tracker := tracking.GetGlobalTracker()
				if tracker != nil {
					wm := watchlist.NewWatchlistManager(tracker)
					media := &models.Media{
						URL:           anime.URL,
						Name:          anime.Name,
						AnilistID:     anime.AnilistID,
						MediaType:     anime.MediaType,
						TotalEpisodes: len(episodes),
						ImageURL:      anime.ImageURL,
						Episodes:      nil, // Not needed for followed_media table basic entry
					}
					_ = wm.FollowMedia(media, "watching")
				}
				// After following, stay in this loop to let user select an episode
				continue
			}
			// Propagate other errors
			return "", "", -1, err
		}
		selectedEpisodeNum, err := strconv.Atoi(player.ExtractEpisodeNumber(numStr))
		if err != nil {
			// If it's a movie/single episode, ExtractEpisodeNumber returns "1"
			return url, numStr, 1, nil
		}
		return url, numStr, selectedEpisodeNum, nil
	}
}

func handleUserNavigation(input string, episodes []models.Episode, currentNum, totalEpisodes int) (string, string, int) {
	switch input {
	case "e":
		return SelectEpisodeWithFuzzy(episodes)
	case "p":
		newNum := currentNum - 1
		if newNum < 1 {
			newNum = 1
		}
		return FindEpisodeByNumber(episodes, newNum)
	default: // 'n' or default
		newNum := currentNum + 1
		if newNum > totalEpisodes {
			newNum = totalEpisodes
		}
		return FindEpisodeByNumber(episodes, newNum)
	}
}


func CheckIfSeries(url string) (bool, int) {
	series, totalEpisodes, err := api.IsSeries(url)
	if err != nil {
		// Instead of killing the app, assume series unknown -> treat as single episode (movie)
		log.Printf("Error checking if the anime is a series: %v", util.ErrorHandler(err))
		return false, 1
	}
	return series, totalEpisodes
}

// CheckIfSeriesEnhanced checks if anime is a series using enhanced API
func CheckIfSeriesEnhanced(anime *models.Anime) (bool, int) {
	series, totalEpisodes, err := api.IsSeriesEnhanced(anime)
	if err != nil {
		log.Printf("Error checking if the anime is a series: %v", util.ErrorHandler(err))
		return false, 1
	}
	return series, totalEpisodes
}

// ChangeAnimeLocal allows the user to search for and select a new anime (local implementation to avoid circular imports)
func ChangeAnimeLocal() (*models.Anime, []models.Episode, error) {
	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
		var animeName string

		prompt := huh.NewInput().
			Title("Change Anime").
			Description("Enter the name of the anime you want to watch:").
			Value(&animeName).
			Validate(func(v string) error {
				if len(v) < 2 {
					return fmt.Errorf("anime name must be at least 2 characters")
				}
				return nil
			})

		if err := prompt.Run(); err != nil {
			return nil, nil, err
		}

		// Use the enhanced API to search for anime
		anime, err := api.SearchAnimeEnhanced(animeName, "")
		if err != nil || anime == nil {
			if i < maxRetries-1 {
				util.Errorf("No anime found with the name: %s", animeName)
				util.Infof("Please try again with a different search term. (Attempt %d/%d)", i+2, maxRetries)
				continue
			}
			return nil, nil, fmt.Errorf("failed to find anime after %d attempts", maxRetries)
		}

		// Get episodes for the new anime using enhanced API
		episodes, err := api.GetAnimeEpisodesEnhanced(anime)
		if err != nil {
			if i < maxRetries-1 {
				util.Errorf("Failed to get episodes for: %s", anime.Name)
				util.Infof("Please try searching for a different anime. (Attempt %d/%d)", i+2, maxRetries)
				continue
			}
			return nil, nil, fmt.Errorf("failed to get episodes after %d attempts", maxRetries)
		}

		return anime, episodes, nil
	}

	return nil, nil, fmt.Errorf("failed to change anime after %d attempts", maxRetries)
}
