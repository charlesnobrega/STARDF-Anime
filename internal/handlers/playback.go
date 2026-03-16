package handlers

import (
	"errors"

	"github.com/charlesnobrega/STARDF-Anime/internal/appflow"
	"github.com/charlesnobrega/STARDF-Anime/internal/discord"
	"github.com/charlesnobrega/STARDF-Anime/internal/playback"
	"github.com/charlesnobrega/STARDF-Anime/internal/player"
	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charlesnobrega/STARDF-Anime/internal/version"
)

// HandlePlaybackMode processes normal anime playback
func HandlePlaybackMode(animeName string) error {
	timer := util.StartTimer("PlaybackMode:Total")
	defer timer.Stop()

	// Initialize the beautiful logger
	util.InitLogger()

	tracking.HandleTrackingNotice()
	util.Debugf("[PERF] starting StarDF-Anime v%s", version.Version)

	discordTimer := util.StartTimer("Discord:Initialize")
	discordManager := discord.NewManager()
	if err := discordManager.Initialize(); err != nil {
		util.Debug("Failed to initialize Discord Rich Presence:", "error", err)
	} else {
		defer discordManager.Shutdown()
	}
	discordTimer.Stop()

	currentAnimeName := animeName

	for {
		// Use enhanced search with retry logic
		searchTimer := util.StartTimer("SearchAnime:WithRetry")
		anime, err := appflow.SearchAnimeWithRetry(currentAnimeName)
		searchTimer.Stop()

		if err != nil {
			if errors.Is(err, util.ErrBackToMainMenu) {
				return err
			}
			util.Errorf("Failed to search for anime: %v", err)
			return err
		}

		detailsTimer := util.StartTimer("FetchAnimeDetails")
		appflow.FetchAnimeDetails(anime)
		detailsTimer.Stop()

		episodesTimer := util.StartTimer("GetAnimeEpisodes")
		episodes, err := appflow.GetAnimeEpisodes(anime)
		episodesTimer.Stop()

		if err != nil {
			util.Warnf("Could not load episodes from %s: %v", anime.Source, err)
			util.Infof("Please pick another source from the list.")
			// Setting anime to nil will trigger search again with same term (effectively going back)
			continue
		}

		util.PerfCount("anime_loaded")

		series, totalEpisodes := playback.CheckIfSeriesEnhanced(anime)
		var playbackErr error

		playbackTimer := util.StartTimer("Playback:Handle")
		if series {
			playbackErr = playback.HandleSeries(anime, episodes, totalEpisodes, discordManager.IsEnabled())
		} else {
			playbackErr = playback.HandleMovie(anime, episodes, discordManager.IsEnabled())
		}
		playbackTimer.Stop()

		// Check if user wants to go back to anime selection
		if errors.Is(playbackErr, player.ErrBackToAnimeSelection) {
			util.Infof("Going back to anime selection...")
			// Keep the same search term to show the anime list again
			continue
		}

		// Normal exit or other errors
		if playbackErr != nil && !errors.Is(playbackErr, player.ErrBackToAnimeSelection) {
			return playbackErr
		}
		break
	}
	return nil
}
