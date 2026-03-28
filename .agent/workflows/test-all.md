---
description: Executa a suite completa de testes e validação de conformidade (SOTA 2026.1)
---

# //test-all

Este workflow orquestra a validação técnica, estrutural e soberana do ecossistema Aurélia.

## Quando usar
- Após qualquer mudança significativa no código (`cmd/`, `internal/`)
- Antes de subir PRs ou dar `git ship`
- Após executar `@sincronizar-ai-context`

---

## Execução

// turbo-all

### 1. Sincronizar Contexto Estrutural
Garante que o `codebase-map.json` e os `.context/docs` estão em paridade com o checkout atual.
```
context({ action: "fill", target: "docs" })
```

### 2. Auditoria de Paridade .env
Valida se as variáveis de ambiente seguem o contrato soberano sem segredos expostos.
```bash
./scripts/audit/audit-env-parity.sh
```

### 3. Auditoria de Inferência (Aurelia Smart)
Valida a conectividade do `gemma3:12b` local e a cascata de failover.
```bash
./.agent/skills/aurelia-smart-validator/scripts/audit-llm.sh
```

### 4. Testes de Unidade e Integração (Go)
Executa a suite Core refatorada com `Testify` e `ALog`.
```bash
go test ./...
```

---

## Output esperado
- Relatório de conformidade sem falhas críticas.
- Log de latência da inferência local.
- Cobertura de testes backend validada.
