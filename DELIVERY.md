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

## Checklist de entrega
- [ ] Compilação cross-platform (GOOS/GOARCH) testada
- [ ] Teste manual: `./goanime --source goyabu "naruto"` retorna resultados
- [ ] README.md com comandos básicos
- [ ] CHANGELOG.md com mudanças depuis última release
- [ ] Tag criada e pushada
- [ ] GitHub Release draft criado com assets
- [ ] Nota de release completa

## Rollback
Se Release Discovery indica falha crítica:
1. Reverter para última tag estável (`v1.6.2`)
2. Comunicar emIssue da release

## Pós-entrega
- Monitorar issues novas por 7 dias
- Responder bugs dentro de 48h úteis
