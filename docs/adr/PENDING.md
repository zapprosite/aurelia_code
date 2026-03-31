# ADR PENDING: Slices em Espera

**Data:** 2026-03-31
**Status:** Proposto / Pendente

---

## P1 — Crítico

| # | Slice | Descrição | Status |
|---|-------|-----------|--------|
| 1 | Computer Use E2E | Agent loop autônomo (BUA-style) | 🔴 Pendente |
| 2 | OS Native God Mode | Automação desktop Ubuntu nativa | 🔴 Pendente |
| 3 | Jarvis Voice + Computer | Wake → STT → LLM → TTS → Browser | 🔴 Pendente |

---

## P2 — Alto

| # | Slice | Descrição | Status |
|---|-------|-----------|--------|
| 4 | Dashboard Perplexity | UI estilo Perplexity search | 🟡 Proposto |
| 5 | Supabase Runtime | Supabase integrado ao runtime | 🟡 Proposto |
| 6 | Obsidian Integration | Obsidian no runtime | 🟡 Proposto |
| 7 | Visual Cortex | Detecção objetos + OCR + Screen | 🟡 Proposto |

---

## P3 — Médio

| # | Slice | Descrição | Status |
|---|-------|-----------|--------|
| 8 | TTS BR Industrialize | Voxtral vs Kokoro-82M benchmark | 🟢 Em progresso |
| 9 | Tavily Integration | Web search nativo | 🟢 Configurado |
| 10 | Acessibilidade Universal | Portabilidade universal | 🟢 Proposto |

---

## Model Stack Policy (Ativo)

```
Proibidos (não reintroduzir):
  gemma3:27b, gemma3:12b, groq/whisper, deepseek/chat, bge-m3

Cascade Motor:
  Nível 0: qwen3.5 + faster-whisper-v3 (local)
  Nível 1: minimax-01:free, gemini-2.0-flash (cloud free)
  Nível 2: glm-5, minimax-2.7, kimi-2.5 (paid)
  Nível 3: aurelia-top, aurelia-audio (SOTA)
```

---

## Skills Pendentes de Implementação

| Skill | Status | Prioridade |
|-------|--------|------------|
| `computer-use-agent` | ❌ Não implementada | P1 |
| `god-mode-desktop` | ❌ Não implementada | P1 |
| `visual-cortex` | ❌ Não implementada | P2 |
| `perplexity-dashboard` | ❌ Não implementada | P2 |

---

## Containers Pendentes

| Container | Status | Dependência |
|----------|--------|-------------|
| `computer-use-steel` | ❌ Não criado | P1 |
| `browser-isolated` | ❌ Não criado | P1 |
| `obsidian-runtime` | ❌ Não criado | P2 |

---

## Próximos Passos

1. **ADR-0001**: Criar `computer-use-agent` skill
2. **ADR-0002**: Implementar `god-mode-desktop`
3. **ADR-0003**: Benchmark Voxtral vs Kokoro-82M
4. **ADR-0004**: Dashboard Perplexity UI

---

**Última atualização:** 2026-03-31
