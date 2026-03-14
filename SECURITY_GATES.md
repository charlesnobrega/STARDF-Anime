# SECURITY_GATES.md

## Gates obrigatórios antes de merge/approve

### 1. Secrets scanning
- [ ] Verificação automática via CI (truffleHog/gitleaks)
- [ ] Nenhum segredo encontrado em diffs

### 2. Dependency check
- [ ] `go mod verify`
- [ ] `go mod tidy` executado (sem mudanças inesperadas)
- [ ] CVEs verificadas (se houver dependências externas)

### 3. Static analysis
- [ ] Lint passando
- [ ] `go vet` sem erros críticos
- [ ] Formatação `go fmt` aplicada

### 4. Test coverage
- [ ] Novas funções com testes unitários
- [ ] Testes existentes passando
- [ ] Cobertura não reduzida (se tracking habilitado)

### 5. Performance/Scalability
- [ ] Não introduzir loops O(n²) sem necessity
- [ ] Limitar uso de memória (especialmente scrapers)

### 6. WAF/Blocks considerations
- [ ] Headers realistas (se scraping)
- [ ] Delays appropriados (evitar rate-limit)
- [ ] Cookie handling correto

## Gates automatizados (GitHub Actions)
- Nome do workflow: `.github/workflows/security.yml`
- Roda em: push para main, PRs
- Falha bloqueia merge automático

## Exception process
Qualquer gate que falhar deve:
1. Documentar motivo na Issue/PR
2. Aplicar label `needs_fix` ou `security_sensitive`
3. Exigir revisão humana para bypass
