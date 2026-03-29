---
description: Define o motor de inferência canônico da Aurélia e proíbe regressão para modelos legados sem ADR aprovado.
id: 13-model-stack-policy
---

# 🔒 Regra 13: Motor de Inferência da Aurélia (IMUTÁVEL SEM ADR)

Esta regra define o **motor interno** do daemon Go `aurelia`.
**Nenhum agente pode alterar este stack sem ADR aprovado registrado em `docs/adr/`.**

> **IMPORTANTE:** Claude, Antigravity e OpenCode são **orquestradores externos** controlados por Will.
> Eles NÃO são modelos internos da Aurélia. O motor de inferência da Aurélia é exclusivamente o stack abaixo.

---

## Motor de Inferência da Aurélia (2026)

| Tier | Modelo | Provedor | Uso |
|------|--------|----------|-----|
| **Tier 0 — Local** | `qwen3.5` | Ollama local (RTX 4090) | Soberania Industrial, Máxima Inteligência |
| **Tier 0 — Local Lab** | `qwen3.5-it-q4_K_M` | Ollama local | Raciocínio profundo, uso manual |
| **Tier 1 — Cheap Remote** | `deepseek/deepseek-chat-v3.1` | OpenRouter | Curation, structured output, routing barato |
| **Tier 2 — Premium Remote** | `minimax/minimax-m2.7` | MiniMax direct | coding_main, critical, execução principal |
| **Tier 2 — Long Context** | `moonshotai/kimi-k2.5` | OpenRouter | Contexto longo, multimodal, análise profunda |
| **Embedding** | `bge-m3` | Ollama local | Qdrant, sempre local, multilingual |
| **STT** | `whisper-large-v3-turbo` | Groq | Transcrição rápida PT-BR |
| **TTS** | Kokoro / voice-proxy | Local CPU | Voz oficial Aurélia |

**Fonte de verdade do roteamento:** `internal/gateway/policy.go`

---

## Orquestradores Externos (NÃO são motor da Aurélia)

| Motor | Papel | Controlado por |
|-------|-------|----------------|
| **Claude Code** | IDE/CLI de desenvolvimento | Will |
| **Antigravity** | IDE de interface e coordenação | Will |
| **OpenCode** | Executor de implementação | Will |

→ Estes motores **nunca** devem aparecer em `internal/config/config.go`, `internal/gateway/cost.go` ou qualquer config de runtime.

---

## Modelos PROIBIDOS (legado — não reintroduzir)

```
gemma3:27b        ← removido em 2026-03-29, substituído por qwen3.5
gemma3:12b        ← removido em 2026-03-29
gemma/gemma3-*    ← removido do Tier 1/2
google/gemini-*   ← removido do motor interno (era placeholder, substituído por DeepSeek/MiniMax/Kimi)
codex             ← CLI auth removida, não reintroduzir
```

Ao encontrar qualquer referência a esses modelos no código de runtime:
→ **Substituir pelo tier equivalente da tabela acima**
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
.agents/skills/homelab-control/    ← Modelos pinados
```

---

## Referências

- ADR de política de modelos: `docs/adr/` (buscar `politica-modelos` ou `model-stack`)
- Regra de autonomia/tiers: [03-tiers-autonomy.md](./03-tiers-autonomy.md)
- Regra de rede/infra: [12-network-governance.md](./12-network-governance.md)
