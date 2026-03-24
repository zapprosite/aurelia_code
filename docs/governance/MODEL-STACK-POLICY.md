# Model Stack Policy — Stack Canônico de LLMs (2026)

**Status:** ✅ Vigente
**Criado:** 2026-03-24
**Autoridade:** Aurélia (Arquiteta Principal) + Will (Principal Engineer)
**Rule associada:** [`.agents/rules/13-model-stack-policy.md`](../../.agents/rules/13-model-stack-policy.md)

---

## Por que este documento existe

Cada LLM que opera neste repositório (Claude, Gemini, OpenCode, Antigravity) tende a "corrigir" o stack de modelos para o que conhece por padrão. Isso causou regressões repetidas (ex: qwen3.5 reintroduzido múltiplas vezes). Este documento é a **fonte de verdade imutável** — qualquer mudança aqui requer revisão explícita do Will.

---

## Stack Canônico Vigente

| Camada | Modelo | Provedor | Motivo |
|--------|--------|----------|--------|
| **Local residente** | `gemma3:12b` | Ollama (localhost:11434) | Melhor custo/qualidade local com 8GB VRAM |
| **Local laboratório** | `gemma3:27b-it-q4_K_M` | Ollama | Raciocínio profundo, uso manual |
| **Embedding** | `bge-m3` | Ollama | Qdrant, sempre local, multilingual |
| **Cloud rápido** | `google/gemini-2.5-flash` | OpenRouter | Baixo custo, alta velocidade |
| **Cloud profundo** | `google/gemini-2.5-pro` | OpenRouter | Análise longa, raciocínio avançado |
| **STT** | `whisper-large-v3-turbo` | Groq | Transcrição rápida PT-BR |
| **TTS** | Kokoro / voice-proxy | Local CPU | Voz oficial Aurélia, sem custo |

---

## Modelos Removidos (PROIBIDO reintroduzir)

| Modelo | Removido em | Substituto |
|--------|------------|-----------|
| `qwen3.5:9b` | 2026-03-24 | `gemma3:12b` |
| `qwen3.5:4b` | 2026-03-24 | `gemma3:12b` |
| `qwen/qwen3.5-flash-02-23` | 2026-03-24 | `google/gemini-2.5-flash` |
| `qwen/qwen3.5-9b` | 2026-03-24 | `google/gemini-flash-1.5` |
| Codex CLI auth | 2026-02 | API key direto (openai provider) |

---

## Arquivos que definem o stack (não alterar sem ADR)

```
internal/gateway/cost.go           ← modelCosts map (tiers 0–4)
internal/config/config.go          ← AppConfig.LLMModel default
cmd/aurelia/onboard.go             ← onboarding default model
scripts/update-ollama.sh           ← MAIN_MODEL / LIGHT_MODEL
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

Se um agente reintroduzir modelo proibido, o sinal aparece em:
```bash
grep -r "qwen3.5\|codex" internal/ cmd/ scripts/ .agents/ .opencode/
```
Se retornar algo além de docs históricos → **reverter imediatamente**.
