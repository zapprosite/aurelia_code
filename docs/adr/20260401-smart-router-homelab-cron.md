# ADR-20260401: Smart Router 3-Camadas + Homelab Cron 5h

**Status:** ✅ Accepted  
**Data:** 01/04/2026  
**Autor:** Aurélia (Agente)  
**Slice:** S-35

---

## Contexto

O projeto Aurélia Sovereign 2026 opera em home lab com GPU RTX 4090. O objetivo é implementar:

1. **Smart Router LiteLLM** com roteamento automático em 3 camadas
2. **Cron job de home lab** enviando resumo a cada 5 horas via Telegram

---

## Decisão

### Router 3-Camadas

| Camada | Modelo | Propósito |
|--------|--------|-----------|
| **Layer 0** | `gemma3:27b-it-qat` (Ollama local) | ops/cron/fiscal — dados operacionais não saem do host |
| **Layer 1** | `nemotron-3-super-120b:free` → `qwen3.6-plus:free` | Chat público free — fallback entre modelos free |
| **Layer 2** | `minimax-m2.7` → `glm-4-flash` → `kimi-k2.5` | Long-context pago — último recurso |

### Por que Nemotron 3 Super como free primary?

- **262K context window** — suporte a conversas longas
- **Agentic capabilities** — otimizado para reasoning e tool use
- **Apache 2.0 weights** — sem restrições comerciais
- **Gratuito via OpenRouter** — custo zero para chat público

### Por que ops-cron usa gemma3:27b local?

- **Dados operacionais sensíveis** — métricas, logs, status de containers
- **Zero cloud exposure** — everything stays on-premise
- **Baixa latência** — ~200ms vs ~2s para cloud
- **Custo zero** — sem API key ou quota limits

### Cron Job Homelab 5h

- **Frequência:** `0 */5 * * *` (a cada 5 horas)
- **Alias LiteLLM:** `ops-cron` (gemma3:27b local)
- **Métricas coletadas:**
  - CPU/GPU Temperature
  - VRAM/RAM usage
  - Container count
  - Ollama/LiteLLM/Qdrant health
  - ZFS available space

---

## Modelo de Configuração

```yaml
model_list:
  - model_name: chat-default
    litellm_params:
      model: openrouter/nvidia/nemotron-3-super-120b-a12b:free
      
  - model_name: ops-cron
    litellm_params:
      model: ollama/gemma3:27b-it-qat
      api_base: http://172.17.0.1:11434

router_settings:
  routing_strategy: latency-based-routing
  fallbacks:
    - {"ops-cron": ["chat-free-second", "chat-paid-glm"]}
```

---

## Tech Debt

1. **Métricas de rota ativa no Grafana** — tracking de qual camada foi usada por request
2. **Rate limiting por tier** — evitar quota exhaustion no free tier
3. **Fallback visual indicator** — mostrar ao usuário qual modelo respondeu

---

## Alternativas Consideradas

| Alternativa | Motivo de Rejeição |
|-------------|-------------------|
| Usar GPT-4o free | Custo + latência alta |
| Usar apenas cloud | Exposição de dados operacionais |
| Cron 1h | Too frequent, noise |
| Cron 24h | Perda de detecção de anomalies |
