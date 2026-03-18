# Engineering Log - Scraper Repair Mission

## Mission Overview

**Objective**: Restore functionality to Goyabu, SuperAnimes, and AnimesOnlineCC scrapers which were previously marked as OFFLINE or returning empty results.

---

## 🛠️ Status Tracker

### 1. Goyabu (Status: ✅ FIXED)

- **Problem**: 
  - Direct DOM parsing for episodes and streams failed because data is now dynamically injected via JavaScript.
  - Captcha/Bot protection triggered by inconsistent headers.
- **Root Cause**: The site moved episode data to a `const allEpisodes` JSON array and stream options to `var playersData` inside `<script>` tags.
- **Solution**:
  - **Regex Extraction**: Switched from `goquery` DOM parsing to regex-based JSON extraction from the page source.
  - **Header Optimization**: Removed hardcoded `Accept-Encoding: gzip, deflate, br` which was causing decompression issues in the Go client.
  - **Re-enabled**: Un-commented in `internal/scraper/unified.go`.
- **Validation**: Verified with `debug/fix_goyabu.go`. Found 5 results, 8 episodes for "One Piece: A Série", and extracted m3u8 stream URL.

### 2. SuperAnimes (Status: 🚧 IN PROGRESS)

- **Problem**: 
  - Search results are obfuscated (titles show as "i" or "I").
  - Aggressive anti-bot (ADEX captcha) detected on episode/player pages.
- **Observations**: The site seems to detect non-browser patterns quickly. DOM selectors have been updated but data returned is still garbled or blocked.
- **Next Steps**: 
  - Analyze the `i` / `I` redirection logic.
  - Test cookie persistence to see if visiting Home first bypasses protection.

### 3. AnimesOnlineCC (Status: ✅ FIXED)

- **Problem**: 
  - Basic `curl`/`http.Client` requests returned 0 results or empty bodies.
- **Root Cause**: Aggressive header checks and problematic `Accept-Encoding: gzip, deflate, br` header preventing proper decompression in the Go client.
- **Solution**:
  - **Header Synchronization**: Updated headers to match modern Chrome (`Sec-Fetch-*`, `Accept` string).
  - **Header Cleanup**: Removed the problematic `Accept-Encoding` header.
  - **Selector Update**: Implemented confirmed selectors: Results (`article.item`), Episodes (`ul.episodios li`), Video (`iframe.metaframe`).
  - **Re-enabled**: Un-commented in `internal/scraper/unified.go`.
- **Validation**: Verified with `debug/fix_animesonlinecc.go`. Found results and episodes for "Naruto Shippuden".

---

## 📝 Technical Decisions (Akita Style)

1. **Prefer Data over DOM**: When sites hide data in JS variables (JSON), we use regex extraction instead of fragile CSS selectors.
2. **Minimal Headers**: We only send essential headers (`User-Agent`, `Referer`, `Accept-Language`) to avoid triggering fingerprinting mismatches.
3. **Debug First**: Every fix is validated by a standalone script in `debug/` before being committed to the main codebase.

---

## 📅 Timeline

- **2026-03-17**: Mission started. Goyabu fixed. Documentation established.
- **2026-03-17**: AnimesOnlineCC fixed. SuperAnimes still in progress.
