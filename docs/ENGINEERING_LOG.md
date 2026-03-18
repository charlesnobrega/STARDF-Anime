# Engineering Log - Scraper Repair Mission

## Mission Overview

**Objective**: Restore functionality to Goyabu, SuperAnimes, CineGratis and AnimesOnlineCC scrapers.

---

## 🛠️ Status Tracker

### 1. Goyabu (Status: ✅ FIXED)

- **Problem**: DOM parsing failed due to dynamic content. Bot protection.
- **Solution**: 
  - Switched to regex JSON extraction (`const allEpisodes`, `var playersData`).
  - Optimized headers (removed `Accept-Encoding`).
- **Validation**: Verified. Found 5 results and extracted streams for One Piece.

### 2. AnimesOnlineCC (Status: ✅ FIXED)

- **Problem**: 403 Forbidden or empty results.
- **Solution**: 
  - Standardized headers with modern Chrome fingerprint.
  - Removed `Accept-Encoding`.
  - Updated selectors: Results (`article.item`), Episodes (`ul.episodios li`).
- **Validation**: Verified. 500 episodes found for Naruto Shippuden.

### 3. CineGratis (Status: ✅ FIXED*)

- **Problem**: Outdated selectors and basic headers.
- **Solution**: 
  - Fixed `MediaType` detection (now includes `/series-hd-online/`).
  - Updated headers to use `util.UserAgentList()`.
- **Note**: Series episode list is handled internally by the player iframe, which remains a limitation for the external selector-based list.

### 4. SuperAnimes (Status: ❌ DEPRECATED/BLOCKED)

- **Problem**: Search results return obfuscated "i" titles and 404 links even in browsers. 
- **Cause**: Sophisticated anti-bot/obfuscation system currently blocking automated search results.
- **Action**: Remaining OFFLINE in `internal/scraper/unified.go` until a bypass is found.

---

## 📅 Timeline

- **2026-03-17**: Goyabu fixed.
- **2026-03-17**: AnimesOnlineCC and CineGratis fixed. SuperAnimes search obfuscation confirmed.
- **2026-03-17**: Debug scripts organized into subdirectories to fix IDE "main redeclared" errors.
