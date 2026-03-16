package util

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Help styles using lipgloss
var (
	// Professional and modern color palette
	lightGreen  = lipgloss.Color("#90EE90") // Soft light green
	gray        = lipgloss.Color("#A9A9A9") // Medium gray
	darkGray    = lipgloss.Color("#5A5A5A") // Dark gray for details
	brightGreen = lipgloss.Color("#00FF7F") // Bright green for highlights
	blue        = lipgloss.Color("#6366F1") // Modern blue (matches logger prefix)

	// Text styles
	titleStyle = lipgloss.NewStyle().
			Foreground(blue). // Title in blue (matching StarDF-Anime prefix)
			Bold(true).
			PaddingBottom(1).
			MarginLeft(2)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(gray).
			Italic(true).
			PaddingBottom(1).
			MarginLeft(2)

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(lightGreen). // Section titles in light green
				Bold(true).
				PaddingLeft(2)

	commandStyle = lipgloss.NewStyle().
			Foreground(brightGreen). // Commands in bright green
			Bold(true).
			PaddingLeft(4)

	optionStyle = lipgloss.NewStyle().
			Foreground(brightGreen). // Options in bright green
			Bold(true).
			PaddingLeft(4)

	parameterStyle = lipgloss.NewStyle().
			Foreground(gray). // Parameters in gray to differentiate
			Italic(true)

	descriptionStyle = lipgloss.NewStyle().
				Foreground(gray). // Descriptions in gray
				PaddingLeft(6).
				Width(80 - 6) // Adjust width for line wrapping

	exampleStyle = lipgloss.NewStyle().
			Foreground(darkGray). // Examples in dark gray
			Italic(true).
			PaddingLeft(8)

	separatorStyle = lipgloss.NewStyle().
			Foreground(darkGray) // Separators in dark gray
)

// ShowBeautifulHelp displays a beautifully formatted help message
func ShowBeautifulHelp() {
	var helpContent strings.Builder

	// Program title
	helpContent.WriteString(titleStyle.Render("StarDF-Anime - Beautiful Anime Streaming CLI"))
	helpContent.WriteString("\n")
	helpContent.WriteString(subtitleStyle.Render("Watch your favorite anime directly from the terminal with style and ease."))
	helpContent.WriteString("\n\n")

	// Usage section
	helpContent.WriteString(separatorStyle.Render(strings.Repeat("─", 80)))
	helpContent.WriteString("\n")
	helpContent.WriteString(sectionTitleStyle.Render("Usage:"))
	helpContent.WriteString("\n")
	helpContent.WriteString(commandStyle.Render("  stardf-anime"))
	helpContent.WriteString("\n")
	helpContent.WriteString(descriptionStyle.Render("    Interactive mode - search and select anime from a beautiful menu"))
	helpContent.WriteString("\n")
	helpContent.WriteString(commandStyle.Render("  stardf-anime ") + parameterStyle.Render("[options]"))
	helpContent.WriteString("\n")
	helpContent.WriteString(descriptionStyle.Render("    Run with specific options"))
	helpContent.WriteString("\n")
	helpContent.WriteString(commandStyle.Render("  stardf-anime ") + parameterStyle.Render("[options] [anime name]"))
	helpContent.WriteString("\n")
	helpContent.WriteString(descriptionStyle.Render("    Direct search for anime (use spaces, not hyphens)"))
	helpContent.WriteString("\n")
	helpContent.WriteString(exampleStyle.Render("Example: stardf-anime \"one piece\" (not \"one-piece\")"))
	helpContent.WriteString("\n\n")

	// Options section
	helpContent.WriteString(separatorStyle.Render(strings.Repeat("─", 80)))
	helpContent.WriteString("\n")
	helpContent.WriteString(sectionTitleStyle.Render("Options:"))
	helpContent.WriteString("\n")
	addOption(&helpContent, "--debug", "Enable debug mode for detailed error information and performance metrics.")
	addOption(&helpContent, "--perf", "Enable performance profiling - shows timing metrics for all operations.")
	addOption(&helpContent, "--help / -h", "Display this beautiful help message with detailed usage information.")
	addOption(&helpContent, "--version", "Show version information and build details.")
	addOption(&helpContent, "--update", "Check for updates and update automatically to the latest version.")
	addOption(&helpContent, "-d", "Download mode - download specific episodes for offline viewing.")
	addOption(&helpContent, "-r", "Range download mode - download multiple episodes (use with -d).")
	addOption(&helpContent, "--source", "Specify anime source (allanime, animefire). Default: search all sources.")
	addOption(&helpContent, "--quality", "Specify video quality (best, worst, 720p, 1080p, etc.). Default: best.")
	addOption(&helpContent, "--allanime-smart", "AllAnime Smart Range: auto-skip intros/outros via AniSkip and use priority mirrors.")
	addOption(&helpContent, "--type", "Specify media type (anime, movie, tv). Default: anime.")
	addOption(&helpContent, "--subs", "Specify subtitle language for movies/TV shows (FlixHQ only: english, spanish, portuguese, etc.).")
	addOption(&helpContent, "--audio", "Specify preferred audio language for movies/TV (FlixHQ only: pt-BR,english,spanish).")
	helpContent.WriteString("\n")

	// Features section
	helpContent.WriteString(separatorStyle.Render(strings.Repeat("─", 80)))
	helpContent.WriteString("\n")
	helpContent.WriteString(sectionTitleStyle.Render("Features:"))
	helpContent.WriteString("\n")

	addFeature(&helpContent, "Multi-Source Support", "Stream from AllAnime, AnimeFire, and FlixHQ (movies/TV) with automatic fallback.")
	addFeature(&helpContent, "Movies & TV Shows", "Watch movies and TV series alongside anime using FlixHQ integration.")
	addFeature(&helpContent, "Smart Search", "Intelligent search with fuzzy matching and suggestions.")
	addFeature(&helpContent, "Quality Selection", "Choose video quality from multiple available sources.")
	addFeature(&helpContent, "Batch Downloads", "Download single episodes or entire seasons for offline viewing.")
	addFeature(&helpContent, "Interactive Controls", "Beautiful terminal interface with keyboard navigation.")
	addFeature(&helpContent, "Discord Rich Presence", "Show your friends what you're watching.")
	addFeature(&helpContent, "Progress Tracking", "Keep track of your watch progress and episode history.")
	addFeature(&helpContent, "Skip Intros", "Automatically skip anime intros and outros.")
	addFeature(&helpContent, "Subtitle Support", "Multilingual subtitle support for movies and TV shows.")
	addFeature(&helpContent, "Audio Track Selection", "Select preferred audio language for movies/TV during playback (FlixHQ only).")
	addFeature(&helpContent, "AllAnime Smart Range", "Exclusive: For AllAnime, download a range with mirror priority and optional intro/outro trimming.")
	helpContent.WriteString("\n")

	// Examples section
	helpContent.WriteString(separatorStyle.Render(strings.Repeat("─", 80)))
	helpContent.WriteString("\n")
	helpContent.WriteString(sectionTitleStyle.Render("Examples:"))
	helpContent.WriteString("\n")
	addExample(&helpContent, "stardf-anime", "Start interactive mode")
	addExample(&helpContent, "stardf-anime \"attack on titan\"", "Search directly for Attack on Titan")
	addExample(&helpContent, "stardf-anime --debug \"naruto\"", "Search with debug information")
	addExample(&helpContent, "stardf-anime --update", "Check for updates and update automatically")
	addExample(&helpContent, "stardf-anime --version", "Show version information")
	addExample(&helpContent, "stardf-anime -d \"one piece\" 1", "Download episode 1 of One Piece")
	addExample(&helpContent, "stardf-anime -d -r \"naruto\" 1-5", "Download episodes 1-5 of Naruto")
	addExample(&helpContent, "stardf-anime -d --source allanime \"bleach\" 10", "Download from AllAnime specifically")
	addExample(&helpContent, "stardf-anime -d --quality 720p \"demon slayer\" 1", "Download in 720p quality")
	addExample(&helpContent, "stardf-anime -d --source animefire --quality best \"jujutsu kaisen\" 5", "Use AnimeFire with best quality")
	addExample(&helpContent, "stardf-anime -d -r --source allanime --allanime-smart \"vinland saga\" 1-4", "AllAnime Smart Range for episodes 1-4")
	addExample(&helpContent, "stardf-anime --type movie \"avengers\"", "Search for movies matching 'avengers'")
	addExample(&helpContent, "stardf-anime --type tv \"breaking bad\"", "Search for TV shows matching 'breaking bad'")
	addExample(&helpContent, "stardf-anime --type movie --subs spanish \"spider-man\"", "Search movies with Spanish subtitles")
	addExample(&helpContent, "stardf-anime --type movie --audio \"pt-BR,english\" \"matrix\"", "Play movie with Portuguese audio preference")
	helpContent.WriteString("\n")

	// Footer
	helpContent.WriteString(separatorStyle.Render(strings.Repeat("─", 80)))
	helpContent.WriteString("\n")
	helpContent.WriteString(subtitleStyle.Render("For more information, visit: https://github.com/charlesnobrega/STARDF-Anime"))
	helpContent.WriteString("\n")
	helpContent.WriteString(subtitleStyle.Render("Made with love for anime lovers everywhere"))
	helpContent.WriteString("\n\n")

	// Print the complete help content
	fmt.Print(helpContent.String())
}

// Helper functions for building help content
func addOption(builder *strings.Builder, opt, desc string) {
	builder.WriteString(optionStyle.Render("  " + opt))
	builder.WriteString("\n")
	builder.WriteString(descriptionStyle.Render("    " + desc))
	builder.WriteString("\n")
}

func addFeature(builder *strings.Builder, feature, desc string) {
	builder.WriteString(commandStyle.Render("  " + feature))
	builder.WriteString("\n")
	builder.WriteString(descriptionStyle.Render("    " + desc))
	builder.WriteString("\n")
}

func addExample(builder *strings.Builder, cmd, desc string) {
	builder.WriteString(commandStyle.Render("  " + cmd))
	builder.WriteString("\n")
	builder.WriteString(descriptionStyle.Render("    " + desc))
	builder.WriteString("\n")
}
