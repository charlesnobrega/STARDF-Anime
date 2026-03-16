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
- [ ] Fontes: Habilitar/Corrigir segundo servidor de Filmes (Cineby/CineGratis) <!-- id: 15 -->
- [ ] Validação: Testes de acesso online aos scrapers (Bloqueado por rede) <!-- id: 16 -->

## Próximas Propostas (Aguardando Aprovação)
- [ ] **Cache de Metadados Local**: Implementar SQLite para salvar resultados de busca e listas de episódios (reduz latência).
- [ ] **Sistema de Plugins para Scrapers**: Modularizar scrapers para facilitar adição/remoção sem alterar o core.
- [ ] **Migração para Bubbletea**: Interface TUI mais rica e informativa com painéis e progresso visual.
- [ ] **Integração com MAL/AniList (Sync)**: Sincronizar progresso de visualização automaticamente com contas do usuário.

## Rollback
Se Release Discovery indica falha crítica:
1. Reverter para última tag estável (`v1.6.2`)
2. Comunicar emIssue da release

## Pós-entrega
- Monitorar issues novas por 7 dias
- Responder bugs dentro de 48h úteis
