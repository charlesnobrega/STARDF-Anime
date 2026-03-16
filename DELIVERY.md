# DELIVERY.md — Definição de "Finalizar/Publicar"

## Escopo de entrega

Para este projeto (STARDF-Anime), a entrega significa:

### Tipo: GitHub Release
- Versão tag: `v1.6.3`
- Assets:
  - `stardf-anime-linux-amd64`
  - `stardf-anime-windows-amd64.exe`
  - `stardf-anime-darwin-amd64` (se aplicável)
- CHANGELOG atualizado
- Instruções de uso atualizadas no README

### Critérios
- Build bem-sucedido em CI para todas as plataformas alvo
- Testes de integração passam (pelo menos uma fonte funciona)
- Nenhum erro de compilação
- SECURITY_GATES aprovados

## Checklist de entrega (v1.6.3 - Nova Identidade STARDF)
- [x] Migração total de namespace para `charlesnobrega/STARDF-Anime` <!-- id: 4 -->
- [x] Correção de imports em todos os arquivos `.go` <!-- id: 5 -->
- [x] Compilação cross-platform (GOOS/GOARCH) testada
- [x] Teste manual: `./stardf-anime --source animefire "naruto"` retorna resultados
- [x] Correção das Issues #1 e #2 (Scrapers WordPress)
- [x] Atualização do README.md com a nova identidade
- [ ] Versão Mobile para Android <!-- id: 6 -->
- [ ] Instalador Windows (.exe) <!-- id: 7 -->
- [x] Binários para Linux e macOS <!-- id: 8 -->
- [x] Suporte a NixOS Flakes <!-- id: 9 -->
- [ ] Criação de Comunidade Oficial (Discord) <!-- id: 10 -->
- [x] Reavaliação de Scrapers Instáveis (Goyabu/SuperAnimes - Atualmente OFFLINE) <!-- id: 11 -->
- [ ] Tag `v1.6.3` criada e pushada para o novo repositório <!-- id: 12 -->
- [ ] GitHub Release draft criado com os novos assets <!-- id: 13 -->
- [x] Navegação: Permitir retornar ao menu inicial (Animes/Filmes) <!-- id: 14 -->
- [x] Fontes: Habilitar/Corrigir segundo servidor de Filmes (Cineby/CineGratis) <!-- id: 15 -->
- [ ] Validação: Testes de acesso online aos scrapers (Bloqueado por rede) <!-- id: 16 -->

## Implementações Recentes (v1.6.4-preview)
- [x] **Expansão SQLite**: Implementação de tabelas para Watchlist e Monitoramento de Plugins.
- [x] **Menu Principal (TUI)**: Novo fluxo de navegação com Busca, Minha Lista, Continuar Assistindo e Saúde dos Plugins.
- [x] **Acompanhamento (Watchlist)**: Sistema de "Follow/Acompanhar" integrado à seleção de episódios.
- [x] **Relatório de Plugin Health**: Dashboard para desenvolvedores monitorarem falhas nos scrapers.
- [x] **Sincronização de Episódios**: Verificação automática de novos episódios ao iniciar o app.
- [x] **Auto-sync AniList**: Marcação automática de episódio assistido quando >85% concluído.
- [x] **Notificações Desktop**: Alertas nativos (Windows/macOS/Linux) ao encontrar novos episódios.
- [x] **--anilist-login / --anilist-logout**: Flags CLI para conectar/desconectar conta AniList.
- [x] **Token Store Cross-platform**: Token OAuth2 salvo em AppData/Application Support/~/.config.

## Próximas Propostas (Aguardando Aprovação)
- [x] **Integração AniList Real Sync**: Recuperar lista completa do AniList para popular watchlist local.
- [x] **Status AniList no Menu**: Exibir nome do usuário logado no menu principal.
- [x] **Rating Local**: Avaliar obras (1-10) e salvar no SQLite local.
- [x] **Exportar Watchlist**: Exportar lista como JSON/CSV para backup.
- [ ] **Streaming Progressivo P2P**: Investigar suporte a BitTorrent (WebTorrent) para algumas fontes.

## Rollback
Se Release Discovery indica falha crítica:
1. Reverter para última tag estável (`v1.6.2`)
2. Comunicar em Issue da release

## Pós-entrega
- Monitorar issues novas por 7 dias
- Responder bugs dentro de 48h úteis
