package util

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the color palette and styles for the TUI
type Theme struct {
	Name string

	// Base colors
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Background lipgloss.Color
	Surface    lipgloss.Color
	Text       lipgloss.Color
	Subtext    lipgloss.Color
	Muted      lipgloss.Color

	// Semantic colors
	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color
	Info    lipgloss.Color
}

// Built-in themes
var (
	ThemeTokyoNight = Theme{
		Name:       "Tokyo Night",
		Primary:    lipgloss.Color("#7AA2F7"),
		Secondary:  lipgloss.Color("#BB9AF7"),
		Accent:     lipgloss.Color("#7DCFFF"),
		Background: lipgloss.Color("#1A1B2E"),
		Surface:    lipgloss.Color("#24283B"),
		Text:       lipgloss.Color("#C0CAF5"),
		Subtext:    lipgloss.Color("#A9B1D6"),
		Muted:      lipgloss.Color("#565F89"),
		Success:    lipgloss.Color("#9ECE6A"),
		Warning:    lipgloss.Color("#E0AF68"),
		Error:      lipgloss.Color("#F7768E"),
		Info:       lipgloss.Color("#7DCFFF"),
	}

	ThemeCatppuccinMocha = Theme{
		Name:       "Catppuccin Mocha",
		Primary:    lipgloss.Color("#CBA6F7"),
		Secondary:  lipgloss.Color("#89B4FA"),
		Accent:     lipgloss.Color("#89DCEB"),
		Background: lipgloss.Color("#1E1E2E"),
		Surface:    lipgloss.Color("#313244"),
		Text:       lipgloss.Color("#CDD6F4"),
		Subtext:    lipgloss.Color("#BAC2DE"),
		Muted:      lipgloss.Color("#585B70"),
		Success:    lipgloss.Color("#A6E3A1"),
		Warning:    lipgloss.Color("#FAB387"),
		Error:      lipgloss.Color("#F38BA8"),
		Info:       lipgloss.Color("#89DCEB"),
	}

	ThemeGruvbox = Theme{
		Name:       "Gruvbox",
		Primary:    lipgloss.Color("#83A598"),
		Secondary:  lipgloss.Color("#D3869B"),
		Accent:     lipgloss.Color("#8EC07C"),
		Background: lipgloss.Color("#282828"),
		Surface:    lipgloss.Color("#3C3836"),
		Text:       lipgloss.Color("#EBDBB2"),
		Subtext:    lipgloss.Color("#D5C4A1"),
		Muted:      lipgloss.Color("#928374"),
		Success:    lipgloss.Color("#B8BB26"),
		Warning:    lipgloss.Color("#FABD2F"),
		Error:      lipgloss.Color("#FB4934"),
		Info:       lipgloss.Color("#83A598"),
	}

	ThemeDracula = Theme{
		Name:       "Dracula",
		Primary:    lipgloss.Color("#BD93F9"),
		Secondary:  lipgloss.Color("#FF79C6"),
		Accent:     lipgloss.Color("#8BE9FD"),
		Background: lipgloss.Color("#282A36"),
		Surface:    lipgloss.Color("#44475A"),
		Text:       lipgloss.Color("#F8F8F2"),
		Subtext:    lipgloss.Color("#BFBFBF"),
		Muted:      lipgloss.Color("#6272A4"),
		Success:    lipgloss.Color("#50FA7B"),
		Warning:    lipgloss.Color("#FFB86C"),
		Error:      lipgloss.Color("#FF5555"),
		Info:       lipgloss.Color("#8BE9FD"),
	}

	AvailableThemes = []Theme{
		ThemeTokyoNight,
		ThemeCatppuccinMocha,
		ThemeGruvbox,
		ThemeDracula,
	}
)

// GlobalTheme is the active theme, defaulting to Tokyo Night
var GlobalTheme = ThemeTokyoNight

// SetTheme sets the active global theme
func SetTheme(t Theme) {
	GlobalTheme = t
}

// GetThemeByName finds a theme by name (case-insensitive)
func GetThemeByName(name string) (Theme, bool) {
	for _, t := range AvailableThemes {
		if t.Name == name {
			return t, true
		}
	}
	return GlobalTheme, false
}

// ─── Styled Builders ───────────────────────────────────────────────────────

// TitleStyle returns a bold title style using the primary color
func TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(GlobalTheme.Primary).
		Bold(true).
		MarginBottom(1)
}

// SubtitleStyle returns a secondary style for subtitles/labels
func SubtitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(GlobalTheme.Secondary).
		Bold(true)
}

// BadgeStyle returns a badge/pill style for status indicators
func BadgeStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(color).
		Foreground(GlobalTheme.Background).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)
}

// BoxStyle returns a bordered card box
func BoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GlobalTheme.Primary).
		Padding(0, 1)
}

// MutedStyle returns a dimmed muted text style
func MutedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(GlobalTheme.Muted)
}

// SuccessStyle returns green success text
func SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(GlobalTheme.Success).Bold(true)
}

// WarningStyle returns yellow warning text
func WarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(GlobalTheme.Warning).Bold(true)
}

// ErrorStyle returns red error text  
func ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(GlobalTheme.Error).Bold(true)
}

// InfoStyle returns cyan info text
func InfoStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(GlobalTheme.Info)
}

// AccentStyle returns accent highlighted text
func AccentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(GlobalTheme.Accent).Bold(true)
}

// HeaderBannerStyle returns the full-width app header banner style
func HeaderBannerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(GlobalTheme.Primary).
		Foreground(GlobalTheme.Background).
		Bold(true).
		Padding(0, 2).
		MarginBottom(1)
}
