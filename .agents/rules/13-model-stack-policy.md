---
description: Define o stack canônico de LLMs do projeto e proíbe regressão para modelos legados sem ADR aprovado.
id: 13-model-stack-policy
---

# 🔒 Regra 13: Política de Model Stack (IMUTÁVEL SEM ADR)

Esta regra define o stack oficial de modelos do projeto Aurélia.
**Nenhum agente pode alterar o model stack sem ADR aprovado registrado em `docs/adr/`.**

---

## Stack Canônico (2026)

| Camada | Modelo | Provedor | Alterável? |
|--------|--------|----------|-----------|
| **Local residente** | `gemma3:12b` | Ollama local | ❌ Apenas via ADR |
| **Local laboratório** | `gemma3:27b-it-q4_K_M` | Ollama local | ❌ Apenas via ADR |
| **Cloud rápido** | `google/gemini-2.5-flash` | OpenRouter | ❌ Apenas via ADR |
| **Cloud profundo** | `google/gemini-2.5-pro` | OpenRouter | ❌ Apenas via ADR |
| **Embedding** | `bge-m3` | Ollama local | ❌ Apenas via ADR |
| **STT** | Groq | Remote | ❌ Apenas via ADR |
| **TTS** | Kokoro / voice-proxy | Local CPU | ❌ Apenas via ADR |

---

## Modelos PROIBIDOS (legado — não reintroduzir)

```
qwen3.5:9b        ← removido em 2026-03-24, substituído por gemma3:12b
qwen3.5:4b        ← removido em 2026-03-24
qwen/qwen3.5-*    ← removido do Tier 1/2
codex             ← CLI auth removida, não reintroduzir
```

Ao encontrar qualquer referência a esses modelos no código, docs ou configs:
→ **Substituir por `gemma3:12b` (local) ou `google/gemini-2.5-flash` (cloud)**
→ **Nunca restaurar silenciosamente**

---

## O que NÃO pode ser feito sem ADR

<directives>
1. Mudar o modelo padrão em `internal/config/config.go`, `onboard.go`, ou qualquer config default.
2. Adicionar novo modelo ao Tier 0 (local gratuito) em `internal/gateway/cost.go`.
3. Mudar o provedor padrão (ollama → outro) nos testes ou no onboarding.
4. Reintroduzir qualquer modelo da lista PROIBIDOS acima.
5. Substituir OpenRouter por outro provedor cloud sem consenso.
6. Mudar o modelo de embedding `bge-m3`.
</directives>

---

## Como propor mudança de modelo

1. Criar ADR em `docs/adr/ADR-AAAAMMDD-politica-modelos-<motivo>.md`
2. Justificar: benchmark, custo, VRAM, latência
3. Atualizar esta rule após aprovação
4. Atualizar `internal/gateway/cost.go`, `scripts/update-ollama.sh`, `.opencode/agents/aurelia.md`

---

## Arquivos protegidos (não alterar modelo sem ADR)

```
internal/gateway/cost.go           ← Tier 0–4 model pricing
internal/config/config.go          ← AppConfig defaults
cmd/aurelia/onboard.go             ← onboarding default model
scripts/update-ollama.sh           ← MAIN_MODEL / LIGHT_MODEL
.opencode/agents/aurelia.md        ← Stack de modelos documentado
.agents/skills/homelab-control/    ← Modelos pinados
```

---

## Referências

- ADR de política de modelos: `docs/adr/` (buscar `politica-modelos` ou `model-stack`)
- Regra de autonomia/tiers: [03-tiers-autonomy.md](./03-tiers-autonomy.md)
- Regra de rede/infra: [12-network-governance.md](./12-network-governance.md)
