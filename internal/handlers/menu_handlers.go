package handlers

import (
	"fmt"

	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// HandleWatchlistMode shows the user's watchlist and allows selecting an item to play
func HandleWatchlistMode() error {
	tracker := tracking.GetGlobalTracker()
	if tracker == nil {
		util.Warnf("Tracker não inicializado — dados de watchlist indisponíveis.")
		return util.ErrBackToMainMenu
	}

	list, err := tracker.GetAllFollowedMedia()
	if err != nil {
		return fmt.Errorf("failed to get watchlist: %w", err)
	}

	if len(list) == 0 {
		banner := util.BoxStyle().Render(
			util.TitleStyle().Render("📂 Minha Lista") + "\n\n" +
				util.MutedStyle().Render("Sua lista está vazia.\nPressione ⭐ na seleção de episódios para adicionar obras."),
		)
		fmt.Println(banner)
		var choice bool
		huh.NewConfirm().
			Title("Voltar ao menu?").
			Value(&choice).
			Run()
		return util.ErrBackToMainMenu
	}

	// Build header
	fmt.Println(util.TitleStyle().Render("📂 Minha Lista de Obras"))

	var options []huh.Option[string]
	options = append(options, huh.NewOption("<< Voltar ao Menu Principal", "back"))
	for _, m := range list {
		var statusColor lipgloss.Color
		switch m.Status {
		case "watching":
			statusColor = util.GlobalTheme.Success
		case "planned":
			statusColor = util.GlobalTheme.Info
		case "dropped":
			statusColor = util.GlobalTheme.Error
		default:
			statusColor = util.GlobalTheme.Muted
		}
		badge := lipgloss.NewStyle().Foreground(statusColor).Bold(true).Render("[" + m.Status + "]")
		label := fmt.Sprintf("%s %s — ep. %d/%d", badge, m.Title, m.LastEpisode, m.TotalEpisodes)
		options = append(options, huh.NewOption(label, m.Title))
	}

	var selected string
	err = huh.NewSelect[string]().
		Title("Selecione para assistir:").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil || selected == "back" {
		return util.ErrBackToMainMenu
	}

	return HandlePlaybackMode(selected)
}

// HandleContinueWatchingMode shows media with progress and allows resuming
func HandleContinueWatchingMode() error {
	tracker := tracking.GetGlobalTracker()
	if tracker == nil {
		util.Warnf("Tracker não inicializado.")
		return util.ErrBackToMainMenu
	}

	list, err := tracker.GetAllAnime()
	if err != nil {
		return fmt.Errorf("failed to get progress: %w", err)
	}

	if len(list) == 0 {
		fmt.Println(util.BoxStyle().Render(
			util.TitleStyle().Render("🕒 Continuar Assistindo") + "\n\n" +
				util.MutedStyle().Render("Você ainda não assistiu nada.\nComece a assistir algo e seu progresso aparecerá aqui."),
		))
		return util.ErrBackToMainMenu
	}

	fmt.Println(util.TitleStyle().Render("🕒 Continuar Assistindo"))

	var options []huh.Option[string]
	options = append(options, huh.NewOption("<< Voltar ao Menu Principal", "back"))
	for _, m := range list {
		progress := 0.0
		if m.Duration > 0 {
			progress = (float64(m.PlaybackTime) / float64(m.Duration)) * 100
		}
		var pctStyle lipgloss.Style
		switch {
		case progress >= 85:
			pctStyle = util.SuccessStyle()
		case progress >= 40:
			pctStyle = util.WarningStyle()
		default:
			pctStyle = util.MutedStyle()
		}
		pct := pctStyle.Render(fmt.Sprintf("%.0f%%", progress))
		label := fmt.Sprintf("%s — Ep %d (%s concluído)", m.Title, m.EpisodeNumber, pct)
		options = append(options, huh.NewOption(label, m.Title))
	}

	var selected string
	err = huh.NewSelect[string]().
		Title("Selecione para retomar:").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil || selected == "back" {
		return util.ErrBackToMainMenu
	}

	return HandlePlaybackMode(selected)
}

// HandleScraperHealthMode shows a rich visual health report for all scrapers
func HandleScraperHealthMode() error {
	tracker := tracking.GetGlobalTracker()
	if tracker == nil {
		util.Warnf("Tracker não inicializado — dados de saúde indisponíveis.")
		return util.ErrBackToMainMenu
	}

	records, err := tracker.GetScraperHealthRecords()
	if err != nil {
		return fmt.Errorf("failed to get health records: %w", err)
	}

	title := util.TitleStyle().Render("📊 Saúde dos Plugins de Busca")
	divider := util.MutedStyle().Render("─────────────────────────────────────────────────────────\n")

	report := "\n" + title + "\n" + divider

	if len(records) == 0 {
		report += util.MutedStyle().Render("  Nenhum dado coletado ainda.\n  Use a busca para que o app comece a monitorar os plugins.\n")
	} else {
		for _, r := range records {
			successRate := 100.0
			if r.TotalSearches > 0 {
				successRate = (float64(r.TotalSearches-r.FailedSearches) / float64(r.TotalSearches)) * 100
			}

			var badge, barStr string
			switch {
			case successRate >= 90:
				badge = lipgloss.NewStyle().Background(util.GlobalTheme.Success).Foreground(util.GlobalTheme.Background).Bold(true).Padding(0, 1).Render(" OK ")
				barStr = util.SuccessStyle().Render(renderBar(successRate, 22))
			case successRate >= 60:
				badge = lipgloss.NewStyle().Background(util.GlobalTheme.Warning).Foreground(util.GlobalTheme.Background).Bold(true).Padding(0, 1).Render("INST")
				barStr = util.WarningStyle().Render(renderBar(successRate, 22))
			default:
				badge = lipgloss.NewStyle().Background(util.GlobalTheme.Error).Foreground(util.GlobalTheme.Background).Bold(true).Padding(0, 1).Render("CRIT")
				barStr = util.ErrorStyle().Render(renderBar(successRate, 22))
			}

			name := util.SubtitleStyle().Render(fmt.Sprintf("%-16s", r.ScraperName))
			stats := util.MutedStyle().Render(fmt.Sprintf("  %d/%d ok", r.TotalSearches-r.FailedSearches, r.TotalSearches))

			report += fmt.Sprintf(" %s %s %s%s\n", badge, name, barStr, stats)
			if r.LastFailure != "" && successRate < 90 {
				report += util.ErrorStyle().Render(fmt.Sprintf("     ↳ %s", truncate(r.LastFailure, 70))) + "\n"
			}
			report += "\n"
		}
	}

	report += divider
	report += util.MutedStyle().Render(" Dica: plugins com status CRIT precisam de atualização.\n")

	fmt.Println(report)

	var back bool
	huh.NewConfirm().
		Title("Voltar ao menu principal?").
		Affirmative("Sim").
		Negative("Sair").
		Value(&back).
		Run()

	if back {
		return util.ErrBackToMainMenu
	}
	return nil
}

// renderBar renders a UTF-8 progress bar
func renderBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled
	bar := "["
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}
	return fmt.Sprintf("%s] %5.1f%%", bar, pct)
}

// truncate truncates a string to max length
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
