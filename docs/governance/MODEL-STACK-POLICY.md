# Model Stack Policy — Stack Canônico de LLMs (2026)

**Status:** ✅ Vigente
**Criado:** 2026-03-24
**Autoridade:** Aurélia (Arquiteta Principal) + Will (Principal Engineer)
**Referência operacional associada:** [`.agent/rules/README.md`](../../.agent/rules/README.md)

---

## Por que este documento existe

Cada LLM que opera neste repositório (Claude, Gemini, OpenCode, Antigravity) tende a "corrigir" o stack de modelos para o que conhece por padrão. Isso causou regressões repetidas (ex: qwen3.5 reintroduzido múltiplas vezes). Este documento é a **fonte de verdade imutável** — qualquer mudança aqui requer revisão explícita do Will.

---

## Motor de Inferência da Aurélia (Stack Canônico Pinado 2026)

| Tier | Modelo | Provedor | Prioridade | Uso |
|------|--------|----------|------------|-----|
| **Nível 0 (Local)** | `gemma3:27b` | Ollama (RTX 4090) | **1** | Juiz/Executor, Timeout 10s |
| **Nível 1 (Cloud Free)**| `minimax-01:free` | OpenRouter | **2** | Linha de frente gratuita |
| **Nível 1 (Cloud Free)**| `gemini-2.0-flash` | Google API | **3** | Visão/Arquivos Grandes |
| **Nível 1 (Cloud Free)**| `llama-3.3-70b` | Groq | **4** | Velocidade absoluta LPU |
| **Nível 2 (Paid)** | `glm-5` | OpenRouter | **5** | Engenharia/Lógica Densa |
| **Nível 2 (Paid)** | `minimax-2.7` | OpenRouter | **6** | Precisão Extrema (Refactor) |
| **Nível 2 (Paid)** | `kimi-2.5` | OpenRouter | **7** | Contexto Massivo (Arquitetura) |

**Embedding**: `nomic-embed-text` (Ollama Local)

---

## Modelos Removidos do Runtime (PROIBIDO reintroduzir)

| Modelo | Removido em | Substituto no Runtime |
|--------|------------|----------------------|
| `gemma3:12b` | 2026-03-28 | `gemma3:27b` |
| `deepseek/chat` | 2026-03-28 | `glm-5` (Tier 2) |
| `bge-m3` | 2026-03-26 | `nomic-embed-text` |

---

## Arquivos que definem o stack (não alterar sem ADR)

```
internal/gateway/cost.go           ← modelCosts map (tiers 0–4)
internal/config/config.go          ← AppConfig.LLMModel default
cmd/aurelia/onboard.go             ← onboarding default model
scripts/update-ollama.sh           ← MAIN_MODEL / EMBED_MODEL
.opencode/agents/aurelia.md        ← Stack documentado para OpenCode
.agent/skills/homelab-control/SKILL.md   ← contexto operacional do homelab
.agent/rules/README.md                   ← índice vigente de regras e guias
```

---

## Como propor mudança

1. Abrir issue ou slice com benchmark comparativo
2. Criar ADR: `docs/ADR.md` → nova entrada com justificativa
3. Atualizar este documento + referência operacional consolidada em `.agent/rules/README.md`
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
