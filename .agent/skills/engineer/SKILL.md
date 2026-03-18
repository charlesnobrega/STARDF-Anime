    ---
    name: engineer
    description: Protocolo de engenharia geral: diagnosticar -> planejar -> implementar -> validar -> documentar -> retro (token-safe).
    ---
    # Engineer Protocol (universal + token-safe)

## Quando usar
- Feature/bug/refactor “normal” (end-to-end).
- Repo desconhecido (junto com core-context).

## Protocolo
1) Diagnóstico mínimo
- Identifique: comportamento esperado, onde acontece, como reproduzir (3 bullets).
- Se faltar dado crítico: 1–3 perguntas, no máximo.

2) Bootstrap (se necessário)
- Rodar core-context para descobrir comandos reais e pontos de extensão.

3) Planejamento curto
- Se for grande: use /plan.
- Se for pequeno: escreva 3–6 passos com checkpoint de validação.

4) Implementação em fatias
- Cada fatia termina com evidência: teste rodando / smoke / screenshot / log.

5) Verificação
- Chame QA skill se precisar de estratégia de testes.
- Chame Security skill se tocar auth, input, permissões, secrets.

6) Review e docs
- Rode /review.
- Atualize docs mínimas (como rodar/testar) se impactar DX.

7) Retro
- Propor até 3 melhorias em rules/skills/workflows (texto exato). Não aplicar sem OK.
