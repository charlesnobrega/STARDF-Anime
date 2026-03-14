<!--
PR Template — Siga este guia para acelerar a revisão
-->

## Descrição
[Descreva de forma clara e concisa o que esta PR faz]

## Tipo de mudança
- [ ] Bug fix (não quebra nada)
- [ ] New feature (não quebra nada)
- [ ] Breaking change (precisa de aprovação especial)
- [ ] Refatoração/techdebt

## Checklist
- [ ] Testes adicionados/atualizados
- [ ] Lint passando (`golangci-lint run`)
- [ ] `go test ./...` passando
- [ ] `go mod tidy` executado (sem mudanças inesperadas)
- [ ] README/DELIVERY atualizado (se aplicável)
- [ ] SECURITY_GATES revisados

## O que esta PR NÃO faz
[Lista tudo que ficou fora do escopo para evitar scope creep]

## Como testar
```bash
# Passo a passo para validar a mudança
go build -o goanime cmd/goanime/main.go
./goanime --source goyabu "test"
```

## Revisão
- [ ] CI passando
- [ ] Code review aprovada
- [ ] Antigravity (se necessário)

## Links
- Issue relacionada: #XXXX

---
**Nota:** Esta PR será bloqueada se:
- Introduzir secrets/cves
- Quebrar testes existentes
- Alterar auth/pagamento sem testes dedicados
