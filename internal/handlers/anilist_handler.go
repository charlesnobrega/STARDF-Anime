package handlers

import (
	"fmt"

	"github.com/charlesnobrega/STARDF-Anime/internal/anilist"
	"github.com/charlesnobrega/STARDF-Anime/internal/util"
	"github.com/charmbracelet/huh"
)

// HandleAniListLogin starts the OAuth2 login flow interactively
func HandleAniListLogin() error {
	if anilist.GlobalSession.IsLoggedIn() {
		fmt.Println(util.SuccessStyle().Render("✅ Você já está conectado ao AniList!"))

		var logout bool
		huh.NewConfirm().
			Title("Deseja fazer logout e reconectar?").
			Value(&logout).
			Run()

		if !logout {
			return nil
		}
		_ = anilist.GlobalSession.Logout()
	}

	fmt.Println(util.BoxStyle().Render(
		util.TitleStyle().Render("🔗 Conectar AniList") + "\n\n" +
			util.MutedStyle().Render("Para sincronizar seu progresso automaticamente com o AniList,\nvocê precisa criar um Client ID gratuito em:\n") +
			util.AccentStyle().Render("https://anilist.co/settings/developer") + "\n\n" +
			util.MutedStyle().Render("Crie o app com Redirect URI: ") +
			util.InfoStyle().Render("https://anilist.co/api/v2/oauth/pin"),
	))

	var clientID string
	err := huh.NewInput().
		Title("Client ID do seu App AniList").
		Description("Cole o Client ID do app que você registrou:").
		Value(&clientID).
		Run()
	if err != nil || clientID == "" {
		util.Infof("Login cancelado.")
		return nil
	}

	cfg := anilist.OAuthConfig{
		ClientID:    clientID,
		RedirectURI: "https://anilist.co/api/v2/oauth/pin",
	}

	if err := anilist.GlobalSession.Login(cfg); err != nil {
		return fmt.Errorf("falha no login AniList: %w", err)
	}

	return nil
}

// HandleAniListLogout removes the saved AniList token
func HandleAniListLogout() error {
	if !anilist.GlobalSession.IsLoggedIn() {
		util.Infof("Você não está conectado ao AniList.")
		return nil
	}

	var confirm bool
	huh.NewConfirm().
		Title("Deseja desconectar sua conta AniList?").
		Description("Seu progresso local não será apagado.").
		Affirmative("Sim, desconectar").
		Negative("Cancelar").
		Value(&confirm).
		Run()

	if !confirm {
		return nil
	}

	if err := anilist.GlobalSession.Logout(); err != nil {
		return fmt.Errorf("erro ao fazer logout: %w", err)
	}

	util.Infof("✅ AniList desconectado com sucesso.")
	return nil
}
