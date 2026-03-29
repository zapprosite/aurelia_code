# ADR 20260328: Smart Router LiteLLM com gemma3 27b como Juiz Soberano

## Status
🟡 Proposto (P1 - Requer Implementation)

## Contexto
O ADR `20260327-smart-router-llm` define uma arquitetura de roteamento inteligente, mas o `configs/litellm/config.yaml` não existe ou não implementa as tiers planejadas. O gemma3 27b deve ser o **Juiz Soberano local** (Tier 0) com fallback para OpenRouter (Tier 1).

## Decisões Arquiteturais

### 1. Arquitetura de Tiers

```yaml
# configs/litellm/config.yaml

model_list:
  # Tier 0: Local Sovereign (VRAM Priority)
  - model_name: "gemma3:27b"
    litellm_params:
      model: "ollama/gemma3:27b"
      api_base: "http://127.0.0.1:11434"
      rpm: 999999
      priority: 0

  - model_name: "gemma3:12b"
    litellm_params:
      model: "ollama/gemma3:12b"
      api_base: "http://127.0.0.1:11434"
      rpm: 999999
      priority: 1

  # Tier 1: OpenRouter (Fallback/Scale)
  - model_name: "qwen-coder-32b"
    litellm_params:
      model: "openrouter/qwen/qwen-2.5-coder-32b-instruct"
      api_base: "https://openrouter.ai/api/v1"
      rpm: 60
      priority: 2

  - model_name: "minimax-01"
    litellm_params:
      model: "openrouter/minimax/minimax-01"
      api_base: "https://openrouter.ai/api/v1"
      rpm: 30
      priority: 2

  - model_name: "kimi-k2"
    litellm_params:
      model: "openrouter/moonshotai/kimi-k2.5"
      api_base: "https://openrouter.ai/api/v1"
      rpm: 60
      priority: 2
```

### 2. Router Group "aurelia-smart"

```yaml
litellm_settings:
  drop_params: true
  set_verbose: false

router_settings:
  model_group_alias:
    "aurelia-smart":
      - "gemma3:27b"  # Prioridade 0
      - "gemma3:12b"  # Fallback 1
      - "qwen-coder-32b"  # Tier 1
      - "minimax-01"  # Tier 1
      - "kimi-k2"  # Tier 1
```

### 3. LiteLLM API Endpoints

```yaml
# LiteLLM Gateway expose
litellm_settings:
  master_key: "${LITELLM_MASTER_KEY}"

general_settings:
  master_key: "${LITELLM_MASTER_KEY}"
  database_url: "sqlite:///litellm.db"

environment_variables:
  OPENROUTER_API_KEY: "${OPENROUTER_API_KEY}"
  OLLAMA_API_KEY: "sk-ollama"  # Dummy, Ollama não requer key
```

### 4. Health Checks

```go
// internal/gateway/provider.go - Health check for each tier
func (p *Provider) checkTierHealth(ctx context.Context, tier int) bool {
    switch tier {
    case 0: // Ollama
        return p.checkOllamaHealth(ctx)
    case 1: // OpenRouter
        return p.checkOpenRouterHealth(ctx)
    }
    return false
}
```

### 5. Retry Logic

- **Tier 0 → Tier 1**: 1 retry on timeout/error
- **Tier 1 errors**: Circuit breaker (5 failures = 30s cooldown)
- **Context length**: Auto-escalate to OpenRouter if > 8k tokens

## Consequências

### Positivas
- gemma3 27b como Juiz Soberano: decisões de roteamento feitas localmente
- Zero custo para tarefas Tier 0
- Soberania de dados preservada
- Escalamento automático para modelos mais capaz quando necessário

### Negativas
- Latência adicional quando fallback para OpenRouter
- Custo variável em produção (Tier 1 é pay-per-use)
- Complexidade operacional: múltiplos providers para monitorar

### Trade-offs
- Local vs Cloud: Garante soberania, mas Gemma3 27b pode ser insuficiente para tarefas complexas
- Custo vs Qualidade: Tier 1 (Qwen Coder) é excellent para código, mas adiciona latência

## Dependências
- ✅ `pkg/llm/ollama.go` (rankOllamaModel com gemma3:27b = 0)
- ⚠️ `configs/litellm/config.yaml` (NÃO EXISTE - precisa ser criado)
- ⚠️ `internal/gateway/provider.go` (precisa ler config e health check)
- ❌ Environment vars: `LITELLM_MASTER_KEY`, `OPENROUTER_API_KEY`

## Referências
- [ADR-20260327-smart-router-llm.md](./20260327-smart-router-llm.md)
- [internal/gateway/provider.go](../../internal/gateway/provider.go)
- [pkg/llm/ollama.go](../../pkg/llm/ollama.go)
- [LiteLLM Router Docs](https://docs.litellm.ai/docs/routing)

## Links Obrigatórios
- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)

---
**Data**: 2026-03-28
**Status**: Proposto
**Autor**: Claude (Principal Engineer)
**Slice**: feature/neon-sentinel
**Progress**: 0%
