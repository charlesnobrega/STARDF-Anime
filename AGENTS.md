# AGENTS.md — Governança do Repositório

## Regra de ouro
CI e testes mandam. Sem "achismo".

## Comandos oficiais
```bash
# Build
go build -o goanime cmd/goanime/main.go

# Test (unit)
go test ./...

# Lint
golangci-lint run

# Integration test (local)
./goanime --source goyabu "test"
```

## Convenções
- Commits pequenos, PRs pequenos
- Não mexer em auth/pagamentos/infra sem testes dedicados
- Seguir modelo de Issue (ver .github/ISSUE_TEMPLATE)

## Definition of Done
- Aceite completo checklist ✓
- Testes verdes (CI passing)
- Scans de segurança (se aplicável) ✓
- DELIVERY.md atualizado (quando applicable)

## Tipos de Issue
| Label | Uso |
|-------|-----|
| bug | Comportamento incorreto |
| enhancement | Melhoria/nova feature |
| security | Vulnerabilidade |
| techdebt | Refatoração sem mudança externa |

## Labels de Workflow
| Label | Significado |
|-------|-------------|
| triage | Nova, aguardando análise |
| negotiating | Escopo em definição |
| approved | Pronta para execução |
| in_progress | Em andamento |
| waiting_ci | Aguardando CI |
| needs_fix | Problema encontrado |
| needs_antigravity | Requer revisão profunda |
| ready_for_release | Pronta para deploy |
| delivered | Concluída |

## Controles
- `hold` → PARA TUDO (bloqueia execute)
- `security_sensitive` → exige revisão humana
- `major_change` → exige aprovação adicional

## Protocolo de Execução
1. Issue com label `approved` e sem `hold` dispara executor.
2. Comentário `EXECUTION BRIEF` inicia trabalho.
3. Atualize labels conforme progride.
