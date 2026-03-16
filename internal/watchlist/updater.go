package watchlist

import (
	"fmt"
	"time"

	"github.com/charlesnobrega/STARDF-Anime/internal/api"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/notify"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

// SyncWatchlist checks for updates in the background
func (wm *WatchlistManager) SyncWatchlist() error {
	list, err := wm.tracker.GetAllFollowedMedia()
	if err != nil {
		return err
	}

	util.Infof("Sincronizando sua lista (Watchlist)...")
	
	count := 0
	for _, item := range list {
		// Only check if status is "watching" or "planned"
		if item.Status != "watching" && item.Status != "planned" {
			continue
		}

		util.Debug("Checking updates for", "title", item.Title)
		
		// Create a temporary media object for API call
		m := &models.Anime{
			URL:    item.AllanimeID,
			Source: item.MediaType, // This is a bit hacky, but GetAnimeEpisodesEnhanced uses item.Source
		}
		
		// Determine source from URL if possible
		if m.Source == "" {
			// Try to guess source from URL
		}

		episodes, err := api.GetAnimeEpisodesEnhanced(m)
		if err != nil {
			util.Debug("Failed to check updates", "title", item.Title, "error", err)
			continue
		}

		if len(episodes) > item.LastEpisode {
			diff := len(episodes) - item.LastEpisode
			util.Infof("✨ NOVO: %d novos episódios encontrados para %s!", diff, item.Title)

			// Fire a desktop notification
			notify.Send(
				"StarDF-Anime — Novos Episódios! 🎉",
				fmt.Sprintf("%d novo(s) ep(s) disponível(is) de %s!", diff, item.Title),
			)

			// Update the record
			item.TotalEpisodes = len(episodes)
			item.LastEpisode = len(episodes)
			item.UpdatedAt = time.Now()
			_ = wm.tracker.UpdateFollowedMedia(item)
			count++
		}
	}

	if count > 0 {
		util.Infof("Sincronização concluída. %d obras com atualizações!", count)
	} else {
		util.Debug("Nenhuma atualização encontrada na lista.")
	}

	return nil
}

// BackgroundSync starts a ticker to check for updates periodically
func (wm *WatchlistManager) BackgroundSync(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			_ = wm.SyncWatchlist()
		}
	}()
}
