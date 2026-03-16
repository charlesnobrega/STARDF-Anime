package watchlist

import (
	"fmt"
	"time"

	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

// WatchlistManager handles logic for followed media and notifications
type WatchlistManager struct {
	tracker *tracking.LocalTracker
}

func NewWatchlistManager(tracker *tracking.LocalTracker) *WatchlistManager {
	return &WatchlistManager{tracker: tracker}
}

// FollowMedia adds a media to the local watchlist
func (wm *WatchlistManager) FollowMedia(media *models.Media, status string) error {
	if wm.tracker == nil {
		return fmt.Errorf("tracker not initialized")
	}

	totalEp := media.TotalEpisodes
	if totalEp == 0 && len(media.Episodes) > 0 {
		totalEp = len(media.Episodes)
	}

	followed := tracking.FollowedMedia{
		AllanimeID:    media.URL, // Using URL as the unique identifier for now
		AnilistID:     media.AnilistID,
		Title:         media.Name,
		MediaType:     string(media.MediaType),
		TotalEpisodes: totalEp,
		Status:        status,
		LastChecked:   time.Now(),
		LastEpisode:   totalEp,
		UpdatedAt:     time.Now(),
	}

	err := wm.tracker.UpdateFollowedMedia(followed)
	if err == nil {
		util.Info("Conteúdo adicionado à sua lista!", "title", media.Name)
	}
	return err
}

// GetWatchlist returns all followed media
func (wm *WatchlistManager) GetWatchlist() ([]tracking.FollowedMedia, error) {
	if wm.tracker == nil {
		return nil, fmt.Errorf("tracker not initialized")
	}
	return wm.tracker.GetAllFollowedMedia()
}

// UnfollowMedia removes a media from the watchlist
func (wm *WatchlistManager) UnfollowMedia(allanimeID string) error {
	if wm.tracker == nil {
		return fmt.Errorf("tracker not initialized")
	}
	return wm.tracker.UnfollowMedia(allanimeID)
}

// CheckForUpdates checks for new episodes for all items in the watchlist
// This is a placeholder for actual scraper integration
func (wm *WatchlistManager) CheckForUpdates() ([]string, error) {
	list, err := wm.GetWatchlist()
	if err != nil {
		return nil, err
	}

	var notifications []string
	for _, item := range list {
		// In a real implementation, we would call a scraper here to check current episode count
		// For now, we'll just log that we are checking
		util.Debug("Checking updates for", "title", item.Title)
		
		// If item.Status == "watching" and new episodes are found...
		// notifications = append(notifications, fmt.Sprintf("Novo episódio de %s disponível!", item.Title))
	}

	return notifications, nil
}
