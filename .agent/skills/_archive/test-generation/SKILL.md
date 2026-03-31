---
name: Test Generation
description: Geração de testes abrangentes (Unit, Integration, E2E) para o ecossistema Aurélia.
phases: [E, V]
---

# 🧪 Test Generation: Sovereign Validation 2026

Habilita o Antigravity a criar suítes de teste robustas que garantem a estabilidade do Home Lab e a precisão da lógica agêntica.

## 🏛️ Padrões de Teste (Industrial)

### 1. Testes de Unidade (Go/TS)
- **Go**: Use `testing` padrão + mocks manuais ou `gomock`. Foque em cobertura de 80%+.
- **TS (Zod)**: Teste esquemas contra inputs válidos e inválidos.

### 2. Testes de Gateway (Soberanos)
- **Qwen 3.5 Judge**: Utilize o padrão de mock de juiz local para validar se as respostas do roteador satisfazem os guardrails sem gastar tokens premium.
- **Dry-Run**: Sempre teste os fluxos de fallback (Tier 1 -> Tier 3).

### 3. Testes de Infra (Sudo=1)
- Verifique se os scripts bash retornam `exit 0` e produzem os arquivos esperados no Host.

## 🛠️ Comandos Úteis
- `go test -v -cover ./internal/...`
- `npm test -- --coverage`
- `scripts/run-gateway-bench.sh`

## 📍 Quando usar
- Durante a fase `E` (Execution) ao implementar novas funcionalidades.
- Na fase `V` (Verification) para garantir que bugs corrigidos não retornem.
- Para validar a performance de novos modelos de LLM no gateway.