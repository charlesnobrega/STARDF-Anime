package handlers

import (
	"fmt"

	"github.com/charlesnobrega/STARDF-Anime/internal/updater"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

// HandleUpdateRequest processes update requests
func HandleUpdateRequest() error {
	// Initialize logger for update process
	util.InitLogger()
	util.Info("Checking for updates...")
	if updateErr := updater.CheckAndPromptUpdate(); updateErr != nil {
		return fmt.Errorf("update failed: %w", updateErr)
	}
	return nil
}
