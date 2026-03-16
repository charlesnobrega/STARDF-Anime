# Operational Guidelines & Lessons Learned

## 🛠 Operational Rules
1. **Clean Environment First**: Always run `taskkill /F /IM stardf-anime.exe` before any build or execution to avoid file-in-use errors.
2. **2-Minute Rule**: Never allow a command to run for more than 2 minutes without progress. If it hangs, terminate and investigate immediately.
3. **No log.Fatal**: Replace all `log.Fatalln` or `log.Fatal` with proper error returns to prevent the TUI from crashing abruptly.
4. **Safe Selection**: Always validate slices/arrays length before accessing indices from user input (Fuzzy Finder).

## 💡 Technical Solutions
- **Scraper Timeouts**: Brazilian scrapers (Goyabu, AnimesOnlineCC) must have minimal artificial delays (milliseconds, not seconds) to avoid being cut off by the `ScraperManager` early return logic.
- **Search Relevance**: Use a category selector (Anime, Movies, TV) at startup to filter sources, improving both speed and result quality.
- **Dynamic Search**: When searching, if results are missing, check if HTML selectors (classes/ids) on the target websites have changed.
