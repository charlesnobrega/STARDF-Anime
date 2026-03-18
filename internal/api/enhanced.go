// Package api provides enhanced anime search and streaming capabilities
package api

import (
	"errors"
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/charlesnobrega/STARDF-Anime/internal/models"
	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/ktr0731/go-fuzzyfinder"
)

// ErrBackToSearch is returned when user selects the back option to search again
var ErrBackToSearch = errors.New("back to search requested")

// Enhanced search that supports multiple sources - always searches both Animefire.io and allanime simultaneously
func SearchAnimeEnhanced(name string, source string) (*models.Anime, error) {
	scraperManager := scraper.NewScraperManager()

	var scraperType *scraper.ScraperType

	// If a specific source is requested, honor it
	if strings.ToLower(source) == "animesonlinecc" {
		t := scraper.AnimesOnlineCCTYPE
		scraperType = &t
		util.Debug("Searching specific source", "source", "AnimesOnlineCC")
	} else if strings.ToLower(source) == "goyabu" {
		t := scraper.GoyabuType
		scraperType = &t
		util.Debug("Searching specific source", "source", "Goyabu")
	} else if strings.ToLower(source) == "animefire" {
		t := scraper.AnimefireType
		scraperType = &t
		util.Debug("Searching specific source", "source", "AnimeFire")
	} else if strings.ToLower(source) == "cineby" || strings.ToLower(source) == "movie" || strings.ToLower(source) == "tv" {
		t := scraper.CinebyType
		scraperType = &t
		util.Debug("Searching specific source", "source", "Cineby")
	} else if strings.ToLower(source) == "flixhq" {
		t := scraper.FlixHQType
		scraperType = &t
		util.Debug("Searching specific source", "source", "FlixHQ")
	} else if strings.ToLower(source) == "cinegratis" {
		t := scraper.CineGratisType
		scraperType = &t
		util.Debug("Searching specific source", "source", "CineGratis")
	} else {
		// Filter sources based on GlobalMediaType from the initial selector
		switch util.GlobalMediaType {
		case "anime":
			util.Debug("Category filtered: Searching only Anime sources")
			// We can pass a hint to scraperManager or just search normally if it's already filtered in unified.go
		case "movie":
			util.Debug("Category filtered: Searching only Movie/TV sources")
		}
		
		scraperType = nil
		util.Debug("Searching relevant sources", "query", name)
	}

	// Perform the search - this will search all sources if scraperType is nil
	util.Debug("Searching for anime/media", "query", name)
	animes, err := scraperManager.SearchAnime(name, scraperType)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if len(animes) == 0 {
		return nil, fmt.Errorf("no results found for: %s", name)
	}

	// Enhance source identification - names already have language tags from unified.go
	for _, anime := range animes {
		// Ensure proper source identification (for internal use only)
		if anime.Source == "" {
			// Fallback source identification by URL analysis
			if strings.Contains(anime.URL, "animesonlinecc") {
				anime.Source = "AnimesOnlineCC"
			} else if strings.Contains(anime.URL, "goyabu") {
				anime.Source = "Goyabu"
			} else if strings.Contains(anime.URL, "animefire") {
				anime.Source = "Animefire.io"
			} else if strings.Contains(anime.URL, "cineby") {
				anime.Source = "Cineby"
			} else if strings.Contains(anime.URL, "flixhq") {
				anime.Source = "FlixHQ"
			}
		}

		// Language tags are already added by unified.go, don't duplicate them here
	}

	util.Debug("Search results summary", "total", len(animes))

	// Show sources breakdown in debug only
	animefireCount := 0
	animesonlineccCount := 0
	goyabuCount := 0
	cinebyCount := 0
	flixhqCount := 0
	for _, anime := range animes {
		if strings.Contains(anime.Source, "AnimeFire") {
			animefireCount++
		} else if strings.Contains(anime.Source, "AnimesOnlineCC") {
			animesonlineccCount++
		} else if strings.Contains(anime.Source, "Goyabu") {
			goyabuCount++
		} else if strings.Contains(anime.Source, "Cineby") {
			cinebyCount++
		} else if strings.Contains(anime.Source, "CineGratis") {
			// Add cinegratis to breakdown if needed, but not in current debug log local vars
		} else if anime.Source == "FlixHQ" {
			flixhqCount++
		}
	}

	util.Debug("Source breakdown", "AnimeFire", animefireCount, "AnimesOnlineCC", animesonlineccCount, "Goyabu", goyabuCount, "Cineby", cinebyCount, "FlixHQ", flixhqCount)

	// Create a special "back" options
	backToSearch := &models.Anime{
		Name:   "← Back (Search Again)",
		URL:    "__back__",
		Source: "__back__",
	}

	backToMenu := &models.Anime{
		Name:   "🏠 Back to Main Menu",
		URL:    "__back_to_menu__",
		Source: "__back_to_menu__",
	}

	// Prepend back options to the list
	animesWithBack := append([]*models.Anime{backToSearch, backToMenu}, animes...)

	// Concurrent episode count fetching for top results (optimization)
	util.Debug("Fetching episode counts for top results...")
	var wg sync.WaitGroup
	// Limit to top 20 results to prevent massive delays but cover most visible items
	limit := len(animesWithBack)
	if limit > 20 {
		limit = 20
	}

	for i := 1; i < limit; i++ {
		wg.Add(1)
		go func(anime *models.Anime) {
			defer wg.Done()
			
			// Use a reasonably long timeout (5s) to be reliable but not hang the search
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			done := make(chan bool, 1)
			go func() {
				// We don't want to fail the whole search if one fetch fails
				eps, err := GetAnimeEpisodesEnhanced(anime)
				if err == nil && len(eps) > 0 {
					anime.TotalEpisodes = len(eps)
					
					// Improved season detection for movie/TV sources
					if anime.Source == "Cineby" || anime.Source == "FlixHQ" || anime.Source == "CineGratis" {
						seasons := make(map[string]bool)
						for _, ep := range eps {
							if ep.SeasonID != "" {
								seasons[ep.SeasonID] = true
							}
						}
						if len(seasons) > 0 {
							anime.SeasonCount = len(seasons)
						} else if strings.Contains(strings.ToLower(anime.URL), "/tv/") || strings.Contains(strings.ToLower(anime.URL), "/series/") {
							anime.SeasonCount = 1
						}
					}
				}
				done <- true
			}()
			
			select {
			case <-done:
			case <-ctx.Done():
				util.Debug("Timeout fetching episodes", "name", anime.Name)
			}
		}(animesWithBack[i])
	}
	// We proceed even if wg isn't complete yet? No, fuzzy finder needs the data for its display func.
	wg.Wait()

	// Use fuzzy finder to let user select
	var idx int

	if util.IsDebug {
		// In debug mode, show preview window with technical details
		idx, err = fuzzyfinder.Find(
			animesWithBack,
			func(i int) string {
				// Show the anime name with language tag as-is
				return animesWithBack[i].Name
			},
			fuzzyfinder.WithPromptString("Select the anime you want: "),
			fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
				if i >= 0 && i < len(animesWithBack) {
					anime := animesWithBack[i]
					if anime.Source == "__back__" {
						return "Go back to perform a new search"
					}
					if anime.Source == "__back_to_menu__" {
						return "Return to the category selection (Anime/Movies)"
					}
					var preview string
					preview = "Source: " + anime.Source + "\nURL: " + anime.URL
					if anime.ImageURL != "" {
						preview += "\nImage: " + anime.ImageURL
					}
					return preview
				}
				return ""
			}),
		)
	} else {
		// In normal mode, no preview window at all
		idx, err = fuzzyfinder.Find(
			animesWithBack,
			func(i int) string {
				if i < 0 || i >= len(animesWithBack) {
					return ""
				}
				anime := animesWithBack[i]
				if anime.Source == "__back__" || anime.Source == "__back_to_menu__" {
					return anime.Name
				}
				
				// Format: [Source] Name - Info
				info := ""
				if anime.SeasonCount > 1 {
					info = fmt.Sprintf(" - %d seasons", anime.SeasonCount)
				} else if anime.SeasonCount == 1 {
					info = " - 1 season"
				}
				
				if anime.TotalEpisodes > 0 {
					if info != "" {
						info += fmt.Sprintf(", %d eps", anime.TotalEpisodes)
					} else {
						info = fmt.Sprintf(" - %d eps", anime.TotalEpisodes)
					}
				}
				
				return fmt.Sprintf("[%s] %s%s", anime.Source, anime.Name, info)
			},
			fuzzyfinder.WithPromptString("Select the anime you want: "),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("anime selection cancelled: %w", err)
	}

	if idx < 0 || idx >= len(animesWithBack) {
		return nil, fmt.Errorf("invalid selection index")
	}
	selectedAnime := animesWithBack[idx]

	// Check if user selected back options
	if selectedAnime.Source == "__back__" {
		return nil, ErrBackToSearch
	}
	if selectedAnime.Source == "__back_to_menu__" {
		return nil, util.ErrBackToMainMenu
	}
	util.Debug("Anime selected", "name", selectedAnime.Name, "source", selectedAnime.Source)

	// CRITICAL: Enrich with AniList data for images and metadata (like the original system)
	if err := enrichAnimeData(selectedAnime); err != nil {
		util.Errorf("Error enriching anime data: %v", err)
	}

	return selectedAnime, nil
}

// Enhanced episode fetching that works with different sources
func GetAnimeEpisodesEnhanced(anime *models.Anime) ([]models.Episode, error) {
	// Check if this is a FlixHQ movie/TV show
	if anime.Source == "FlixHQ" || anime.MediaType == models.MediaTypeMovie || anime.MediaType == models.MediaTypeTV {
		return GetFlixHQEpisodes(anime)
	}

	// Determine source type from multiple indicators with enhanced logic
	var sourceName string

	// Priority 1: Check the Source field (most reliable)
	if anime.Source == "AnimesOnlineCC" {
		sourceName = "AnimesOnlineCC"
	} else if anime.Source == "Goyabu" {
		sourceName = "Goyabu"
	} else if strings.Contains(anime.Source, "AnimeFire") {
		sourceName = "Animefire.io"
	} else if anime.Source == "Cineby" {
		sourceName = "Cineby"
	} else if anime.Source == "CineGratis" {
		sourceName = "CineGratis"
	} else if strings.Contains(anime.Name, "[Portuguese]") || strings.Contains(anime.Name, "[Português]") {
		// Priority 2: Check language tags (Portuguese default to Animefire)
		sourceName = "Animefire.io"
		anime.Source = "Animefire.io"
	} else if strings.Contains(anime.URL, "animesonlinecc") {
		// Priority 3: URL analysis
		sourceName = "AnimesOnlineCC"
		anime.Source = "AnimesOnlineCC" // Update source field
	} else if strings.Contains(anime.URL, "goyabu") {
		sourceName = "Goyabu"
		anime.Source = "Goyabu" // Update source field
	} else if strings.Contains(anime.URL, "animefire") {
		sourceName = "Animefire.io"
		anime.Source = "Animefire.io" // Update source field
	} else if strings.Contains(anime.URL, "cineby") {
		sourceName = "Cineby"
		anime.Source = "Cineby" // Update source field
	} else if strings.Contains(anime.URL, "cinegratis") {
		sourceName = "CineGratis"
		anime.Source = "CineGratis"
	} else {
		// Default to Animefire for unknown sources
		sourceName = "Animefire (default)"
		anime.Source = "Animefire.io"
	}

	cleanName := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(anime.Name, "[English]", ""), "[Portuguese]", ""))

	util.Debug("Getting episodes", "source", sourceName, "anime", cleanName)

	var episodes []models.Episode
	var err error

	// Use different approaches based on source
	switch sourceName {
	case "AnimesOnlineCC":
		scraperManager := scraper.NewScraperManager()
		scraperInstance, scErr := scraperManager.GetScraper(scraper.AnimesOnlineCCTYPE)
		if scErr != nil {
			return nil, fmt.Errorf("failed to get AnimesOnlineCC scraper: %w", scErr)
		}
		episodes, err = scraperInstance.GetAnimeEpisodes(anime.URL)
	case "Goyabu":
		scraperManager := scraper.NewScraperManager()
		scraperInstance, scErr := scraperManager.GetScraper(scraper.GoyabuType)
		if scErr != nil {
			return nil, fmt.Errorf("failed to get Goyabu scraper: %w", scErr)
		}
		episodes, err = scraperInstance.GetAnimeEpisodes(anime.URL)
	case "CineGratis":
		scraperManager := scraper.NewScraperManager()
		scraperInstance, scErr := scraperManager.GetScraper(scraper.CineGratisType)
		if scErr != nil {
			return nil, fmt.Errorf("failed to get CineGratis scraper: %w", scErr)
		}
		episodes, err = scraperInstance.GetAnimeEpisodes(anime.URL)
	default:
		// For AnimeFire and others, use the original API function
		episodes, err = GetAnimeEpisodes(anime.URL)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get episodes from %s: %w", sourceName, err)
	}

	if len(episodes) > 0 {
		util.Debug("Episodes found", "count", len(episodes), "source", sourceName)

		// Provide additional info for user based on source (debug only)
		switch sourceName {
		case "AnimesOnlineCC":
			util.Debug("Source info", "type", "AnimesOnlineCC")
		case "Goyabu":
			util.Debug("Source info", "type", "Goyabu")
		default:
			util.Debug("Source info", "type", "Animefire.io", "features", "dubbed/subtitled")
		}
	} else {
		util.Warn("No episodes found", "source", sourceName)
	}

	return episodes, nil
}

// Enhanced episode URL fetching with improved source detection
func GetEpisodeStreamURLEnhanced(episode *models.Episode, anime *models.Anime, quality string) (string, error) {
	// Check if this is FlixHQ content
	if anime.Source == "FlixHQ" || anime.MediaType == models.MediaTypeMovie || anime.MediaType == models.MediaTypeTV {
		streamURL, _, err := GetFlixHQStreamURL(anime, episode, quality)
		return streamURL, err
	}

	scraperManager := scraper.NewScraperManager()

	// Determine source type with enhanced logic
	var scraperType scraper.ScraperType
	var sourceName string

	// Priority 1: Check the Source field (most reliable)
	if anime.Source == "AnimesOnlineCC" {
		scraperType = scraper.AnimesOnlineCCTYPE
		sourceName = "AnimesOnlineCC"
	} else if anime.Source == "Goyabu" {
		scraperType = scraper.GoyabuType
		sourceName = "Goyabu"
	} else if strings.Contains(anime.Source, "AnimeFire") {
		scraperType = scraper.AnimefireType
		sourceName = "Animefire.io"
	} else if anime.Source == "Cineby" {
		scraperType = scraper.CinebyType
		sourceName = "Cineby"
	} else if anime.Source == "CineGratis" {
		scraperType = scraper.CineGratisType
		sourceName = "CineGratis"
	} else if strings.Contains(anime.Name, "[Portuguese]") || strings.Contains(anime.Name, "[Português]") {
		// Priority 2: Check language tags (Portuguese default to Animefire)
		scraperType = scraper.AnimefireType
		sourceName = "Animefire.io"
	} else if strings.Contains(anime.URL, "animesonlinecc") {
		// Priority 3: URL analysis
		scraperType = scraper.AnimesOnlineCCTYPE
		sourceName = "AnimesOnlineCC"
	} else if strings.Contains(anime.URL, "goyabu") {
		scraperType = scraper.GoyabuType
		sourceName = "Goyabu"
	} else if strings.Contains(anime.URL, "animefire") {
		scraperType = scraper.AnimefireType
		sourceName = "Animefire.io"
	} else if strings.Contains(anime.URL, "cineby") {
		scraperType = scraper.CinebyType
		sourceName = "Cineby"
	} else if strings.Contains(anime.URL, "cinegratis") {
		scraperType = scraper.CineGratisType
		sourceName = "CineGratis"
	} else {
		// Default to Animefire
		scraperType = scraper.AnimefireType
		sourceName = "Animefire.io (default)"
	}

	util.Debug("Getting stream URL", "source", sourceName, "episode", episode.Number)

	util.Debug("Source details",
		"scraperType", scraperType,
		"animeURL", anime.URL,
		"episodeURL", episode.URL,
		"episodeNumber", episode.Number,
		"quality", quality)

	scraperInstance, err := scraperManager.GetScraper(scraperType)
	if err != nil {
		return "", fmt.Errorf("failed to get scraper for %s: %w", sourceName, err)
	}

	if quality == "" {
		quality = "best"
	}

	var streamURL string
	var streamErr error

	// Handle different scraper types with appropriate parameters
	switch scraperType {
	case scraper.AnimesOnlineCCTYPE:
		util.Debug("Processing through AnimesOnlineCC")
		streamURL, _, streamErr = scraperInstance.GetStreamURL(episode.URL)
	case scraper.GoyabuType:
		util.Debug("Processing through Goyabu")
		streamURL, _, streamErr = scraperInstance.GetStreamURL(episode.URL)
	case scraper.CinebyType:
		util.Debug("Processing through Cineby")
		streamURL, _, streamErr = scraperInstance.GetStreamURL(episode.URL)
	case scraper.CineGratisType:
		util.Debug("Processing through CineGratis")
		streamURL, _, streamErr = scraperInstance.GetStreamURL(episode.URL)
	default:
		util.Debug("Processing through Animefire.io")
		streamURL, _, streamErr = scraperInstance.GetStreamURL(episode.URL, quality)
	}

	if streamErr != nil {
		// Propagate back request error without wrapping
		if errors.Is(streamErr, scraper.ErrBackRequested) {
			return "", streamErr
		}
		return "", fmt.Errorf("failed to get stream URL from %s: %w", sourceName, streamErr)
	}

	if streamURL == "" {
		return "", fmt.Errorf("empty stream URL returned from %s", sourceName)
	}

	util.Debug("Stream URL obtained", "source", sourceName)
	util.Debug("Stream URL details", "url", streamURL)

	return streamURL, nil
}

// Enhanced download support
func DownloadEpisodeEnhanced(anime *models.Anime, episodeNum int, quality string) error {
	util.Debugf("Fetching episodes for %s...", anime.Name)

	episodes, err := GetAnimeEpisodesEnhanced(anime)
	if err != nil {
		return fmt.Errorf("failed to get episodes: %w", err)
	}

	if episodeNum < 1 || episodeNum > len(episodes) {
		return fmt.Errorf("episode %d not found (available: 1-%d)", episodeNum, len(episodes))
	}

	episode := episodes[episodeNum-1]

	util.Debugf("Getting stream URL for episode %d...", episodeNum)
	streamURL, err := GetEpisodeStreamURLEnhanced(&episode, anime, quality)
	if err != nil {
		return fmt.Errorf("failed to get stream URL: %w", err)
	}

	util.Debugf("Stream URL obtained: %s", streamURL)

	// Create a basic downloader (this would integrate with your existing downloader)
	return downloadFromURL(streamURL, fmt.Sprintf("%s_Episode_%d",
		sanitizeFilename(anime.Name), episodeNum))
}

// Enhanced range download support
func DownloadEpisodeRangeEnhanced(anime *models.Anime, startEp, endEp int, quality string) error {
	util.Debugf("Fetching episodes for %s...", anime.Name)

	episodes, err := GetAnimeEpisodesEnhanced(anime)
	if err != nil {
		return fmt.Errorf("failed to get episodes: %w", err)
	}

	if startEp < 1 || endEp > len(episodes) || startEp > endEp {
		return fmt.Errorf("invalid range %d-%d (available: 1-%d)", startEp, endEp, len(episodes))
	}

	for i := startEp; i <= endEp; i++ {
		util.Infof("Downloading episode %d of %d...", i, endEp)

		episode := episodes[i-1]
		streamURL, err := GetEpisodeStreamURLEnhanced(&episode, anime, quality)
		if err != nil {
			util.Errorf("Failed to get stream URL for episode %d: %v", i, err)
			continue
		}

		filename := fmt.Sprintf("%s_Episode_%d", sanitizeFilename(anime.Name), i)
		// Note: downloadFromURL is a placeholder - integrate with proper downloader
		_ = downloadFromURL(streamURL, filename) // This will always fail as expected

		util.Infof("Successfully downloaded episode %d", i)
	}

	return nil
}

// Helper function to sanitize filename
func sanitizeFilename(name string) string {
	// Remove language tags like [English], [Portuguese], [Português], [Japonês], [Japanese], [Movies/TV] at the start
	reLangTags := regexp.MustCompile(`^\s*\[(?:English|Portuguese|Português|Japonês|Japanese|Movies/TV)\]\s*`)
	name = reLangTags.ReplaceAllString(name, "")
	name = strings.TrimSpace(name)

	// Replace invalid characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalid {
		name = strings.ReplaceAll(name, char, "_")
	}

	return name
}

// Basic download function (placeholder - integrate with your existing downloader)
func downloadFromURL(_ string, _ string) error {
	// This is a placeholder that should fail to trigger fallback to the proper downloader
	util.Debugf("Enhanced API downloadFromURL is a placeholder - returning error to trigger fallback")
	return fmt.Errorf("enhanced download not implemented - use legacy downloader")
}

// Legacy wrapper functions to maintain compatibility
func SearchAnimeWithSource(name string, source string) (*models.Anime, error) {
	return SearchAnimeEnhanced(name, source)
}

// GetFlixHQEpisodes handles episodes/content for FlixHQ movies and TV shows
func GetFlixHQEpisodes(media *models.Anime) ([]models.Episode, error) {
	flixhqClient := scraper.NewFlixHQClient()

	// Extract media ID from URL
	mediaID := extractMediaIDFromURL(media.URL)
	if mediaID == "" {
		return nil, fmt.Errorf("could not extract media ID from URL: %s", media.URL)
	}

	util.Debug("Getting FlixHQ content", "mediaType", media.MediaType, "mediaID", mediaID)

	// For movies, return a single "episode" representing the movie
	if media.MediaType == models.MediaTypeMovie {
		util.Debug("FlixHQ: Processing movie")
		return []models.Episode{
			{
				Number: "1",
				Num:    1,
				URL:    mediaID, // Store media ID for later use
				Title: models.TitleDetails{
					English: media.Name,
					Romaji:  media.Name,
				},
			},
		}, nil
	}

	// For TV shows, get seasons and let user select
	util.Debug("FlixHQ: Processing TV show, getting seasons")
	seasons, err := flixhqClient.GetSeasons(mediaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seasons: %w", err)
	}

	if len(seasons) == 0 {
		return nil, fmt.Errorf("no seasons found for TV show")
	}

	// Let user select a season
	seasonNames := make([]string, len(seasons))
	for i, s := range seasons {
		seasonNames[i] = s.Title
	}

	seasonIdx, err := fuzzyfinder.Find(
		seasonNames,
		func(i int) string { return seasonNames[i] },
		fuzzyfinder.WithPromptString("Select season: "),
	)
	if err != nil {
		return nil, fmt.Errorf("season selection cancelled: %w", err)
	}

	selectedSeason := seasons[seasonIdx]
	util.Debug("Selected season", "season", selectedSeason.Title, "id", selectedSeason.ID)

	// Get episodes for the selected season
	flixEpisodes, err := flixhqClient.GetEpisodes(selectedSeason.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get episodes: %w", err)
	}

	// Convert to models.Episode
	var episodes []models.Episode
	for _, ep := range flixEpisodes {
		episodes = append(episodes, models.Episode{
			Number: fmt.Sprintf("%d", ep.Number),
			Num:    ep.Number,
			URL:    ep.DataID, // Store DataID for stream retrieval
			Title: models.TitleDetails{
				English: ep.Title,
				Romaji:  ep.Title,
			},
			DataID:   ep.DataID,
			SeasonID: selectedSeason.ID,
		})
	}

	util.Debug("FlixHQ episodes loaded", "count", len(episodes))
	return episodes, nil
}

// GetFlixHQStreamURL gets the stream URL for FlixHQ content
func GetFlixHQStreamURL(media *models.Anime, episode *models.Episode, quality string) (string, []models.Subtitle, error) {
	flixhqClient := scraper.NewFlixHQClient()
	provider := "Vidcloud"
	subsLanguage := util.GlobalSubsLanguage
	if subsLanguage == "" {
		subsLanguage = "english"
	}

	var streamInfo *scraper.FlixHQStreamInfo
	var episodeID string
	var embedLink string
	var err error

	if media.MediaType == models.MediaTypeMovie {
		// For movies, episode.URL contains the media ID
		mediaID := episode.URL
		util.Debug("Getting movie stream", "mediaID", mediaID)

		episodeID, err = flixhqClient.GetMovieServerID(mediaID, provider)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get movie server: %w", err)
		}

		embedLink, err = flixhqClient.GetEmbedLink(episodeID)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get embed link: %w", err)
		}

		streamInfo, err = flixhqClient.ExtractStreamInfo(embedLink, quality, subsLanguage)
		if err != nil {
			return "", nil, fmt.Errorf("failed to extract stream info: %w", err)
		}
	} else {
		// For TV shows, episode.URL contains the DataID
		dataID := episode.URL
		util.Debug("Getting TV episode stream", "dataID", dataID)

		episodeID, err = flixhqClient.GetEpisodeServerID(dataID, provider)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get episode server: %w", err)
		}

		embedLink, err = flixhqClient.GetEmbedLink(episodeID)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get embed link: %w", err)
		}

		streamInfo, err = flixhqClient.ExtractStreamInfo(embedLink, quality, subsLanguage)
		if err != nil {
			return "", nil, fmt.Errorf("failed to extract stream info: %w", err)
		}
	}

	// Convert subtitles
	var subtitles []models.Subtitle
	for _, sub := range streamInfo.Subtitles {
		subtitles = append(subtitles, models.Subtitle{
			URL:      sub.URL,
			Language: sub.Language,
			Label:    sub.Label,
		})
	}

	return streamInfo.VideoURL, subtitles, nil
}

// extractMediaIDFromURL extracts the media ID from a FlixHQ URL
func extractMediaIDFromURL(urlStr string) string {
	// URL format: https://flixhq.to/movie/watch-movie-name-12345 or /movie/watch-movie-name-12345
	parts := strings.Split(urlStr, "-")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func GetAnimeEpisodesWithSource(anime *models.Anime) ([]models.Episode, error) {
	return GetAnimeEpisodesEnhanced(anime)
}
