package scraper

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alvarorichard/Goanime/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockScraper implements UnifiedScraper for testing
type MockScraper struct {
	searchFunc      func(query string) ([]*models.Anime, error)
	episodesFunc    func(url string) ([]models.Episode, error)
	streamURLFunc   func(url string) (string, map[string]string, error)
	scraperType     ScraperType
	searchCallCount atomic.Int32
	searchDelay     time.Duration
}

func (m *MockScraper) SearchAnime(query string, options ...interface{}) ([]*models.Anime, error) {
	m.searchCallCount.Add(1)
	if m.searchDelay > 0 {
		time.Sleep(m.searchDelay)
	}
	if m.searchFunc != nil {
		return m.searchFunc(query)
	}
	return nil, nil
}

func (m *MockScraper) GetAnimeEpisodes(animeURL string) ([]models.Episode, error) {
	if m.episodesFunc != nil {
		return m.episodesFunc(animeURL)
	}
	return nil, nil
}

func (m *MockScraper) GetStreamURL(episodeURL string, options ...interface{}) (string, map[string]string, error) {
	if m.streamURLFunc != nil {
		return m.streamURLFunc(episodeURL)
	}
	return "", nil, nil
}

func (m *MockScraper) GetType() ScraperType {
	return m.scraperType
}

// createTestManager creates a ScraperManager with mock scrapers
func createTestManager(cinebyMock, animefireMock *MockScraper) *ScraperManager {
	manager := &ScraperManager{
		scrapers: make(map[ScraperType]UnifiedScraper),
	}
	if cinebyMock != nil {
		cinebyMock.scraperType = CinebyType
		manager.scrapers[CinebyType] = cinebyMock
	}
	if animefireMock != nil {
		animefireMock.scraperType = AnimefireType
		manager.scrapers[AnimefireType] = animefireMock
	}
	return manager
}

// =============================================================================
// Test: Both sources return results successfully
// =============================================================================

func TestSearchAnime_BothSourcesSucceed(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{
				{Name: "Naruto", URL: "cineby-naruto-id"},
				{Name: "Naruto Shippuden", URL: "cineby-shippuden-id"},
			}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{
				{Name: "Naruto", URL: "https://animefire.io/anime/naruto"},
				{Name: "Naruto Classico", URL: "https://animefire.io/anime/naruto-classico"},
			}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("naruto", nil)

	require.NoError(t, err)
	assert.Len(t, results, 4, "Should have results from both sources")

	// Verify both scrapers were called
	assert.Equal(t, int32(1), cinebyMock.searchCallCount.Load())
	assert.Equal(t, int32(1), animefireMock.searchCallCount.Load())

	// Verify language tags are added
	// Cineby uses [Movies/TV], AnimeFire uses [Portuguese]
	cinebyCount := 0
	animefireCount := 0
	for _, anime := range results {
		switch anime.Source {
		case "Cineby":
			cinebyCount++
			assert.Contains(t, anime.Name, "[Movies/TV]")
		case "Animefire.io":
			animefireCount++
			assert.Contains(t, anime.Name, "[Portuguese]")
		}
	}
	assert.Equal(t, 2, cinebyCount)
	assert.Equal(t, 2, animefireCount)
}

// =============================================================================
// Test: AnimeFire fails, AllAnime succeeds (Portuguese results missing)
// =============================================================================

func TestSearchAnime_AnimefireFails_CinebySucceeds(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{
				{Name: "Naruto", URL: "cineby-naruto-id"},
			}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return nil, errors.New("animefire returned a challenge page (try VPN or wait)")
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("naruto", nil)

	// Should still return results from Cineby
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Cineby", results[0].Source)

	// Both scrapers should have been called
	assert.Equal(t, int32(1), cinebyMock.searchCallCount.Load())
	assert.Equal(t, int32(1), animefireMock.searchCallCount.Load())
}

// =============================================================================
// Test: AllAnime fails, AnimeFire succeeds
// =============================================================================

func TestSearchAnime_CinebyFails_AnimefireSucceeds(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return nil, errors.New("connection timeout")
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{
				{Name: "Naruto", URL: "https://animefire.io/anime/naruto"},
			}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("naruto", nil)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Animefire.io", results[0].Source)
}

// =============================================================================
func TestSearchAnime_BothSourcesFail(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return nil, errors.New("API rate limited")
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return nil, errors.New("challenge page detected")
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("naruto", nil)

	require.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "no anime found")
	assert.Contains(t, err.Error(), "some sources failed")
}

// =============================================================================
// Test: Both sources return empty results
// =============================================================================

func TestSearchAnime_BothSourcesReturnEmpty(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	_, err := manager.SearchAnime("xyznonexistent", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no anime found")
	// Should not mention failed sources since they didn't fail
	assert.NotContains(t, err.Error(), "some sources failed")
}

// =============================================================================
// Test: One source returns empty, other returns results
// =============================================================================

func TestSearchAnime_OneSourceEmpty_OtherHasResults(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{}, nil // Vazio mas sem erro
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{
				{Name: "Anime Brasileiro", URL: "https://animefire.io/anime/brasileiro"},
			}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("brasileiro", nil)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Animefire.io", results[0].Source)
}

// =============================================================================
// Test: Concurrent execution - both scrapers run in parallel
// =============================================================================

func TestSearchAnime_ConcurrentExecution(t *testing.T) {
	t.Parallel()

	var cinebyStart, animefireStart time.Time
	var mu sync.Mutex

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			mu.Lock()
			cinebyStart = time.Now()
			mu.Unlock()

			time.Sleep(100 * time.Millisecond)

			return []*models.Anime{{Name: "Cineby Result", URL: "id1"}}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			mu.Lock()
			animefireStart = time.Now()
			mu.Unlock()

			time.Sleep(100 * time.Millisecond)

			return []*models.Anime{{Name: "AnimeFire Result", URL: "https://animefire.io/1"}}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)

	start := time.Now()
	results, err := manager.SearchAnime("test", nil)
	totalDuration := time.Since(start)

	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Se estiver rodando concorrentemente, o tempo total deve ser ~100ms, não ~200ms
	assert.Less(t, totalDuration, 180*time.Millisecond,
		"As buscas devem rodar concorrentemente, não sequencialmente")

	// Verifica se ambos começaram por volta do mesmo tempo
	mu.Lock()
	startDiff := cinebyStart.Sub(animefireStart)
	if startDiff < 0 {
		startDiff = -startDiff
	}
	mu.Unlock()

	assert.Less(t, startDiff, 50*time.Millisecond,
		"Ambas as buscas devem começar quase simultaneamente")
}

// =============================================================================
// Test: Slow source doesn't block fast source results
// =============================================================================

func TestSearchAnime_SlowSourceDoesNotBlockFastSource(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			time.Sleep(200 * time.Millisecond) // Lento
			return []*models.Anime{{Name: "Slow Result", URL: "id1"}}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			time.Sleep(10 * time.Millisecond) // Rápido
			return []*models.Anime{{Name: "Fast Result", URL: "https://animefire.io/1"}}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("test", nil)

	require.NoError(t, err)
	// Ambos os resultados devem estar presentes
	assert.Len(t, results, 2)
}

// =============================================================================
// Test: Specific scraper selection - AnimeFire only
// =============================================================================

func TestSearchAnime_SpecificScraper_AnimefireOnly(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{{Name: "Cineby Result", URL: "id1"}}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{{Name: "AnimeFire Result", URL: "https://animefire.io/1"}}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)

	scraperType := AnimefireType
	results, err := manager.SearchAnime("test", &scraperType)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Animefire.io", results[0].Source)

	// Apenas AnimeFire deve ser chamado
	assert.Equal(t, int32(0), cinebyMock.searchCallCount.Load())
	assert.Equal(t, int32(1), animefireMock.searchCallCount.Load())
}

// =============================================================================
// Test: Specific scraper selection - AllAnime only
// =============================================================================

func TestSearchAnime_SpecificScraper_CinebyOnly(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{{Name: "Cineby Result", URL: "id1"}}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{{Name: "AnimeFire Result", URL: "https://animefire.io/1"}}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)

	scraperType := CinebyType
	results, err := manager.SearchAnime("test", &scraperType)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Cineby", results[0].Source)

	// Only Cineby should be called
	assert.Equal(t, int32(1), cinebyMock.searchCallCount.Load())
	assert.Equal(t, int32(0), animefireMock.searchCallCount.Load())
}

// =============================================================================
// Test: Specific scraper fails - returns error
// =============================================================================

func TestSearchAnime_SpecificScraper_Fails(t *testing.T) {
	t.Parallel()

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return nil, errors.New("Cloudflare challenge")
		},
	}

	manager := createTestManager(nil, animefireMock)

	scraperType := AnimefireType
	results, err := manager.SearchAnime("test", &scraperType)

	require.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "busca falhou")
	assert.Contains(t, err.Error(), "Cloudflare challenge")
}

// =============================================================================
// Test: Source tags are not duplicated
// =============================================================================

func TestSearchAnime_SourceTagsNotDuplicated(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{
				{Name: "[Movies/TV] Naruto", URL: "id1"}, // Já tem tag
			}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{
				{Name: "[Portuguese] Naruto", URL: "https://animefire.io/1"}, // Já tem tag
			}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("naruto", nil)

	require.NoError(t, err)

	for _, anime := range results {
		// Conta ocorrências de tags
		cinebyTagCount := countOccurrences(anime.Name, "[Movies/TV]")
		animefireTagCount := countOccurrences(anime.Name, "[Portuguese]")

		// Nunca deve ter mais de uma de cada tag
		assert.LessOrEqual(t, cinebyTagCount, 1, "Cineby tag duplicada")
		assert.LessOrEqual(t, animefireTagCount, 1, "AnimeFire tag duplicada")
	}
}

// =============================================================================
// Test: Race condition - multiple concurrent searches
// =============================================================================

func TestSearchAnime_NoConcurrentRaceConditions(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			time.Sleep(10 * time.Millisecond)
			return []*models.Anime{{Name: "Result " + query, URL: "id-" + query}}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			time.Sleep(10 * time.Millisecond)
			return []*models.Anime{{Name: "AF Result " + query, URL: "https://animefire.io/" + query}}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)

	// Run multiple concurrent searches
	var wg sync.WaitGroup
	errChan := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results, err := manager.SearchAnime("query", nil)
			if err != nil {
				errChan <- err
				return
			}
			if len(results) != 2 {
				errChan <- errors.New("unexpected result count")
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("Concurrent search error: %v", err)
	}
}

// =============================================================================
// Test: Network timeout simulation
// =============================================================================

func TestSearchAnime_NetworkTimeout(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{{Name: "Quick Result", URL: "id1"}}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			// Simula erro de timeout de rede
			return nil, errors.New("connection timeout after 30s")
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("test", nil)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Cineby", results[0].Source)
}

// =============================================================================
// Test: VPN required error from AnimeFire
// =============================================================================

func TestSearchAnime_VPNRequired(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{{Name: "English Result", URL: "id1"}}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return nil, errors.New("access restricted: VPN may be required")
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("test", nil)

	require.NoError(t, err, "Should return results from working source")
	assert.Len(t, results, 1)
}

// =============================================================================
// Test: Cloudflare challenge detection
// =============================================================================

func TestSearchAnime_CloudflareChallenge(t *testing.T) {
	t.Parallel()

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return []*models.Anime{{Name: "Result", URL: "id1"}}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			return nil, errors.New("animefire returned a challenge page (try VPN or wait)")
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	results, err := manager.SearchAnime("test", nil)

	require.NoError(t, err)
	assert.Len(t, results, 1)
}

// =============================================================================
// Test: Query is passed correctly to scrapers
// =============================================================================

func TestSearchAnime_QueryPassedCorrectly(t *testing.T) {
	t.Parallel()

	var capturedQueries []string
	var mu sync.Mutex

	cinebyMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			mu.Lock()
			capturedQueries = append(capturedQueries, "cineby:"+query)
			mu.Unlock()
			return []*models.Anime{}, nil
		},
	}

	animefireMock := &MockScraper{
		searchFunc: func(query string) ([]*models.Anime, error) {
			mu.Lock()
			capturedQueries = append(capturedQueries, "animefire:"+query)
			mu.Unlock()
			return []*models.Anime{}, nil
		},
	}

	manager := createTestManager(cinebyMock, animefireMock)
	_, _ = manager.SearchAnime("Shingeki no Kyojin", nil)

	mu.Lock()
	defer mu.Unlock()

	assert.Len(t, capturedQueries, 2)
	assert.Contains(t, capturedQueries, "cineby:Shingeki no Kyojin")
	assert.Contains(t, capturedQueries, "animefire:Shingeki no Kyojin")
}

// =============================================================================
// Helper functions
// =============================================================================

func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			count++
		}
	}
	return count
}
