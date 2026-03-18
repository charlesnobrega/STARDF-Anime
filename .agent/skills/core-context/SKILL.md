    ---
    name: core-context
    description: Bootstrap de contexto do projeto. Use quando iniciar trabalho num repo novo ou confuso.
    ---
    # Core Context Bootstrap (token-safe)

## Objetivo
Criar/atualizar um contexto mínimo e reutilizável do projeto, sem inflar a conversa.

## Protocolo
1) Descobrir “como roda” (comandos reais do repo): build, test, lint, dev server.
2) Mapear arquitetura em 10 bullets: pastas, módulos, fluxos principais, dependências.
3) Registrar um CONTEXTO curto em `docs/AI_CONTEXT.md` (ou equivalente).
4) Se faltar info crítica: pergunte 1–3 perguntas, no máximo.

## Saída padrão (sempre)
- “Como rodar” (3 linhas)
- “Onde mexer” (pastas/arquivos)
- “Riscos” (3 bullets)
