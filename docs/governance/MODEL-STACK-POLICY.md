# Model Stack Policy — Stack Canônico de LLMs (2026)

**Status:** ✅ Vigente
**Criado:** 2026-03-24
**Autoridade:** Aurélia (Arquiteta Principal) + Will (Principal Engineer)
**Rule associada:** [`.agents/rules/13-model-stack-policy.md`](../../.agents/rules/13-model-stack-policy.md)

---

## Por que este documento existe

Cada LLM que opera neste repositório (Claude, Gemini, OpenCode, Antigravity) tende a "corrigir" o stack de modelos para o que conhece por padrão. Isso causou regressões repetidas (ex: qwen3.5 reintroduzido múltiplas vezes). Este documento é a **fonte de verdade imutável** — qualquer mudança aqui requer revisão explícita do Will.

---

## Motor de Inferência da Aurélia (Stack Canônico Vigente)

> Claude, Antigravity e OpenCode são **orquestradores externos** controlados por Will — não são modelos internos da Aurélia.

| Tier | Modelo | Provedor | Uso |
|------|--------|----------|-----|
| **Tier 0 — Local** | `gemma3:12b` | Ollama (RTX 4090) | Fallback universal, custo zero |
| **Tier 1 — Cheap Remote** | `deepseek/deepseek-chat-v3.1` | OpenRouter | Curation, structured output, routing |
| **Tier 2 — Premium Remote** | `minimax/minimax-m2.7` | MiniMax direct | coding_main, critical, execução principal |
| **Tier 2 — Long Context** | `moonshotai/kimi-k2.5` | OpenRouter | Contexto longo, multimodal |
| **Embedding** | `nomic-embed-text` | Ollama (RTX 4090) | Qdrant, sempre local |
| **STT** | `whisper-large-v3-turbo` | Groq | Transcrição rápida PT-BR |
| **TTS** | `kokoro` | Local CPU | Voz oficial Aurélia, sem XTTS nem proxy |

**Fonte de verdade do roteamento:** `internal/gateway/policy.go`

---

## Modelos Removidos do Runtime (PROIBIDO reintroduzir)

| Modelo | Removido em | Substituto no Runtime |
|--------|------------|----------------------|
| `gemma3:27b-it-q4_K_M` | 2026-03-26 | `gemma3:12b` |
| `bge-m3` | 2026-03-26 | `nomic-embed-text` |
| `qwen3.5:9b` | 2026-03-24 | `gemma3:12b` |
| `qwen3.5:4b` | 2026-03-24 | `gemma3:12b` |
| `qwen/qwen3*` | 2026-03-26 | `deepseek/deepseek-chat-v3.1` |
| `voice-proxy` | 2026-03-26 | `kokoro` em `127.0.0.1:8012` |
| `xtts` / `openedai-speech` | 2026-03-26 | `kokoro` em CPU |
| `google/gemini-2.5-flash` | 2026-03-25 | `deepseek/deepseek-chat-v3.1` (Tier 1) |
| `google/gemini-2.5-pro` | 2026-03-25 | `minimax/minimax-m2.7` (Tier 2) |
| Codex CLI auth | 2026-02 | OpenCode (orquestrador externo) |

---

## Arquivos que definem o stack (não alterar sem ADR)

```
internal/gateway/cost.go           ← modelCosts map (tiers 0–4)
internal/config/config.go          ← AppConfig.LLMModel default
cmd/aurelia/onboard.go             ← onboarding default model
scripts/update-ollama.sh           ← MAIN_MODEL / EMBED_MODEL
.opencode/agents/aurelia.md        ← Stack documentado para OpenCode
.agents/skills/homelab-control/SKILL.md  ← Modelos pinados
.agents/rules/13-model-stack-policy.md   ← Esta policy em form de rule
```

---

## Como propor mudança

1. Abrir issue ou slice com benchmark comparativo
2. Criar ADR: `docs/ADR.md` → nova entrada com justificativa
3. Atualizar este documento + rule 13
4. Atualizar todos os arquivos protegidos listados acima
5. Rodar `go test ./...` para validar testes com novo modelo default

---

## Detecção de regressão

Se um agente reintroduzir modelo proibido ou confundir orquestrador com motor interno:
```bash
# Detectar modelos legados no runtime
grep -r "qwen3\.5\|qwen/qwen3\|bge-m3\|gemma3:27b\|gemini-2\.5\|gemini-flash\|gemini-pro\|google/gemini" internal/ cmd/ scripts/
# Detectar orquestradores no runtime (não deveriam estar aqui)
grep -r "anthropic/claude\|opencode\|antigravity" internal/gateway/ internal/config/
```
Se retornar algo além de docs históricos → **reverter imediatamente**.
