package watchlist

import (
	"fmt"
	"time"

	"github.com/charlesnobrega/STARDF-Anime/internal/anilist"
	"github.com/charlesnobrega/STARDF-Anime/internal/api"
	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/notify"
	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
)

// SyncWatchlist checks for updates in the background
// SyncWatchlist checks for updates in the background
func (wm *WatchlistManager) SyncWatchlist() error {
	// Pull from AniList first if logged in
	if err := wm.PullFromAniList(); err != nil {
		util.Warnf("Erro ao sincronizar com AniList: %v", err)
	}

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

// PullFromAniList fetches the user's AniList watchlist and merges it into local tracking
func (wm *WatchlistManager) PullFromAniList() error {
	session := anilist.GlobalSession
	if !session.IsLoggedIn() {
		return nil // Not logged in, skip
	}

	user := session.CurrentUser
	if user == nil {
		var err error
		user, err = session.Client.GetViewer()
		if err != nil {
			return err
		}
		session.CurrentUser = user
	}

	util.Infof("Sincronizando com o AniList...")

	list, err := session.Client.GetUserList(user.ID)
	if err != nil {
		return fmt.Errorf("falha ao buscar lista do AniList: %w", err)
	}

	count := 0
	localMediaList, _ := wm.tracker.GetAllFollowedMedia()

	for _, l := range list.Data.MediaListCollection.Lists {
		for _, entry := range l.Entries {
			status := "planning"
			switch entry.Status {
			case "CURRENT":
				status = "watching"
			case "COMPLETED":
				status = "completed"
			case "DROPPED":
				status = "dropped"
			case "PAUSED":
				status = "paused"
			case "PLANNING":
				status = "planning"
			}

			title := entry.Media.Title.Romaji
			if title == "" {
				title = entry.Media.Title.English
			}

			var existing *tracking.FollowedMedia
			for i, m := range localMediaList {
				if m.AnilistID == entry.MediaID || m.Title == title {
					existing = &localMediaList[i]
					break
				}
			}

			var allanimeID string
			if existing != nil {
				allanimeID = existing.AllanimeID
			} else {
				allanimeID = fmt.Sprintf("anilist_%d", entry.MediaID)
			}

			totalEp := 0
			if entry.Media.Episodes != nil {
				totalEp = *entry.Media.Episodes
			}

			updatedAt := time.Unix(int64(entry.UpdatedAt), 0)

			// Skip older updates using AniList timestamp (though AniList might just sync once)
			if existing != nil && existing.UpdatedAt.After(updatedAt) {
				continue
			}

			newFollow := tracking.FollowedMedia{
				AllanimeID:    allanimeID,
				AnilistID:     entry.MediaID,
				Title:         title,
				MediaType:     "anime",
				TotalEpisodes: totalEp,
				Status:        status,
				LastChecked:   time.Now(),
				LastEpisode:   entry.Progress,
				UpdatedAt:     updatedAt,
			}

			err := wm.tracker.UpdateFollowedMedia(newFollow)
			if err == nil {
				count++
			}
		}
	}

	if count > 0 {
		util.Infof("Sincronização com AniList concluída. %d itens atualizados!", count)
	}
	return nil
}
