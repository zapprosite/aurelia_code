# ADR Slice: LiteLLM Cascade — Qwen 3.6 Free + Priorities

## Status
✅ Implementado

## Contexto
LiteLLM proxy já rodando via Docker Compose (serviço `smart-router` na porta 4000). O config.yaml anterior usava `latency-based-routing` sem prioridade explícita. A cascade precisa ser determinística.

## Cascade de Prioridades

| Priority | Modelo | Provider | Tipo | Notes |
|---|---|---|---|---|
| 1 | `ollama/qwen3.5:9b` | Local RTX 4090 | Free | Padrão, ~6GB VRAM, think:false |
| 2 | `openrouter/qwen/qwen-3.6-plus-preview:free` | OpenRouter | Free | Qwen 3.6 Free |
| 3 | `openrouter/minimax/minimax-2.5:free` | OpenRouter | Free | Minimax via OpenRouter (sem API direta) |
| 4 | `groq/llama-3.3-70b-versatile` | Groq | Free | 14.400 req/dia |
| 10 | `openrouter/minimax/minimax-m2.7` | OpenRouter | Paid | Último recurso |
| 11 | `openrouter/moonshotai/kimi-k2.5` | OpenRouter | Paid | |

## Configuração (`configs/litellm/config.yaml`)
```yaml
# Note: routing_strategy "priority" NÃO existe no LiteLLM main-stable
# Solução: least-busy (prefere menos ocupado = local) + model_fallbacks explícitos
# Fix: 2026-03-30 — extra_body: {think: false} dentro de litellm_params (não no nível do model_list)

router_settings:
  routing_strategy: least-busy   # prefere menos ocupado = prefere local RTX 4090
  num_retries: 3
  allowed_fails: 2
  cooldown_time: 30
  timeout: 90
```

### Cascade de modelos (model_fallbacks chain)
```
aurelia-smart (local) ─[fallback]→ qwen3.6-plus-preview:free
                                    ─[fallback]→ minimax-2.5:free
                                    ─[fallback]→ groq/llama-3.3-70b-versatile
                                    ─[fallback]→ minimax-m2.7 (paid)
                                    ─[fallback]→ kimi-k2.5 (paid)
```

## Integração Go
- `cmd/aurelia/app.go`: gateway.NewProvider() aponta para localhost:4000
- Não adicionou caso `litellm` em buildLLMProvider — gateway.NewProvider já cobre isso

## ADR
`docs/adr/20260330-qwen36-free-litellm-cascade.md`

## Consequências
- **Positivo**: Cascade determinística, Qwen 3.6 Free antes do Minimax Free
- **Negativo**: OpenRouter free tier tem rate limits; monitorar
