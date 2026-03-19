package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/charlesnobrega/STARDF-Anime/internal/anilist"
	"github.com/charlesnobrega/STARDF-Anime/internal/handlers"
	"github.com/charlesnobrega/STARDF-Anime/internal/player"
	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charlesnobrega/STARDF-Anime/internal/webui"
	"github.com/charlesnobrega/STARDF-Anime/internal/watchlist"
	"time"
)

func main() {
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		util.RunCleanup()
		os.Exit(0)
	}()

	// Ensure cleanup runs on normal exit
	defer util.RunCleanup()

	// Start total execution timer
	timer := util.StartTimer("TotalExecution")
	defer timer.Stop()

	// Initialize tracker early in background
	player.InitTrackerAsync()

	// Sync watchlist in background
	go func() {
		tracker := tracking.GetGlobalTracker()
		// Wait a bit for tracker to be ready
		for i := 0; i < 10 && tracker == nil; i++ {
			time.Sleep(500 * time.Millisecond)
			tracker = tracking.GetGlobalTracker()
		}
		if tracker != nil {
			wm := watchlist.NewWatchlistManager(tracker)
			_ = wm.SyncWatchlist()
		}
	}()

	// Parse initial flags/CLI arguments
	cliName, err := util.ParseFlags()
	if err != nil {
		// Handle special requests
		if errors.Is(err, util.ErrUpdateRequested) {
			if updateErr := handlers.HandleUpdateRequest(); updateErr != nil {
				fmt.Println(util.ErrorHandler(updateErr))
			}
			return
		}
		if errors.Is(err, util.ErrDownloadRequested) {
			if downloadErr := handlers.HandleDownloadRequest(); downloadErr != nil {
				fmt.Println(util.ErrorHandler(downloadErr))
			}
			return
		}
		if errors.Is(err, util.ErrAniListLoginRequested) {
			if loginErr := handlers.HandleAniListLogin(); loginErr != nil {
				fmt.Println(util.ErrorHandler(loginErr))
			}
			return
		}
		if errors.Is(err, util.ErrAniListLogoutRequested) {
			if logoutErr := handlers.HandleAniListLogout(); logoutErr != nil {
				fmt.Println(util.ErrorHandler(logoutErr))
			}
			return
		}
		if errors.Is(err, util.ErrWebRequested) {
			port := 8080
			portStr := os.Getenv("STARDF_PORT")
			if portStr != "" {
				fmt.Sscanf(portStr, "%d", &port)
			}
			if webErr := webui.StartWebUI(port); webErr != nil {
				fmt.Println(util.ErrorHandler(webErr))
			}
			return
		}
		if errors.Is(err, util.ErrHelpRequested) {
			return
		}
		fmt.Println(util.ErrorHandler(err))
		return
	}

	// Setup AniList user display
	util.GetMenuSubtitleFunc = func() string {
		if anilist.GlobalSession.IsLoggedIn() && anilist.GlobalSession.CurrentUser != nil {
			return " [👤 " + anilist.GlobalSession.CurrentUser.Name + "]"
		}
		return ""
	}

	var menuResult util.MenuResult
	
	// If name from CLI, set initial action to search
	if cliName != "" {
		menuResult = util.MenuResult{Action: util.ActionSearch, SearchTerm: cliName}
	}

	// Loop for main menu navigation
	for {
		// If no result yet or returning from menu, prompt user
		if menuResult.Action == 0 && menuResult.SearchTerm == "" {
			var promptErr error
			menuResult, promptErr = util.PromptInteractive()
			if promptErr != nil {
				// Handle exit/cancel from prompt
				if errors.Is(promptErr, context.Canceled) || errors.Is(promptErr, util.ErrExitRequested) || strings.Contains(promptErr.Error(), "user") {
					return
				}
				fmt.Println(util.ErrorHandler(promptErr))
				return
			}
		}

		// Route based on action
		var err error
		switch menuResult.Action {
		case util.ActionSearch:
			err = handlers.HandlePlaybackMode(menuResult.SearchTerm)
		case util.ActionWatchlist:
			err = handlers.HandleWatchlistMode()
		case util.ActionContinue:
			err = handlers.HandleContinueWatchingMode()
		case util.ActionHealth:
			err = handlers.HandleScraperHealthMode()
		case util.ActionTheme:
			err = handlers.HandleThemeMode()
		case util.ActionExit:
			return
		}

		if err != nil {
			if errors.Is(err, util.ErrBackToMainMenu) {
				util.Infof("Returning to main menu...")
				menuResult = util.MenuResult{} // Reset to force interactive prompt
				continue
			}
			fmt.Println(util.ErrorHandler(err))
			return
		}
		
		// If we reached here without error, we might want to return to menu or exit
		// For now, let's return to menu after playback/action finishes
		menuResult = util.MenuResult{}
	}
}
