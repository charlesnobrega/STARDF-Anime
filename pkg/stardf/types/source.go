package types

import (
	"fmt"
	"strings"

	"github.com/charlesnobrega/STARDF-Anime/internal/scraper"
)

// Source represents an anime scraper source
type Source int

const (
	// SourceAnimeFire represents the AnimeFire source
	SourceAnimeFire Source = iota
	// SourceFlixHQ represents the FlixHQ source (movies/TV)
	SourceFlixHQ
	// SourceCineby represents the Cineby source (movies)
	SourceCineby
	// SourceAnimesOnlineCC represents the AnimesOnlineCC source
	SourceAnimesOnlineCC
	// SourceGoyabu represents the Goyabu source
	SourceGoyabu
)

// String returns the string representation of the source
func (s Source) String() string {
	switch s {
	case SourceAnimeFire:
		return "AnimeFire"
	case SourceFlixHQ:
		return "FlixHQ"
	case SourceCineby:
		return "Cineby"
	case SourceAnimesOnlineCC:
		return "AnimesOnlineCC"
	case SourceGoyabu:
		return "Goyabu"
	default:
		return "Unknown"
	}
}

func (s Source) ToScraperType() scraper.ScraperType {
	switch s {
	case SourceAnimeFire:
		return scraper.AnimefireType
	case SourceFlixHQ:
		return scraper.FlixHQType
	case SourceCineby:
		return scraper.CinebyType
	case SourceAnimesOnlineCC:
		return scraper.AnimesOnlineCCTYPE
	case SourceGoyabu:
		return scraper.GoyabuType
	default:
		return scraper.CinebyType
	}
}

// ParseSource parses a string into a Source type
func ParseSource(s string) (Source, error) {
	switch strings.ToLower(s) {
	case "animefire", "fire":
		return SourceAnimeFire, nil
	case "flixhq", "flix", "movies", "tv":
		return SourceFlixHQ, nil
	case "cineby", "cine":
		return SourceCineby, nil
	case "animesonlinecc", "animesonline", "cc":
		return SourceAnimesOnlineCC, nil
	case "goyabu", "goy":
		return SourceGoyabu, nil
	default:
		return SourceAnimeFire, fmt.Errorf("unknown source: %s", s)
	}
}
