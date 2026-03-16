package anilist

import (
	"fmt"
	"os"

	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charmbracelet/huh"
)

// AniListSession wraps the client + token store for a user session
type AniListSession struct {
	Client      *Client
	TokenStore  *TokenStore
	CurrentUser *User
}

// GlobalSession is the application's AniList session
var GlobalSession = &AniListSession{
	Client:     NewClient(),
	TokenStore: NewTokenStore(),
}

// Initialize loads a saved token if available
func (s *AniListSession) Initialize() {
	token, err := s.TokenStore.Load()
	if err == nil && token != "" {
		s.Client.SetToken(token)
		util.Debug("AniList token loaded from disk")
		s.FetchViewerAsync()
	}
}

// FetchViewerAsync fetches the user profile asynchronously without blocking
func (s *AniListSession) FetchViewerAsync() {
	go func() {
		if s.IsLoggedIn() {
			user, err := s.Client.GetViewer()
			if err == nil {
				s.CurrentUser = user
			}
		}
	}()
}

// Login starts the OAuth2 login flow
// clientID must be provided (from the user's registered AniList app)
func (s *AniListSession) Login(cfg OAuthConfig) error {
	authURL := GetAuthorizationURL(cfg.ClientID, cfg.RedirectURI)

	fmt.Println(util.TitleStyle().Render("🔗 Login no AniList"))
	fmt.Println(util.InfoStyle().Render("Abra este link no seu navegador:"))
	fmt.Println(util.AccentStyle().Render(authURL))
	fmt.Println()
	fmt.Println(util.MutedStyle().Render("Após autorizar, cole o código de autorização abaixo."))

	var code string
	err := huh.NewInput().
		Title("Código de Autorização").
		Description("Cole o código que o AniList te forneceu:").
		Value(&code).
		Run()
	if err != nil {
		return err
	}

	if code == "" {
		return fmt.Errorf("código vazio, login cancelado")
	}

	token, err := ExchangeCode(cfg, code)
	if err != nil {
		return fmt.Errorf("falha ao trocar código por token: %w", err)
	}

	s.Client.SetToken(token)
	if err := s.TokenStore.Save(token); err != nil {
		util.Warnf("Não foi possível salvar o token localmente: %v", err)
	}

	user, err := s.Client.GetViewer()
	if err == nil {
		s.CurrentUser = user
		util.Infof("✅ Login no AniList realizado com sucesso como %s!", user.Name)
	} else {
		util.Infof("✅ Login no AniList realizado com sucesso!")
	}
	return nil
}

// Logout removes the saved token and clears the session
func (s *AniListSession) Logout() error {
	s.Client.SetToken("")
	return s.TokenStore.Delete()
}

// IsLoggedIn returns true if the user is authenticated
func (s *AniListSession) IsLoggedIn() bool {
	return s.Client.IsAuthenticated()
}

// SyncProgress syncs a local watchlist entry to AniList
// mediaID = AniList media ID, progress = episodes watched
func (s *AniListSession) SyncProgress(anilistID int, progress int) error {
	if !s.IsLoggedIn() {
		return ErrNotAuthenticated
	}
	return s.Client.SaveProgress(anilistID, StatusCurrent, progress)
}

// PrintStatus prints a friendly login status to the UI
func (s *AniListSession) PrintStatus() {
	if s.IsLoggedIn() {
		if s.CurrentUser != nil {
			fmt.Println(util.SuccessStyle().Render(fmt.Sprintf("✅ AniList: Conectado como %s", s.CurrentUser.Name)))
		} else {
			fmt.Println(util.SuccessStyle().Render("✅ AniList: Conectado"))
		}
	} else {
		fmt.Println(util.MutedStyle().Render("○  AniList: Não conectado (use --anilist-login para sincronizar)"))
	}
}

// init loads the token automatically on package initialization
func init() {
	// Only auto-load if not in test environment
	if os.Getenv("GO_TEST") == "" {
		GlobalSession.Initialize()
	}
}
