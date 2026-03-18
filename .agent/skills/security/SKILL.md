    ---
    name: security
    description: Auditoria de segurança prática: secrets, authn/z, OWASP, dependências, hardening.
    ---
    # Security Protocol

- Secrets: nunca em código; env/secret manager.
- AuthZ: checar permissões por ação, não por tela.
- Input: sanitize/validate; evitar eval/exec.
- Dependências: sinalizar libs críticas desatualizadas.
Se risco alto: parar e pedir OK antes de prosseguir.
