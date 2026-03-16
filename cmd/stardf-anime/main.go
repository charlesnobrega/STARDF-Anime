package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/charlesnobrega/STARDF-Anime/internal/handlers"
	"github.com/charlesnobrega/STARDF-Anime/internal/player"
	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
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
				log.Fatalln(util.ErrorHandler(updateErr))
			}
			return
		}
		if errors.Is(err, util.ErrDownloadRequested) {
			if downloadErr := handlers.HandleDownloadRequest(); downloadErr != nil {
				log.Fatalln(util.ErrorHandler(downloadErr))
			}
			return
		}
		if errors.Is(err, util.ErrAniListLoginRequested) {
			if loginErr := handlers.HandleAniListLogin(); loginErr != nil {
				log.Fatalln(util.ErrorHandler(loginErr))
			}
			return
		}
		if errors.Is(err, util.ErrAniListLogoutRequested) {
			if logoutErr := handlers.HandleAniListLogout(); logoutErr != nil {
				log.Fatalln(util.ErrorHandler(logoutErr))
			}
			return
		}
		if errors.Is(err, util.ErrHelpRequested) {
			return
		}
		log.Fatalln(util.ErrorHandler(err))
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
				log.Fatalln(util.ErrorHandler(promptErr))
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
			log.Fatalln(util.ErrorHandler(err))
		}
		
		// If we reached here without error, we might want to return to menu or exit
		// For now, let's return to menu after playback/action finishes
		menuResult = util.MenuResult{}
	}
}
