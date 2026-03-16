package handlers

import (
	"fmt"
	"strings"

	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charmbracelet/huh"
)

// HandleWatchlistMode shows the user's watchlist and allows selecting an item to play
func HandleWatchlistMode() error {
	tracker := tracking.GetGlobalTracker()
	if tracker == nil {
		return fmt.Errorf("tracker not initialized")
	}

	list, err := tracker.GetAllFollowedMedia()
	if err != nil {
		return fmt.Errorf("failed to get watchlist: %w", err)
	}

	if len(list) == 0 {
		util.Infof("Sua lista está vazia. Adicione algo enquanto navega!")
		var choice bool
		huh.NewConfirm().
			Title("Sua lista está vazia.").
			Description("Deseja voltar ao menu principal?").
			Affirmative("Sim").
			Negative("Sair").
			Value(&choice).
			Run()
		if choice {
			return util.ErrBackToMainMenu
		}
		return nil
	}

	var options []huh.Option[string]
	options = append(options, huh.NewOption("<< Voltar ao Menu Principal", "back"))
	for _, m := range list {
		label := fmt.Sprintf("[%s] %s (Ep: %d)", m.Status, m.Title, m.LastEpisode)
		options = append(options, huh.NewOption(label, m.Title))
	}

	var selected string
	err = huh.NewSelect[string]().
		Title("Sua Lista de Obras").
		Description("Selecione para assistir:").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return err
	}

	if selected == "back" {
		return util.ErrBackToMainMenu
	}

	// Route to playback mode with the selected name
	return HandlePlaybackMode(selected)
}

// HandleContinueWatchingMode shows media with progress and allows resuming
func HandleContinueWatchingMode() error {
	tracker := tracking.GetGlobalTracker()
	if tracker == nil {
		return fmt.Errorf("tracker not initialized")
	}

	list, err := tracker.GetAllAnime()
	if err != nil {
		return fmt.Errorf("failed to get progress: %w", err)
	}

	if len(list) == 0 {
		util.Infof("Você ainda não assistiu nada.")
		return util.ErrBackToMainMenu
	}

	// Sort by last updated (already done by SQL ideally, but let's be sure)
	var options []huh.Option[string]
	options = append(options, huh.NewOption("<< Voltar ao Menu Principal", "back"))
	for _, m := range list {
		progress := (float64(m.PlaybackTime) / float64(m.Duration)) * 100
		label := fmt.Sprintf("%s - Ep %d (%.0f%%)", m.Title, m.EpisodeNumber, progress)
		options = append(options, huh.NewOption(label, m.Title))
	}

	var selected string
	err = huh.NewSelect[string]().
		Title("Continuar Assistindo").
		Description("Selecione para retomar:").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return err
	}

	if selected == "back" {
		return util.ErrBackToMainMenu
	}

	return HandlePlaybackMode(selected)
}

// HandleScraperHealthMode shows stats about scraper reliability
func HandleScraperHealthMode() error {
	tracker := tracking.GetGlobalTracker()
	if tracker == nil {
		return fmt.Errorf("tracker not initialized")
	}

	records, err := tracker.GetScraperHealthRecords()
	if err != nil {
		return fmt.Errorf("failed to get health records: %w", err)
	}

	if len(records) == 0 {
		util.Infof("Ainda não há dados de performance dos plugins.")
		return util.ErrBackToMainMenu
	}

	var builder strings.Builder
	builder.WriteString("Status dos Plugins (Scrapers):\n\n")
	for _, r := range records {
		successRate := 100.0
		if r.TotalSearches > 0 {
			successRate = (float64(r.TotalSearches-r.FailedSearches) / float64(r.TotalSearches)) * 100
		}
		
		status := "✅ OK"
		if successRate < 70 {
			status = "⚠️ INSTÁVEL"
		}
		if successRate < 30 {
			status = "❌ CRÍTICO (Necessita Atualização)"
		}

		builder.WriteString(fmt.Sprintf("%s %-15s | Sucesso: %5.1f%% (%d/%d)\n", 
			status, r.ScraperName, successRate, r.TotalSearches-r.FailedSearches, r.TotalSearches))
		if r.LastFailure != "" {
			builder.WriteString(fmt.Sprintf("   Ultimo erro: %s\n", r.LastFailure))
		}
		builder.WriteString("\n")
	}

	fmt.Println(builder.String())
	
	var back bool
	huh.NewConfirm().
		Title("Relatório de Performance").
		Description("Deseja voltar ao menu principal?").
		Affirmative("Sim").
		Negative("Sair").
		Value(&back).
		Run()

	if back {
		return util.ErrBackToMainMenu
	}
	return nil
}
