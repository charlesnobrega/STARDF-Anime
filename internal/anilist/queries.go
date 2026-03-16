package anilist

import (
	"fmt"
	"time"
)

// MediaStatus represents a user's watch status on AniList
type MediaStatus string

const (
	StatusCurrent   MediaStatus = "CURRENT"
	StatusCompleted MediaStatus = "COMPLETED"
	StatusPlanning  MediaStatus = "PLANNING"
	StatusDropped   MediaStatus = "DROPPED"
	StatusPaused    MediaStatus = "PAUSED"
)

// User represents an AniList user
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MediaListEntry represents a single entry on a user's AniList
type MediaListEntry struct {
	ID              int         `json:"id"`
	MediaID         int         `json:"mediaId"`
	Status          MediaStatus `json:"status"`
	Score           float64     `json:"score"`
	Progress        int         `json:"progress"` // Episodes watched
	UpdatedAt       time.Time   `json:"updatedAt"`
	StartedAt       *FuzzyDate  `json:"startedAt"`
	CompletedAt     *FuzzyDate  `json:"completedAt"`
}

// FuzzyDate represents a potentially incomplete date (year/month/day may be 0)
type FuzzyDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

// MediaSearchResult is a partial media result from AniList search
type MediaSearchResult struct {
	ID    int    `json:"id"`
	IDMal int    `json:"idMal"`
	Title struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
		Native  string `json:"native"`
	} `json:"title"`
	CoverImage struct {
		Large  string `json:"large"`
		Medium string `json:"medium"`
	} `json:"coverImage"`
	Episodes  *int   `json:"episodes"` // pointer because can be null
	Status    string `json:"status"`
	Format    string `json:"format"` // TV, MOVIE, OVA, etc.
	AverageScore int `json:"averageScore"`
}

// graphQL queries
const querySearchMedia = `
query ($search: String, $type: MediaType) {
  Media(search: $search, type: $type) {
    id
    idMal
    title {
      romaji
      english
      native
    }
    coverImage {
      large
    }
    episodes
    status
    format
    averageScore
  }
}
`

const queryGetUserList = `
query ($userId: Int, $type: MediaType) {
  MediaListCollection(userId: $userId, type: $type) {
    lists {
      entries {
        id
        mediaId
        status
        score
        progress
        updatedAt
        media {
          id
          title {
            romaji
            english
          }
          episodes
          coverImage { large }
        }
      }
    }
  }
}
`

const mutationSaveProgress = `
mutation ($mediaId: Int, $status: MediaListStatus, $progress: Int) {
  SaveMediaListEntry(mediaId: $mediaId, status: $status, progress: $progress) {
    id
    status
    progress
  }
}
`

const queryGetViewer = `
query {
  Viewer {
    id
    name
  }
}
`

// SearchAnimeByName searches for an anime by name. No auth required.
func (c *Client) SearchAnimeByName(name string) (*MediaSearchResult, error) {
	var resp struct {
		Data struct {
			Media MediaSearchResult `json:"Media"`
		} `json:"data"`
	}

	err := c.query(querySearchMedia, map[string]interface{}{
		"search": name,
		"type":   "ANIME",
	}, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Data.Media.ID == 0 {
		return nil, nil
	}
	return &resp.Data.Media, nil
}

// SaveProgress saves watch progress to AniList. Requires authentication.
func (c *Client) SaveProgress(mediaID int, status MediaStatus, progress int) error {
	if !c.IsAuthenticated() {
		return ErrNotAuthenticated
	}

	var resp struct {
		Data struct {
			SaveMediaListEntry struct {
				ID int `json:"id"`
			} `json:"SaveMediaListEntry"`
		} `json:"data"`
	}

	return c.query(mutationSaveProgress, map[string]interface{}{
		"mediaId":  mediaID,
		"status":   string(status),
		"progress": progress,
	}, &resp)
}

// GetViewer gets the currently authenticated user. Requires authentication.
func (c *Client) GetViewer() (*User, error) {
	if !c.IsAuthenticated() {
		return nil, ErrNotAuthenticated
	}

	var resp struct {
		Data struct {
			Viewer User `json:"Viewer"`
		} `json:"data"`
	}

	err := c.query(queryGetViewer, nil, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Data.Viewer.ID == 0 {
		return nil, fmt.Errorf("could not get viewer")
	}

	return &resp.Data.Viewer, nil
}

// UserListResponse holds the complete watchlist data from AniList
type UserListResponse struct {
	Data struct {
		MediaListCollection struct {
			Lists []struct {
				Entries []struct {
					ID        int         `json:"id"`
					MediaID   int         `json:"mediaId"`
					Status    string      `json:"status"`
					Score     float64     `json:"score"`
					Progress  int         `json:"progress"`
					UpdatedAt int         `json:"updatedAt"`
					Media     struct {
						ID    int `json:"id"`
						Title struct {
							Romaji  string `json:"romaji"`
							English string `json:"english"`
						} `json:"title"`
						Episodes   *int `json:"episodes"`
						CoverImage struct {
							Large string `json:"large"`
						} `json:"coverImage"`
					} `json:"media"`
				} `json:"entries"`
			} `json:"lists"`
		} `json:"MediaListCollection"`
	} `json:"data"`
}

// GetUserList gets the user's entire watchlist from AniList.
func (c *Client) GetUserList(userID int) (*UserListResponse, error) {
	var resp UserListResponse

	err := c.query(queryGetUserList, map[string]interface{}{
		"userId": userID,
		"type":   "ANIME",
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

