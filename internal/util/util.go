package util

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"sync"

	"github.com/charlesnobrega/STARDF-Anime/internal/version"
	"github.com/charmbracelet/huh"
)

var (
	IsDebug             bool
	minNameLength       = 4
	ErrHelpRequested    = errors.New("help requested") // Custom error for help
	GlobalSource        string                         // Global variable to store selected source
	GlobalQuality       string                         // Global variable to store selected quality
	GlobalMediaType     string                         // Global variable to store media type (anime, movie, tv)
	GlobalSubsLanguage  string                         // Global variable to store subtitle language
	GlobalAudioLanguage string                         // Global variable to store preferred audio language
	flagsOnce           sync.Once                      // Ensure flags are only defined/parsed once
	flagsParsed         bool                           // Track if flags have been parsed
	GetMenuSubtitleFunc func() string                  // Function to get the menu subtitle
)

// Cleanup function to be called on program exit
var cleanupFuncs []func()

// RegisterCleanup registers a function to be called on program exit
func RegisterCleanup(fn func()) {
	cleanupFuncs = append(cleanupFuncs, fn)
}

// RunCleanup runs all registered cleanup functions
func RunCleanup() {
	for _, fn := range cleanupFuncs {
		fn()
	}
	// Print performance report if enabled
	if PerfEnabled {
		GetPerfTracker().PrintReport()
	}
}

// ErrorHandler returns a string with the error message, if debug mode is enabled, it will return the full error with details.
func ErrorHandler(err error) string {
	if IsDebug {
		return fmt.Sprintf("%+v", err)
	} else {
		return fmt.Sprintf("%v -- run the program with -debug to see details", err)
	}
}

// Helper prints the beautiful help message
func Helper() {
	ShowBeautifulHelp()
}

// Custom error types for different exit conditions
var (
	ErrUpdateRequested   = errors.New("update requested")
	ErrDownloadRequested = errors.New("download requested")
	ErrBackToMainMenu    = errors.New("back to main menu")
	ErrExitRequested     = errors.New("exit requested")
	ErrAniListLoginRequested  = errors.New("anilist login requested")
	ErrAniListLogoutRequested = errors.New("anilist logout requested")
	ErrBackRequested          = errors.New("back requested")
)

// MenuAction defines possible actions from the main menu
type MenuAction int

const (
	ActionSearch MenuAction = iota
	ActionWatchlist
	ActionContinue
	ActionHealth
	ActionTheme
	ActionExit
)

// MenuResult holds the result of the interactive menu selection
type MenuResult struct {
	Action     MenuAction
	SearchTerm string
}

// DownloadRequest holds download command parameters
type DownloadRequest struct {
	AnimeName     string
	EpisodeNum    int
	IsRange       bool
	StartEpisode  int
	EndEpisode    int
	Source        string // Added source field for specifying anime source
	Quality       string // Added quality field for video quality
}

// Global variable to store download request
var GlobalDownloadRequest *DownloadRequest

// ParseFlags parses the -flags and returns the anime name if provided via CLI
func ParseFlags() (string, error) {
	var animeName string
	var err error

	flagsOnce.Do(func() {
		// Define flags
		debug := flag.Bool("debug", false, "enable debug mode")
		perf := flag.Bool("perf", false, "enable performance profiling")
		help := flag.Bool("help", false, "show help message")
		altHelp := flag.Bool("h", false, "show help message")
		versionFlag := flag.Bool("version", false, "show version information")
		updateFlag := flag.Bool("update", false, "check for updates and update if available")
		downloadFlag := flag.Bool("d", false, "download mode")
		rangeFlag := flag.Bool("r", false, "download episode range (use with -d)")
		sourceFlag := flag.String("source", "", "specify media source (animefire, betteranime, topanimes/animesdigital, flixhq, cinegratis)")
		qualityFlag := flag.String("quality", "best", "specify video quality (best, worst, 720p, 1080p, etc.)")
		mediaTypeFlag := flag.String("type", "", "specify media type (anime, movie, tv)")
		subsLanguageFlag := flag.String("subs", "english", "specify subtitle language for movies/TV (FlixHQ only)")
		audioLanguageFlag := flag.String("audio", "pt-BR,pt,english", "specify preferred audio language for movies/TV (FlixHQ only)")

		anilistLoginFlag  := flag.Bool("anilist-login",  false, "connect your AniList account for automatic progress sync")
		anilistLogoutFlag := flag.Bool("anilist-logout", false, "disconnect AniList account and remove saved token")

		// Parse the flags early
		flag.Parse()

		// Set debug mode
		IsDebug = *debug

		// Set performance profiling mode
		PerfEnabled = *perf
		if PerfEnabled {
			IsDebug = true
			Debug("Performance profiling enabled")
		}

		// Store global configurations
		GlobalSource = *sourceFlag
		GlobalQuality = *qualityFlag
		GlobalMediaType = *mediaTypeFlag
		GlobalSubsLanguage = *subsLanguageFlag
		GlobalAudioLanguage = *audioLanguageFlag

		if *versionFlag || version.HasVersionArg() {
			version.ShowVersion()
			err = ErrHelpRequested
			return
		}

		if *help || *altHelp {
			Helper()
			err = ErrHelpRequested
			return
		}

		if *updateFlag {
			err = ErrUpdateRequested
			return
		}

		if *anilistLoginFlag {
			err = ErrAniListLoginRequested
			return
		}

		if *anilistLogoutFlag {
			err = ErrAniListLogoutRequested
			return
		}

		// Handle download mode
		if *downloadFlag {
			animeName, err = handleDownloadMode(*rangeFlag, *sourceFlag, *qualityFlag)
			return
		}

		if *debug {
			Debug("Debug mode is enabled")
		}

		// If the user has provided an anime name as an argument
		if len(flag.Args()) > 0 {
			name := strings.Join(flag.Args(), " ")
			if strings.Contains(name, "-") {
				name = strings.Split(name, "-")[0]
			}
			Debug("Anime name from CLI", "name", name)
			if len(name) < minNameLength {
				err = fmt.Errorf("anime name must have at least %d characters, you entered: %v", minNameLength, name)
				return
			}
			animeName = TreatingAnimeName(name)
		}
		flagsParsed = true
	})

	return animeName, err
}

// PromptInteractive shows the main menu and returns the user's choice
func PromptInteractive() (MenuResult, error) {
	items := []huh.Option[MenuAction]{
		huh.NewOption("🔍 Buscar Novo Conteúdo", ActionSearch),
		huh.NewOption("📂 Minha Lista (Watchlist)", ActionWatchlist),
		huh.NewOption("🕒 Continuar Assistindo", ActionContinue),
		huh.NewOption("📊 Saúde dos Plugins", ActionHealth),
		huh.NewOption("🎨 Temas Visuais", ActionTheme),
		huh.NewOption("❌ Sair", ActionExit),
	}

	var action MenuAction

	subtitle := ""
	if GetMenuSubtitleFunc != nil {
		subtitle = GetMenuSubtitleFunc()
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[MenuAction]().
				Title("StarDF-Anime - Menu Principal" + subtitle).
				Description("Escolha uma opção:").
				Options(items...).
				Value(&action),
		),
	)

	if err := form.Run(); err != nil {
		return MenuResult{}, err
	}

	switch action {
	case ActionSearch:
		mediaType, err := selectMediaType()
		if err != nil {
			if errors.Is(err, ErrBackRequested) {
				return PromptInteractive()
			}
			return MenuResult{}, err
		}
		GlobalMediaType = mediaType

		searchTerm, err := getUserInput("Digite o nome da obra")
		if err != nil {
			if errors.Is(err, ErrBackRequested) {
				return PromptInteractive()
			}
			return MenuResult{}, err
		}
		return MenuResult{Action: ActionSearch, SearchTerm: TreatingAnimeName(searchTerm)}, nil

	case ActionWatchlist:
		return MenuResult{Action: ActionWatchlist}, nil

	case ActionContinue:
		return MenuResult{Action: ActionContinue}, nil

	case ActionHealth:
		return MenuResult{Action: ActionHealth}, nil

	case ActionTheme:
		return MenuResult{Action: ActionTheme}, nil

	case ActionExit:
		return MenuResult{}, ErrExitRequested

	default:
		return MenuResult{}, fmt.Errorf("opção inválida")
	}
}

// FlagParser (Legacy/Deprecated) - now just calls ParseFlags and PromptInteractive if needed
func FlagParser() (MenuResult, error) {
	name, err := ParseFlags()
	if err != nil {
		return MenuResult{}, err
	}

	if name != "" {
		return MenuResult{Action: ActionSearch, SearchTerm: name}, nil
	}

	return PromptInteractive()
}

// selectMediaType asks the user what type of content they want
func selectMediaType() (string, error) {
	var choice string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("O que vamos assistir hoje?").
				Options(
					huh.NewOption("Animes (PT-BR)", "anime"),
					huh.NewOption("Filmes e Séries (PT-BR/Multi)", "movie"),
					huh.NewOption("<< Voltar", "back"),
					// huh.NewOption("Canais de TV (Em breve)", "tv"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}
	if choice == "back" {
		return "", ErrBackRequested
	}
	return choice, nil
}

// getUserInput prompts the user for input the anime name and returns it
func getUserInput(label string) (string, error) {
	var animeName string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(label).
				Description("Digite o nome e pressione Enter (vazio para voltar)").
				Value(&animeName),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	if strings.TrimSpace(animeName) == "" {
		return "", ErrBackRequested
	}

	if len(strings.TrimSpace(animeName)) < minNameLength {
		return "", fmt.Errorf("o nome deve ter pelo menos %d caracteres", minNameLength)
	}

	return animeName, nil
}

// TreatingAnimeName removes special characters and spaces from the anime name.
func TreatingAnimeName(animeName string) string {
	loweredName := strings.ToLower(animeName)
	return strings.ReplaceAll(loweredName, " ", "-")
}

func handleDownloadMode(isRange bool, source, quality string) (string, error) {
	args := flag.Args()

	if len(args) == 0 {
		return "", fmt.Errorf("download mode requires anime name and episode number/range")
	}

	if isRange {
		// Range download: stardf-anime -d -r "anime name" start-end
		if len(args) < 2 {
			return "", fmt.Errorf("range download requires anime name and episode range (e.g., '1-5')")
		}

		animeName := strings.Join(args[:len(args)-1], " ")
		rangeStr := args[len(args)-1]

		// Parse range (e.g., "1-5")
		rangeParts := strings.Split(rangeStr, "-")
		if len(rangeParts) != 2 {
			return "", fmt.Errorf("invalid range format. Use 'start-end' (e.g., '1-5')")
		}

		startEp, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
		if err != nil {
			return "", fmt.Errorf("invalid start episode number: %s", rangeParts[0])
		}

		endEp, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
		if err != nil {
			return "", fmt.Errorf("invalid end episode number: %s", rangeParts[1])
		}

		if startEp > endEp {
			return "", fmt.Errorf("start episode (%d) cannot be greater than end episode (%d)", startEp, endEp)
		}

		if startEp < 1 {
			return "", fmt.Errorf("episode numbers must be positive")
		}

		// Store download request
		GlobalDownloadRequest = &DownloadRequest{
			AnimeName:     animeName,
			IsRange:       true,
			StartEpisode:  startEp,
			EndEpisode:    endEp,
			Source:        source,
			Quality:       quality,
		}

		return TreatingAnimeName(animeName), ErrDownloadRequested

	} else {
		// Single episode download: stardf-anime -d "anime name" episode_number
		if len(args) < 2 {
			return "", fmt.Errorf("single episode download requires anime name and episode number")
		}

		animeName := strings.Join(args[:len(args)-1], " ")
		episodeStr := args[len(args)-1]

		episodeNum, err := strconv.Atoi(episodeStr)
		if err != nil {
			return "", fmt.Errorf("invalid episode number: %s", episodeStr)
		}

		if episodeNum < 1 {
			return "", fmt.Errorf("episode number must be positive")
		}

		// Store download request
		GlobalDownloadRequest = &DownloadRequest{
			AnimeName:     animeName,
			EpisodeNum:    episodeNum,
			IsRange:       false,
			Source:        source,
			Quality:       quality,
		}

		return TreatingAnimeName(animeName), ErrDownloadRequested
	}
}
