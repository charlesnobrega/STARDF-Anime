# AI Context: STARDF-Anime

## 🚀 Como Rodar
- **Dev (TUI)**: `go run cmd/stardf-anime/main.go`
- **Build**: `go build -o stardf-anime.exe cmd/stardf-anime/main.go`
- **Testes**: `go test ./...`
- **Validação de Scrapers**: `go run debug/validate/main.go`

## 📂 Onde Mexer
- **Scrapers**: `internal/scraper/` (Adicionar/corrigir fontes)
- **Interface (TUI)**: `internal/appflow/` e subpastas de UI (Bubble Tea)
- **Modelos**: `internal/models/`
- **Player/Streaming**: `internal/player/` e `internal/playback/`
- **Mobile**: `mobile/` (Dart/Flutter)

## 🏗️ Arquitetura (10 Bullets)
1. **Linguagem**: Go (Backend/TUI) + Dart (Mobile).
2. **Framework TUI**: Charmbracelet (Bubble Tea, Lip Gloss) para uma UX rica no terminal.
3. **Scraper Manager**: Centralizado em `internal/scraper/unified.go`.
4. **Interface Unificada**: `UnifiedScraper` facilita a adição de novas fontes.
5. **Streaming**: Integração com `mpv` e `yt-dlp`.
6. **Rastreio**: Módulo `tracking` para SQLite local e possivelmente AniList.
7. **Notificações**: Sistema em `internal/notify`.
8. **Atualizações**: Auto-update via `internal/updater`.
9. **Fuzzy Search**: `go-fuzzyfinder` para seleção rápida de itens.
10. **Concorrência**: Busca paralela em múltiplos scrapers.

## ⚠️ Riscos
- **Quebra de Seletores**: Scrapers dependem de HTML externo que muda frequentemente.
- **Dependências Externas**: Requer `mpv` e `yt-dlp` instalados no PATH do sistema.
- **Múltiplos OS**: Diferenças de path e permissões entre Windows/Linux/Android.
- **Token Economy**: Consultas pesadas de scraping podem demorar e gastar tokens de contexto se logadas em excesso.
