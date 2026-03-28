# 16. Testing & Validation Governance (v2026.1)

A integridade do ecossistema Aurélia é garantida por auditorias contínuas.

## Obrigatoriedades

1. **Protocolo //test-all**: É obrigatório rodar a suite completa de testes antes de qualquer transição de fase (P -> R -> E -> V).
2. **Sincronia com ai-context**: O comando `//test-all` deve ser invocado imediatamente após `@sincronizar-ai-context` se houver mudanças em `internal/`.
3. **Contrato de Falha**: Nenhuma feature pode ser considerada "Complete" se houver regressão no Smart Router Audit ou nos testes Go Core.
4. **Logs Estruturados**: Todo novo teste deve utilizar a biblioteca `internal/purity/alog` para garantir rastreabilidade SOTA 2026.1.

## Gatilhos de Auditoria

- **Subdomain Change**: Rodar `audit-llm.sh` após adicionar subdomínios.
- **Model Switch**: Validar `config.yaml` via roteador antes de novos deploys.
- **Secret Rotation**: Rodar `audit-env-parity.sh` imadiatamente após atualizar chaves.
