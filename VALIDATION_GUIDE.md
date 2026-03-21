# 🛡️ Guia de Validação StarDF-Anime

Este documento é obrigatório. Toda modificação no código deve ser validada contra este checklist para evitar regressões.

## 🔍 1. Fluxo de Busca (Search)
- [ ] **Fuzzy Match**: O título do AniList deve encontrar correspondência no scraper (ex: Jujutsu S3 -> Shimetsu Kaiyuu).
- [ ] **Unificação**: Resultados de múltiplos sites para o mesmo anime devem aparecer em um ÚNICO card.
- [ ] **Metadados**: Capa e descrição DEVEM vir do AniList, nunca do site de mídia.

## 📋 2. Fluxo de Episódios (Listing)
- [ ] **AnimeFire**: Validar se `.lEp` retorna a lista completa (URL atual: animefire.plus/io).
- [ ] **AnimePlayer**: Validar se o clique no disparador de episódios funciona.
- [ ] **TopAnimes**: Verificar se a lista de episódios não está vazia (mesmo com CDN instável).

## 🎬 3. Fluxo de Playback (Stream)
- [ ] **Extração de Link**: O sistema deve extrair o `.mp4` ou `.m3u8` final (procurar tag video/source/iframe `link=`).
- [ ] **MPV Launch**: O player deve abrir no Windows com os headers corretos (Referer/UA).
- [ ] **Fallback**: Se a Fonte A falhar, o sistema deve sugerir a Fonte B.

## 📊 4. Dashboard & Estado
- [ ] **Trending**: Clicar em um item da dashboard deve disparar a busca automática por fontes (Não usar o nome como URL).
- [ ] **Persistência**: O histórico de episódios assistidos deve ser mantido.

# 🚫 Relatório de Falhas Conhecidas (Log de Erros)

| Data | Fonte | Falha Detectada | Status |
| :--- | :--- | :--- | :--- |
| 20/03 | TopAnimes | CDN Kojima desativado no site original (falta de fundos). | ⚠️ Monitorando |
| 20/03 | AnimeFire | Mudança de classe de `.div_episodes` para `.lEp`. | 🛠️ Corrigindo Seletor |
| 20/03 | Cineby | Domínio .gd estacionado / Conteúdo migrou. | ❌ Obsoleto |
| 20/03 | Dashboard | Erro de lógica: o sistema tentava tratar o nome do anime como uma URL. | 🛠️ Corrigindo Handler |
