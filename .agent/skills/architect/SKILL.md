    ---
    name: architect
    description: Arquitetura e refatoração segura: modularidade, boundaries, ADR, migração incremental.
    ---
    # Architecture Protocol

Checklist:
- Definir objetivo da refatoração (performance? manutenção? bugs?)
- Mapear dependências antes de mover coisas.
- Criar “estrangulamento”: nova rota/módulo coexistindo com antigo.
- Exigir evidência: testes / smoke / build verde.
Saída: propor ADR curto em `docs/adr/NNN-*.md` quando decisão for relevante.
