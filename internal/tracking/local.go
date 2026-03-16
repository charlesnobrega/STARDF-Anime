// tracking/sqlite_tracker.go
package tracking

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// IsCgoEnabled indicates whether CGO is enabled for SQLite support
var IsCgoEnabled = true

// Error constants
var (
	ErrCgoDisabled      = errors.New("CGO disabled: sqlite tracking not available")
	ErrTrackerNotInited = errors.New("tracker not initialized")
)

/*
────────────────────────────────────────────────────────────────────────────*
│  Constantes de Configuração                                                │
*────────────────────────────────────────────────────────────────────────────
*/
const (
	defaultCacheSize  = -32000    // 32MB (increased from 20MB for better performance)
	mmapSize          = 536870912 // 512MB (increased from 256MB)
	busyTimeout       = 3000      // 3 seconds (reduced from 5s for faster response)
	walAutoCheckpoint = 500       // pages (reduced for more frequent checkpoints)
	maxOpenConns      = 3         // reduced for lower overhead
	maxIdleConns      = 2         // keep idle connections
	avgAnimePerUser   = 100       // pre-allocation for slices
)

/*
────────────────────────────────────────────────────────────────────────────*
│  Tipos e Estruturas                                                        │
*────────────────────────────────────────────────────────────────────────────
*/

// Anime represents tracked media (anime, movie, or TV show)
type Anime struct {
	AnilistID     int       `json:"anilist_id"`
	AllanimeID    string    `json:"allanime_id"` // Primary key - unique per content
	EpisodeNumber int       `json:"episode_number"`
	PlaybackTime  int       `json:"playback_time"`
	Duration      int       `json:"duration"`
	Title         string    `json:"title"`
	MediaType     string    `json:"media_type"` // "anime" or "movie"
	LastUpdated   time.Time `json:"last_updated"`
}

// FollowedMedia represents a media item the user is following
type FollowedMedia struct {
	AllanimeID    string    `json:"allanime_id"`
	AnilistID     int       `json:"anilist_id"`
	Title         string    `json:"title"`
	MediaType     string    `json:"media_type"`
	TotalEpisodes int       `json:"total_episodes"`
	Status        string    `json:"status"` // watching, planning, completed, dropped
	Score         int       `json:"score"`  // User rating 1-10
	LastChecked   time.Time `json:"last_checked"`
	LastEpisode   int       `json:"last_episode"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ScraperHealth represents the reliability of a scraper plugin
type ScraperHealth struct {
	ScraperName    string    `json:"scraper_name"`
	TotalSearches  int       `json:"total_searches"`
	FailedSearches int       `json:"failed_searches"`
	LastFailure    string    `json:"last_failure"`
	LastUpdate     time.Time `json:"last_update"`
}

type LocalTracker struct {
	db              *sql.DB
	upsertPS        *sql.Stmt
	getPS           *sql.Stmt
	allPS           *sql.Stmt
	deletePS        *sql.Stmt
	followUpsertPS  *sql.Stmt
	followGetPS     *sql.Stmt
	followAllPS     *sql.Stmt
	followDeletePS  *sql.Stmt
	healthUpdatePS  *sql.Stmt
	healthAllPS     *sql.Stmt
}

/*
────────────────────────────────────────────────────────────────────────────*
│  Singleton/Cache Global do Tracker                                         │
*────────────────────────────────────────────────────────────────────────────
*/
var (
	globalTracker     *LocalTracker
	globalTrackerPath string
	trackerMutex      = &sync.Mutex{}
)

// GetGlobalTracker returns the cached global tracker instance.
// This avoids repeatedly opening the database connection which is slow.
func GetGlobalTracker() *LocalTracker {
	trackerMutex.Lock()
	defer trackerMutex.Unlock()
	return globalTracker
}

// CloseGlobalTracker closes the global tracker and clears the cache.
// Should be called on application shutdown.
func CloseGlobalTracker() error {
	trackerMutex.Lock()
	defer trackerMutex.Unlock()

	if globalTracker != nil {
		err := globalTracker.Close()
		globalTracker = nil
		globalTrackerPath = ""
		return err
	}
	return nil
}

/*
────────────────────────────────────────────────────────────────────────────*
│  Construtor e Inicialização                                                │
*────────────────────────────────────────────────────────────────────────────
*/
var NewLocalTracker func(dbPath string) *LocalTracker

func newLocalTrackerImpl(dbPath string) *LocalTracker {
	// Use singleton pattern to avoid repeatedly opening the database
	trackerMutex.Lock()
	defer trackerMutex.Unlock()

	// Return cached tracker if path matches
	if globalTracker != nil && globalTrackerPath == dbPath {
		return globalTracker
	}
	// Check if CGO is disabled (SQLite not available)
	if !IsCgoEnabled {
		fmt.Println("Warning: CGO is disabled, anime progress tracking will be unavailable")
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0700); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		return nil
	}
	// Build DSN with optimized pragmas for faster operations
	var dsn string
	if runtime.GOOS == "windows" {
		// Use URI format with escape for Windows
		escapedPath := strings.ReplaceAll(dbPath, "\\", "/")
		dsn = fmt.Sprintf(
			"file:%s?_journal_mode=WAL&_synchronous=OFF&_wal_autocheckpoint=%d&"+
				"_busy_timeout=%d&_cache_size=%d&_mmap_size=%d&_mode=rwc&_temp_store=MEMORY",
			escapedPath,
			walAutoCheckpoint,
			busyTimeout,
			defaultCacheSize,
			mmapSize,
		)
	} else {
		dsn = fmt.Sprintf(
			"file:%s?_journal_mode=WAL&_synchronous=OFF&_wal_autocheckpoint=%d&"+
				"_busy_timeout=%d&_cache_size=%d&_mmap_size=%d&_temp_store=MEMORY",
			dbPath,
			walAutoCheckpoint,
			busyTimeout,
			defaultCacheSize,
			mmapSize,
		)
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return nil
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(0) // Keep connections alive indefinitely

	if err := initializeDatabase(db); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Printf("Error closing database: %v\n", closeErr)
		}
		fmt.Printf("Error initializing database: %v\n", err)
		return nil
	}

	statements, err := prepareStatements(db)
	if err != nil {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Printf("Error closing database: %v\n", closeErr)
		}
		fmt.Printf("Error preparing statements: %v\n", err)
		return nil
	}

	tracker := &LocalTracker{
		db:              db,
		upsertPS:        statements.upsert,
		getPS:           statements.get,
		allPS:           statements.all,
		deletePS:        statements.delete,
		followUpsertPS:  statements.followUpsert,
		followGetPS:     statements.followGet,
		followAllPS:     statements.followAll,
		followDeletePS:  statements.followDelete,
		healthUpdatePS:  statements.healthUpdate,
		healthAllPS:     statements.healthAll,
	}

	// Cache the tracker globally for reuse
	globalTracker = tracker
	globalTrackerPath = dbPath

	return tracker
}

/*
────────────────────────────────────────────────────────────────────────────*
│  Inicialização do Banco de Dados                                           │
*────────────────────────────────────────────────────────────────────────────
*/
func initializeDatabase(db *sql.DB) error {
	// New schema: allanime_id is the primary key (works for both movies and anime)
	// media_type: 'anime' or 'movie' to distinguish content type
	schema := `
	CREATE TABLE IF NOT EXISTS media_progress (
		allanime_id    TEXT    PRIMARY KEY NOT NULL,
		anilist_id     INTEGER DEFAULT 0,
		episode_number INTEGER NOT NULL,
		playback_time  INTEGER NOT NULL CHECK(playback_time >= 0),
		duration       INTEGER NOT NULL CHECK(duration > 0),
		title          TEXT,
		media_type     TEXT    DEFAULT 'anime',
		last_updated   INTEGER NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS followed_media (
		allanime_id    TEXT    PRIMARY KEY NOT NULL,
		anilist_id     INTEGER DEFAULT 0,
		title          TEXT,
		media_type     TEXT    DEFAULT 'anime',
		total_episodes INTEGER DEFAULT 0,
		status         TEXT    DEFAULT 'planning',
		score          INTEGER DEFAULT 0,
		last_checked   INTEGER NOT NULL,
		last_episode   INTEGER DEFAULT 0,
		updated_at     INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS scraper_health (
		scraper_name    TEXT    PRIMARY KEY NOT NULL,
		total_searches  INTEGER DEFAULT 0,
		failed_searches INTEGER DEFAULT 0,
		last_failure    TEXT,
		last_update     INTEGER NOT NULL
	);`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("schema creation failed: %w", err)
	}

	// Migrate old data if anime_progress table exists
	migrateOldData(db)

	// Add score column if it doesn't exist (migration for older v1.6.4 clients)
	_, _ = db.Exec("ALTER TABLE followed_media ADD COLUMN score INTEGER DEFAULT 0")

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_media_lookup 
		ON media_progress(allanime_id, last_updated DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_media_type 
		ON media_progress(media_type, last_updated DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_followed_status
		ON followed_media(status, updated_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_health_failure
		ON scraper_health(failed_searches DESC)`,
	}

	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			return fmt.Errorf("index creation '%s' failed: %w", idx, err)
		}
	}

	if _, err := db.Exec(`PRAGMA optimize`); err != nil {
		return fmt.Errorf("initial optimization failed: %w", err)
	}

	return nil
}

// migrateOldData migrates data from old anime_progress table to new media_progress table
func migrateOldData(db *sql.DB) {
	// Check if old table exists
	var tableName string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='anime_progress'").Scan(&tableName)
	if err != nil {
		return // Old table doesn't exist, nothing to migrate
	}

	// Migrate data: for each allanime_id, keep the entry with highest playback_time
	_, err = db.Exec(`
		INSERT OR REPLACE INTO media_progress (allanime_id, anilist_id, episode_number, playback_time, duration, title, media_type, last_updated)
		SELECT 
			allanime_id,
			MAX(anilist_id),
			episode_number,
			MAX(playback_time),
			MAX(duration),
			title,
			CASE WHEN title LIKE '[Movies/TV]%' OR title LIKE '[Movie]%' THEN 'movie' ELSE 'anime' END,
			MAX(last_updated)
		FROM anime_progress
		GROUP BY allanime_id
	`)
	if err != nil {
		fmt.Printf("Warning: migration failed: %v\n", err)
		return
	}

	// Drop old table after successful migration
	_, _ = db.Exec("DROP TABLE IF EXISTS anime_progress")
	fmt.Println("✓ Migrated tracking data to new format (movies + anime support)")
}

/*
────────────────────────────────────────────────────────────────────────────*
│  Preparação de Statements                                                  │
*────────────────────────────────────────────────────────────────────────────
*/
type preparedStatements struct {
	upsert       *sql.Stmt
	get          *sql.Stmt
	all          *sql.Stmt
	delete       *sql.Stmt
	followUpsert *sql.Stmt
	followGet    *sql.Stmt
	followAll    *sql.Stmt
	followDelete *sql.Stmt
	healthUpdate *sql.Stmt
	healthAll    *sql.Stmt
}

func prepareStatements(db *sql.DB) (*preparedStatements, error) {
	// New schema: allanime_id is the primary key (unique per content)
	upsert, err := db.Prepare(`INSERT INTO media_progress (
		allanime_id,
		anilist_id, 
		episode_number, 
		playback_time, 
		duration, 
		title,
		media_type,
		last_updated
	) VALUES (?,?,?,?,?,?,?,?) 
	ON CONFLICT(allanime_id) DO UPDATE SET
		anilist_id = CASE WHEN excluded.anilist_id > 0 THEN excluded.anilist_id ELSE media_progress.anilist_id END,
		episode_number = excluded.episode_number,
		playback_time = excluded.playback_time,
		duration = excluded.duration,
		title = excluded.title,
		media_type = excluded.media_type,
		last_updated = excluded.last_updated`)

	if err != nil {
		return nil, fmt.Errorf("upsert preparation failed: %w", err)
	}

	// Get by allanime_id only (works for both movies and anime)
	get, err := db.Prepare(`SELECT 
		episode_number, 
		playback_time, 
		duration, 
		title, 
		last_updated 
	FROM media_progress 
	WHERE allanime_id = ?`)

	if err != nil {
		return nil, fmt.Errorf("get preparation failed: %w", err)
	}

	all, err := db.Prepare(`SELECT 
		anilist_id, 
		allanime_id, 
		episode_number, 
		playback_time, 
		duration, 
		title, 
		last_updated 
	FROM media_progress`)

	if err != nil {
		return nil, fmt.Errorf("all preparation failed: %w", err)
	}

	delete, err := db.Prepare(`DELETE FROM media_progress 
		WHERE allanime_id = ?`)

	if err != nil {
		return nil, fmt.Errorf("delete preparation failed: %w", err)
	}

	// Followed Media Statements
	followUpsert, err := db.Prepare(`INSERT INTO followed_media (
		allanime_id, anilist_id, title, media_type, total_episodes, status, score, last_checked, last_episode, updated_at
	) VALUES (?,?,?,?,?,?,?,?,?,?)
	ON CONFLICT(allanime_id) DO UPDATE SET
		anilist_id = CASE WHEN excluded.anilist_id > 0 THEN excluded.anilist_id ELSE followed_media.anilist_id END,
		title = excluded.title,
		media_type = excluded.media_type,
		total_episodes = excluded.total_episodes,
		status = excluded.status,
		score = excluded.score,
		last_checked = excluded.last_checked,
		last_episode = CASE WHEN excluded.last_episode > 0 THEN excluded.last_episode ELSE followed_media.last_episode END,
		updated_at = excluded.updated_at`)
	if err != nil {
		return nil, fmt.Errorf("followUpsert preparation failed: %w", err)
	}

	followGet, err := db.Prepare(`SELECT anilist_id, title, media_type, total_episodes, status, score, last_checked, last_episode, updated_at 
		FROM followed_media WHERE allanime_id = ?`)
	if err != nil {
		return nil, fmt.Errorf("followGet preparation failed: %w", err)
	}

	followAll, err := db.Prepare(`SELECT allanime_id, anilist_id, title, media_type, total_episodes, status, score, last_checked, last_episode, updated_at 
		FROM followed_media ORDER BY updated_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("followAll preparation failed: %w", err)
	}

	followDelete, err := db.Prepare(`DELETE FROM followed_media WHERE allanime_id = ?`)
	if err != nil {
		return nil, fmt.Errorf("followDelete preparation failed: %w", err)
	}

	// Scraper Health Statements
	healthUpdate, err := db.Prepare(`INSERT INTO scraper_health (
		scraper_name, total_searches, failed_searches, last_failure, last_update
	) VALUES (?,1,?,?,?)
	ON CONFLICT(scraper_name) DO UPDATE SET
		total_searches = scraper_health.total_searches + 1,
		failed_searches = scraper_health.failed_searches + ?,
		last_failure = CASE WHEN ? != '' THEN ? ELSE scraper_health.last_failure END,
		last_update = excluded.last_update`)
	if err != nil {
		return nil, fmt.Errorf("healthUpdate preparation failed: %w", err)
	}

	healthAll, err := db.Prepare(`SELECT scraper_name, total_searches, failed_searches, last_failure, last_update 
		FROM scraper_health ORDER BY failed_searches DESC`)
	if err != nil {
		return nil, fmt.Errorf("healthAll preparation failed: %w", err)
	}

	return &preparedStatements{
		upsert:       upsert,
		get:          get,
		all:          all,
		delete:       delete,
		followUpsert: followUpsert,
		followGet:    followGet,
		followAll:    followAll,
		followDelete: followDelete,
		healthUpdate: healthUpdate,
		healthAll:    healthAll,
	}, nil
}

/*
────────────────────────────────────────────────────────────────────────────*
│  Operações Principais                                                      │
*────────────────────────────────────────────────────────────────────────────
*/
func (t *LocalTracker) UpdateProgress(a Anime) error {
	// Safety check for when tracker is not initialized
	if t == nil || t.db == nil || t.upsertPS == nil {
		return ErrTrackerNotInited
	}

	// Validate duration to prevent constraint errors
	if a.Duration <= 0 {
		return fmt.Errorf("invalid duration value (%d): must be greater than 0", a.Duration)
	}

	// Validate playback time (shouldn't be negative)
	if a.PlaybackTime < 0 {
		a.PlaybackTime = 0
	}

	// Determine media type from title if not set
	if a.MediaType == "" {
		if strings.Contains(a.Title, "[Movies/TV]") || strings.Contains(a.Title, "[Movie]") {
			a.MediaType = "movie"
		} else {
			a.MediaType = "anime"
		}
	}

	// New order: allanime_id first (primary key)
	_, err := t.upsertPS.Exec(
		a.AllanimeID,
		a.AnilistID,
		a.EpisodeNumber,
		a.PlaybackTime,
		a.Duration,
		a.Title,
		a.MediaType,
		a.LastUpdated.Unix(),
	)
	return err
}

// GetAnime retrieves tracking data by allanime_id (works for both movies and anime)
// The anilistID parameter is kept for backwards compatibility but is ignored
func (t *LocalTracker) GetAnime(anilistID int, allanimeID string) (*Anime, error) {
	// Safety check for when tracker is not initialized
	if t == nil || t.db == nil || t.getPS == nil {
		return nil, ErrTrackerNotInited
	}

	var a Anime
	var ts int64

	// Query by allanime_id only (primary key)
	err := t.getPS.QueryRow(allanimeID).Scan(
		&a.EpisodeNumber,
		&a.PlaybackTime,
		&a.Duration,
		&a.Title,
		&ts,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	a.AnilistID = anilistID
	a.AllanimeID = allanimeID
	a.LastUpdated = time.Unix(ts, 0)
	return &a, nil
}

func (t *LocalTracker) GetAllAnime() ([]Anime, error) {
	// Safety check for when tracker is not initialized
	if t == nil || t.db == nil || t.allPS == nil {
		return nil, ErrTrackerNotInited
	}

	rows, err := t.allPS.Query()
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	list := make([]Anime, 0, avgAnimePerUser)
	for rows.Next() {
		var a Anime
		var ts int64
		if err := rows.Scan(
			&a.AnilistID,
			&a.AllanimeID,
			&a.EpisodeNumber,
			&a.PlaybackTime,
			&a.Duration,
			&a.Title,
			&ts,
		); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		a.LastUpdated = time.Unix(ts, 0)
		list = append(list, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return list, nil
}

// DeleteAnime removes tracking data by allanime_id
// The anilistID parameter is kept for backwards compatibility but is ignored
func (t *LocalTracker) DeleteAnime(anilistID int, allanimeID string) error {
	_, err := t.deletePS.Exec(allanimeID)
	return err
}

// UpdateFollowedMedia adds or updates a media in the watchlist
func (t *LocalTracker) UpdateFollowedMedia(m FollowedMedia) error {
	if t == nil || t.followUpsertPS == nil {
		return ErrTrackerNotInited
	}
	_, err := t.followUpsertPS.Exec(
		m.AllanimeID,
		m.AnilistID,
		m.Title,
		m.MediaType,
		m.TotalEpisodes,
		m.Status,
		m.Score,
		m.LastChecked.Unix(),
		m.LastEpisode,
		m.UpdatedAt.Unix(),
	)
	return err
}

// GetFollowedMedia retrieves a single followed media
func (t *LocalTracker) GetFollowedMedia(allanimeID string) (*FollowedMedia, error) {
	if t == nil || t.followGetPS == nil {
		return nil, ErrTrackerNotInited
	}
	var m FollowedMedia
	var lc, ua int64
	err := t.followGetPS.QueryRow(allanimeID).Scan(
		&m.AnilistID,
		&m.Title,
		&m.MediaType,
		&m.TotalEpisodes,
		&m.Status,
		&m.Score,
		&lc,
		&m.LastEpisode,
		&ua,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	m.AllanimeID = allanimeID
	m.LastChecked = time.Unix(lc, 0)
	m.UpdatedAt = time.Unix(ua, 0)
	return &m, nil
}

// GetAllFollowedMedia retrieves all items in the watchlist
func (t *LocalTracker) GetAllFollowedMedia() ([]FollowedMedia, error) {
	if t == nil || t.followAllPS == nil {
		return nil, ErrTrackerNotInited
	}
	rows, err := t.followAllPS.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []FollowedMedia
	for rows.Next() {
		var m FollowedMedia
		var lc, ua int64
		if err := rows.Scan(
			&m.AllanimeID,
			&m.AnilistID,
			&m.Title,
			&m.MediaType,
			&m.TotalEpisodes,
			&m.Status,
			&m.Score,
			&lc,
			&m.LastEpisode,
			&ua,
		); err != nil {
			return nil, err
		}
		m.LastChecked = time.Unix(lc, 0)
		m.UpdatedAt = time.Unix(ua, 0)
		list = append(list, m)
	}
	return list, nil
}

// UnfollowMedia removes an item from the watchlist
func (t *LocalTracker) UnfollowMedia(allanimeID string) error {
	if t == nil || t.followDeletePS == nil {
		return ErrTrackerNotInited
	}
	_, err := t.followDeletePS.Exec(allanimeID)
	return err
}

// TrackScraperAction records a success or failure for a scraper
func (t *LocalTracker) TrackScraperAction(name string, success bool, lastErr string) error {
	if t == nil || t.healthUpdatePS == nil {
		return ErrTrackerNotInited
	}
	failInc := 0
	if !success {
		failInc = 1
	}
	now := time.Now().Unix()
	_, err := t.healthUpdatePS.Exec(
		name,
		lastErr,
		now,
		failInc,
		lastErr,
		lastErr,
	)
	return err
}

// GetScraperHealthRecords retrieves all scraper health data
func (t *LocalTracker) GetScraperHealthRecords() ([]ScraperHealth, error) {
	if t == nil || t.healthAllPS == nil {
		return nil, ErrTrackerNotInited
	}
	rows, err := t.healthAllPS.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ScraperHealth
	for rows.Next() {
		var h ScraperHealth
		var lu int64
		if err := rows.Scan(
			&h.ScraperName,
			&h.TotalSearches,
			&h.FailedSearches,
			&h.LastFailure,
			&lu,
		); err != nil {
			return nil, err
		}
		h.LastUpdate = time.Unix(lu, 0)
		list = append(list, h)
	}
	return list, nil
}

/*
────────────────────────────────────────────────────────────────────────────*
│  Finalização                                                               │
*────────────────────────────────────────────────────────────────────────────
*/
func (t *LocalTracker) Close() error {
	if t == nil {
		return nil
	}

	var finalErr error

	closeStmt := func(stmt *sql.Stmt, name string) {
		if stmt != nil {
			if err := stmt.Close(); err != nil {
				finalErr = fmt.Errorf("%s statement close error: %w", name, err)
			}
		}
	}

	closeStmt(t.upsertPS, "upsert")
	closeStmt(t.getPS, "get")
	closeStmt(t.allPS, "all")
	closeStmt(t.deletePS, "delete")
	closeStmt(t.followUpsertPS, "followUpsert")
	closeStmt(t.followGetPS, "followGet")
	closeStmt(t.followAllPS, "followAll")
	closeStmt(t.followDeletePS, "followDelete")
	closeStmt(t.healthUpdatePS, "healthUpdate")
	closeStmt(t.healthAllPS, "healthAll")

	if err := t.db.Close(); err != nil {
		finalErr = fmt.Errorf("database close error: %w", err)
	}

	return finalErr
}

func init() {
	// This will be replaced at build time with false if CGO is disabled
	// When using CGO_ENABLED=0
	IsCgoEnabled = isCgoEnabled()
	NewLocalTracker = newLocalTrackerImpl // Initialize the public variable
}

// The implementation of isCgoEnabled is defined in local_cgo.go and local_nocgo.go
// based on build tags. We don't define it here to avoid duplicate declarations.
