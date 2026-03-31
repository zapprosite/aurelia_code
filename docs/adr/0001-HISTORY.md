# ADR 0001-HISTORY: Decisões Históricas ✅

**Status:** Arquivado (implementado e validado)
**Data:** 2026-03-24 a 2026-03-30
**Última atualização:** 2026-03-31

---

## Infraestrutura ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| ZFS + Docker | Pool `tank` em nvme0n1, 12 containers | ✅ Implementado |
| Git Cleanup | Limpeza 10k+ arquivos (2.4GB) | ✅ Concluído |
| Claude Code Installer | NPM → binário nativo | ✅ Migrado |

---

## Skills & Agents ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Master Skill Global | Orquestração centralizada | ✅ Implementado |
| Voice Capture | Pipeline Kokoro + Whisper | ✅ Pronto |
| Runtime Governance | Enforcement canônico no Qdrant | ✅ Ativo |
| Team Orchestration | Transparência entre agentes | ✅ Honesto |
| ADR Semparar | Workflow slices nonstop | ✅ Workflow |

---

## Streaming & Multimodal ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Sovereign Streaming | Pipeline streaming soberano | ✅ Implementado |
| VAD Monitor | Voice Activity Detection | ✅ Contínuo |
| Multimodal GPU | GPU optimization | ✅ Otimizado |

---

## Smart Router & LLM ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Smart Router LiteLLM | Cascade: Qwen→MiniMax→GLM | ✅ Operacional |
| Rate Limiting | Smart scheduler | ✅ Ativo |
| Fallback Gateway | Chain: local→free→paid | ✅ Testado |

---

## Segurança & Defesa ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Zero Hardcode | Segredos nunca em código | ✅ Policy |
| Porteiro Sentinel | Defesa Prompt Injection | ✅ Ativo |
| Segurança Sentinel | Hardening | ✅ Consolidado |

---

## Telegram & Comunicação ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Telegram Formatter SOTA | JSON → Markdown 2026 | ✅ Industrializado |
| Multi-Bot | Dashboard + Personas | ✅ Operacional |

---

## Voice & TTS ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Kokoro TTS | Local PT-BR | ✅ Deployado |
| Voxtral (substitui Kokoro) | SOTA TTS BR | ✅ Substituiu |
| Whisper Groq | STT via Groq | ✅ Budget otimizado |

---

## Browser & MCP ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Playwright + Context7 | Automação + docs | ✅ Integrado |
| MCP Go Client | Stagehand Go client | ✅ Feito |
| Rod Browser Layer | CDP stealth mode | ✅ Implementado |

---

## Data & Storage ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Data Stack Contract | Templates SQLite/Qdrant | ✅ Documentado |
| OpenClaw Vault | Isolamento skill vault | ✅ Ishado |

---

## Jarvis & Autonomous ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Jarvis Tutor | Loop wake→STT→TTS | ✅ 24/7 |
| E2E Jarvis Loop | Pipeline completo | ✅ Testado |
| Autonomous Visual | Detecção + OCR | ✅ Ativo |

---

## Computador Use ✅

| ADR | Decisão | Status |
|-----|---------|--------|
| Computer Use E2E | Agent loop autônomo | ✅ Implementado |
| Container Steel | Isolamento browser | ✅ Ativo |
| HITL Safety | Normalized coords | ✅ Confirm gate |

---

## Consolidado em 2026-03-31

| Métrica | Valor |
|---------|-------|
| ADRs históricos | 50 |
| Implementados | 100% |
| Arquivo | `0001-HISTORY.md` |
