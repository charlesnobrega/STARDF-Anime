# DELIVERY.md — Definição de "Finalizar/Publicar"

## Escopo de entrega

Para este projeto (STARDF-Anime), a entrega significa:

### Tipo: GitHub Release
- Versão tag: `v1.6.3` (ou maior)
- Assets:
  - `goanime-linux-amd64`
  - `goanime-windows-amd64.exe`
  - `goanime-darwin-amd64` (se aplicável)
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
- [ ] Compilação cross-platform (GOOS/GOARCH) testada
- [ ] Teste manual: `./stardf-anime --source animefire "naruto"` retorna resultados
- [ ] Correção das Issues #1 e #2 (Scrapers WordPress)
- [ ] Atualização do README.md com a nova identidade
- [ ] CHANGELOG.md atualizado com a desvinculação do fork
- [ ] Tag `v1.6.3` criada e pushada para o novo repositório
- [ ] GitHub Release draft criado com os novos assets

## Rollback
Se Release Discovery indica falha crítica:
1. Reverter para última tag estável (`v1.6.2`)
2. Comunicar emIssue da release

## Pós-entrega
- Monitorar issues novas por 7 dias
- Responder bugs dentro de 48h úteis
