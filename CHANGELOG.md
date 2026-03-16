# StarDF-Anime Release Notes - Version 1.6.3

Release date: 2026-03-16

## Highlights

- **Rebranding**: The project has been officially renamed from GoAnime to **StarDF-Anime**.
- **New Identity**: Total migration of namespace to `charlesnobrega/STARDF-Anime`.
- **Improved Stability**: Fixed compilation errors in WordPress scrapers and updated internal components to reflect the new identity.

## Features

- Updated all internal logs and help menus to reflect the new **StarDF-Anime** branding.
- Migrated all repository links and references in the documentation.
- Bumped version to `v1.6.3`.

## Bug Fixes

- Fixed "rand undefined" compilation errors in `Goyabu`, `SuperAnimes`, and `AnimesOnlineCC` scrapers.
- Updated the auto-updater to correctly point to the new repository location at `charlesnobrega/STARDF-Anime`.
- Resolved file locking issues during local builds on Windows.

## Scraper Status

- **AnimeFire**: Fully functional.
- **FlixHQ**: Fully functional (Movies/TV).
- **Goyabu/SuperAnimes**: Currently OFFLINE/Unstable (re-evaluated and kept disabled for stability).

---

# StarDF-Anime Release Notes - Version 1.6.2

Release date: 2026-01-19

## Highlights

- **SQLite Local Tracking Enabled**: All release binaries are now compiled with CGO enabled, providing full SQLite-based local tracking support for watch history and progress.

## Features

- All platform binaries (Linux, macOS, Windows) now include SQLite local tracking support.
- Native builds for each platform ensure optimal performance and compatibility.

## Improvements

- Release workflow updated to use native runners for each platform (ubuntu, macos, windows) for CGO support.
- Improved debug logging in the auto-updater to show available release assets.

## Bug Fixes

- Fixed release workflow to avoid duplicate asset uploads.
- Fixed AUR publish workflow secrets check (moved from job-level to step-level).
- Fixed updater debug output showing available assets for troubleshooting.

---

# GoAnime Release Notes - Version 1.6

Release date: 2026-01-18

## Features


- FlixHQ integration for movies and TV shows, enabling searching and playback of FlixHQ content.
- TMDB and OMDb integration for improved media enrichment and metadata (posters, descriptions, external IDs).
- Concurrent anime search with exponential backoff for faster, more reliable search results across sources.
- Episode data providers with fallback support to improve episode lookup resilience.
- Enhanced playback features for movies and TV shows, including HLS stream handling and better MPV integration.
- Improved Discord Rich Presence: shows clean title and precise timing, removing media-type tags from titles.
- Fuzzy server selection for AnimeDrive video options.
- Restored selection option for episode and anime in playback menus.
- Added AUR package support and publishing workflow for Arch Linux users.
- Automated release workflow via GitHub Actions to streamline builds and releases.

## Improvements

- General search and playback experience improvements and UI text consistency (error messages/prompts now in English).
- Updated AnimeFire source references and other scraper consistency fixes.
- Improved title sanitization and retrieval logic to avoid noisy tags in titles.
- Windows installer improvements: better configuration handling, MPV DLL inclusion, and PATH improvements.
- CI and release workflow restructuring: binary builds, artifact upload, and RELEASE_NOTES.md support.
- Dependency updates across the codebase for improved stability and performance.
- Code formatting and readability improvements across multiple files.

## Bug Fixes

- Fixed Discord invite link in README files.
- Corrected base URL in anime parser tests and other test fixes.
- Fixed formatting and path-detection issues in platform-specific MPV helper functions.
- Fixed PKGBUILD source URL and added optional dependencies for packaging.

- AnimeDrive: integration worked for several days but the source enabled Cloudflare protections; AnimeDrive support is temporarily on standby (integration commented/disabled) until a reliable, compliant workaround is found.

## Developer Notes

- CI: removed ARM64 Windows build from the release workflow and added AUR publishing workflow.
- Added tests for AnimeDrive client, `CleanTitle`, and search variation generation.
- Continued refactors to streamline PATH handling and improve test coverage.
- Many small refactors, formatting (go fmt), and chore updates to keep dependencies current.

---

For upgrade instructions and binary downloads, see the project releases and the updated release workflow in the repository.

Thank you to all contributors for this release.

