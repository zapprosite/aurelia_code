# ADR — Authority Documents

**Diretório:** `docs/adr/`
**Última compactação:** 2026-03-31
**Total:** 8 arquivos + índice

---

## Arquivos

| Arquivo | Conteúdo | Status |
|---------|----------|--------|
| `README.md` | Este índice | — |
| `20260331-architecture-mcp-a2a.md` | Arquitetura + MCP/A2A roadmap | Ativo |
| `20260331-duplicate-response-fix.md` | Correção respostas duplicadas (S-34) | ✅ Implementado |
| `20260331-telegram-fix-s34.md` | Deduplicação + TTS (S-33, S-34) | ✅ Implementado |
| `20260331-sota-industrializacao.md` | Industrialização SOTA 2026 | Ativo |
| `ADR-20260331-polda-2-industrializacao.md` | Police v2 | Arquivado |
| `20260331-telegram-duplicate-fix.md` | Duplicado (substituído) | Obsoleto |
| `0001-HISTORY.md` | Decisões implementadas | Histórico |
| `PENDING.md` | Slices em espera | Ativo |

---

## 0001-HISTORY.md — Implementado ✅

Decisões de 24-30/03/2026 que foram implementadas e validadas:

| Categoria | Count | Status |
|-----------|-------|--------|
| Infraestrutura | 3 | ✅ |
| Skills & Agents | 5 | ✅ |
| Streaming & Multimodal | 3 | ✅ |
| Smart Router & LLM | 3 | ✅ |
| Segurança & Defesa | 3 | ✅ |
| Telegram & Comunicação | 2 | ✅ |
| Voice & TTS | 3 | ✅ |
| Browser & MCP | 3 | ✅ |
| Data & Storage | 2 | ✅ |
| Jarvis & Autonomous | 3 | ✅ |
| Computer Use | 3 | ✅ |
| **Total** | **33** | **100% ✅** |

---

## PENDING.md — Pendente 📋

| Prioridade | Count | Status |
|------------|-------|--------|
| P1 Crítico | 3 | 🔴 |
| P2 Alto | 4 | 🟡 |
| P3 Médio | 3 | 🟢 |
| **Total** | **10** | **📋** |

### P1 Crítico
- Computer Use E2E (BUA-style)
- OS Native God Mode
- Jarvis Voice + Computer

---

## Model Stack Policy (Ativo)

```
Proibidos:
  gemma3:27b, gemma3:12b, groq/whisper, deepseek/chat, bge-m3

Cascade:
  Nível 0: qwen3.5 + faster-whisper-v3 (local)
  Nível 1: minimax-01:free, gemini-2.0-flash (cloud free)
  Nível 2: glm-5, minimax-2.7, kimi-2.5 (paid)
  Nível 3: aurelia-top, aurelia-audio (SOTA)
```

---

## Referências

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../governance/REPOSITORY_CONTRACT.md)
- [SKILL-CATALOG.md](../governance/SKILL-CATALOG.md)

---

**Última atualização:** 2026-03-31
