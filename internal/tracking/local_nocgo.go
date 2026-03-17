//go:build !cgo

package tracking

// isCgoEnabled returns true because we use modernc.org/sqlite (pure Go)
func isCgoEnabled() bool {
	return true
}
