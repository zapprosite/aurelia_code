# ADR 20260330: Cascade LiteLLM — Qwen 3.6 Free + Prioridades Explícitas

## Status
✅ Aceito

## Contexto
O `config.yaml` do LiteLLM usa `routing_strategy: latency-based-routing` sem prioridades explícitas. A ordem de fallback não é determinística. A tarefa pede adicionar `openrouter/qwen/qwen-3.6-plus-preview:free` como segundo tier gratuito (logo após local) e numerar prioridades explicitamente.

## Decisão

### Cascade de Prioridades

| Priority | Modelo | Provider | Tipo |
|---|---|---|---|
| 1 | `ollama/qwen3.5:9b` | Local (RTX 4090) | Free |
| 2 | `openrouter/qwen/qwen-3.6-plus-preview:free` | OpenRouter | Free |
| 3 | `openrouter/minimax/minimax-2.5:free` | OpenRouter | Free |
| 4 | `groq/llama-3.3-70b-versatile` | Groq | Free |
| 10 | `openrouter/minimax/minimax-2.7` | OpenRouter | Paid |
| 11 | `openrouter/moonshotai/kimi-k2.5` | OpenRouter | Paid |

### Configuração LiteLLM
- `routing_strategy: usage-based-routing` (prioridade numérica ao invés de só latência)
- `priority` em cada modelo
- API Key reutiliza `OPENROUTER_API_KEY` já configurado no `.env`

## Consequências
- **Positivo**: Cascade determinística, Qwen 3.6 Free entra na rotação antes do Minimax Free
- **Negativo**: OpenRouter free tier pode ter rate limits; monitorar
