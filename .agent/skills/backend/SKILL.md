    ---
    name: backend
    description: Backend seguro e testável: APIs, banco, migração, auth, integrações.
    ---
    # Backend Protocol

- Contrato antes de código: endpoints, payloads, erros.
- Validação de entrada sempre.
- Banco: migração reversível + verificação de performance básica.
- Integrações: timeouts/retry/backoff; logs úteis.
Sempre sugerir teste mínimo (unit ou integração) pro caminho feliz + borda.
