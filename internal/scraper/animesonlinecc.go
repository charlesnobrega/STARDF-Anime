package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alvarorichard/Goanime/internal/models"
	"github.com/alvarorichard/Goanime/internal/util"
)

const (
	AnimesOnlineCCReferer = "https://animesonlinecc.to"
	AnimesOnlineCCBase    = "animesonlinecc.to"
	AnimesOnlineCCAPI     = "https://api.animesonlinecc.to/api"
	UserAgent             = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0"
)

type AnimesOnlineCCClient struct {
	client  *http.Client
	referer string
}

func NewAnimesOnlineCCClient() *AnimesOnlineCCClient {
	return &AnimesOnlineCCClient{
		client:  util.GetFastClient(),
		referer: AnimesOnlineCCReferer,
	}
}

func (c *AnimesOnlineCCClient) SearchAnime(query string) ([]*models.Anime, error) {
	// Implementação similar ao AllAnime, ajustada para o site
	// Placeholder - implementar depois
	return nil, fmt.Errorf("not implemented")
}

func (c *AnimesOnlineCCClient) GetEpisodes(animeURL string) ([]models.Episode, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *AnimesOnlineCCClient) GetStreamURL(episodeURL string) (string, map[string]string, error) {
	return "", nil, fmt.Errorf("not implemented")
}