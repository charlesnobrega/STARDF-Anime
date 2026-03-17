package handlers

import (
	"fmt"

	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charmbracelet/huh"
)

// HandleThemeMode shows theme selection UI
func HandleThemeMode() error {
	themes := util.AvailableThemes

	var options []huh.Option[string]
	for _, t := range themes {
		current := ""
		if t.Name == util.GlobalTheme.Name {
			current = " ✓"
		}
		options = append(options, huh.NewOption(t.Name+current, t.Name))
	}
	options = append(options, huh.NewOption("<< Voltar ao Menu Principal", "back"))

	var selected string
	err := huh.NewSelect[string]().
		Title("🎨 Selecionar Tema").
		Description("Escolha um tema visual para a interface:").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return err
	}

	if selected == "back" {
		return util.ErrBackToMainMenu
	}

	if t, ok := util.GetThemeByName(selected); ok {
		util.SetTheme(t)
		util.InitLogger() // Re-initialize logger with new theme colors
		util.Infof("✨ Tema aplicado: %s", t.Name)

		// Persist theme choice
		if tracker := tracking.GetGlobalTracker(); tracker != nil {
			_ = tracker.SetConfig("theme", t.Name)
		}
	}

	return util.ErrBackToMainMenu
}

// RenderHealthDashboard renders a pretty health report using the active theme
func RenderHealthDashboard(records []interface{ GetName() string; GetRate() float64; GetTotal() int; GetFailed() int; GetLastErr() string }) string {
	title := util.TitleStyle().Render("📊 Relatório de Saúde dos Plugins")
	divider := util.MutedStyle().Render("──────────────────────────────────────────────────\n")

	result := "\n" + title + "\n" + divider

	for _, r := range records {
		rate := r.GetRate()
		var statusBadge string
		var bar string

		switch {
		case rate >= 90:
			statusBadge = util.BadgeStyle(util.GlobalTheme.Success).Render("OK")
			bar = util.SuccessStyle().Render(progressBar(rate, 20))
		case rate >= 60:
			statusBadge = util.BadgeStyle(util.GlobalTheme.Warning).Render("INST")
			bar = util.WarningStyle().Render(progressBar(rate, 20))
		default:
			statusBadge = util.BadgeStyle(util.GlobalTheme.Error).Render("CRIT")
			bar = util.ErrorStyle().Render(progressBar(rate, 20))
		}

		name := util.SubtitleStyle().Render(fmt.Sprintf("%-16s", r.GetName()))
		stats := util.MutedStyle().Render(fmt.Sprintf("%d/%d ok", r.GetTotal()-r.GetFailed(), r.GetTotal()))

		result += fmt.Sprintf("%s %s %s %s\n", statusBadge, name, bar, stats)

		if r.GetLastErr() != "" && rate < 90 {
			result += util.ErrorStyle().Render("   ↳ "+r.GetLastErr()) + "\n"
		}
	}

	result += "\n" + divider
	return result
}

// progressBar renders a simple ASCII progress bar
func progressBar(pct float64, width int) string {
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
	bar += fmt.Sprintf("] %.0f%%", pct)
	return bar
}
