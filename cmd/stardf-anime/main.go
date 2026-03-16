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
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
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

	// Initialize tracker early in background to avoid delays when playing movies
	player.InitTrackerAsync()

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
		if errors.Is(err, util.ErrHelpRequested) {
			return
		}
		log.Fatalln(util.ErrorHandler(err))
	}

	targetName := cliName

	// Loop for main menu navigation
	for {
		// If no name from CLI or returning from menu, prompt user
		if targetName == "" {
			var promptErr error
			targetName, promptErr = util.PromptInteractive()
			if promptErr != nil {
				// Handle exit/cancel from prompt
				if errors.Is(promptErr, context.Canceled) || strings.Contains(promptErr.Error(), "user") {
					return
				}
				log.Fatalln(util.ErrorHandler(promptErr))
			}
		}

		// Handle normal playback mode
		err = handlers.HandlePlaybackMode(targetName)
		if err != nil {
			if errors.Is(err, util.ErrBackToMainMenu) {
				util.Infof("Returning to main menu...")
				targetName = "" // Reset to force interactive prompt
				continue
			}
			log.Fatalln(util.ErrorHandler(err))
		}
		break
	}
}
