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
- [x] Correção de imports e nomenclatura (GoAnime -> StarDF-Anime) <!-- id: 5 -->
- [x] Padronização de nomes de binários nos scripts de build <!-- id: 20 -->
- [x] Atualização de toda a documentação interna (.md) <!-- id: 21 -->
- [x] Compilação cross-platform (GOOS/GOARCH) concluída (v1.6.3) <!-- id: 24 -->
- [ ] Teste manual: `./stardf-anime --source animefire "naruto"` retorna resultados <!-- id: 22 -->
- [x] Correção das Issues #1 e #2 (Scrapers WordPress) - CONCLUÍDO <!-- id: 1 -->
- [x] Atualização do README.md e README_pt-br.md com a nova identidade
- [x] Versão Mobile para Android (Projeto Iniciado / Mockups) <!-- id: 6 -->
- [ ] Instalador Windows (.exe) (Gerar via Inno Setup) <!-- id: 7 -->
- [x] Binários para Linux e macOS (Liberados na build/) <!-- id: 8 -->
- [x] Suporte a NixOS Flakes (Atualizado) <!-- id: 9 -->
- [x] Criação de Comunidade Oficial (Discord Badge Adicionada) <!-- id: 10 -->
- [x] Reavaliação de Scrapers Instáveis (Goyabu/SuperAnimes - Atualmente OFFLINE) <!-- id: 11 -->
- [x] Tag `v1.6.3` criada e pushada para o repositório <!-- id: 12 -->
- [x] Assets de Release v1.6.3 gerados em `build/` <!-- id: 13 -->
- [x] Navegação: Permitir retornar ao menu inicial (Animes/Filmes) <!-- id: 14 -->
- [x] Fontes: Habilitar/Corrigir segundo servidor de Filmes (Cineby/CineGratis) <!-- id: 15 -->
- [ ] Validação: Testes de acesso online aos scrapers (Bloqueado por rede) <!-- id: 16 -->
- [!] INVESTIGAÇÃO: Erro "Access denied" nos testes de socket (Bloqueio de segurança do OS para binários temporários de teste) <!-- id: 23 -->

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
