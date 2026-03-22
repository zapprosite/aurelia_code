---
title: Aurelia — Roadmap Mestre de Slices
status: active
date: 2026-03-20
decision-makers: [humano, aurelia]
supersedes: todos-ADRs-20260319, todos-blueprints-20260319
---

# ADR-20260320: Aurelia — Roadmap Mestre de Slices

## Contexto

Hardware alvo: AMD Ryzen 9 7900X + NVIDIA RTX 4090 (24 GiB VRAM), 32 GiB DDR5.
Stack: Go, Telegram, Ollama local, Qdrant, Supabase, Postgres, SQLite.

---

## Decisões Arquiteturais Permanentes

| Área | Decisão Fechada |
|------|----------------|
| Autoridade | AGENTS.md soberano → Aurélia → agentes |
| Modelos | gemma3:12b residente, qwen3.5:9b fallback, bge-m3 embedding |
| Voz | Groq STT (remoto) + Kokoro TTS (CPU local) |
| Memória | Qdrant + Postgres + SQLite nativo Go, sem Python runtime |
| LangGraph | Padrão de referência apenas — implementação Go via canais nativos |
| Browser | Playwright para login guiado; desktop (xdotool) apenas como fallback |
| Segurança | secrets.env centralizado, chmod 600 + secret-audit.sh pré-push |
| Commits | `--no-verify` padrão de agentes (hooks lentos); CI valida |
| Sudo | sudo=1 habilitado Tier C com log obrigatório |

**Detalhe de modelos e VRAM:** ver `ADR-20260320-politica-modelos-hardware-vram.md`
**Arquitetura JARVIS completa:** ver `ADR-20260320-plano-mestre-jarvis-local-first.md`

---

## Slices por Prioridade

### 🔴 [P1] Voz E2E — Onda 1

**Status:** 🔄 Em progresso (90%)
**Dependências:** hardware AMD + Groq API key + Kokoro instalado

- [ ] Validar E2E: wake word → Groq STT → gemma3:12b → resposta Telegram
- [ ] Provar com evidência humana (log estruturado + áudio)
- [ ] Validar `GET /v1/voice/status` e métricas live
- [ ] Fechar rollback claro de voz

**Aceite:** Bot responde por voz em PT-BR a partir de comando falado sem intervenção.

---

### 🔴 [P2] Antigravity Handoff E2E — Onda 2

**Status:** 🔄 Em progresso (85%)

- [ ] Fechar handoff Aurélia → Antigravity → resultado registrado
- [ ] Validar task `light` → prompt automático → execução → retorno estruturado
- [ ] Login guiado seguro no browser via Playwright com kill-switch

**Aceite:** Aurélia delega tarefa ao Antigravity e recebe resultado de volta estruturado.

---

### 🟠 [P3] Agent-to-Agent Communication — Go Native

**Status:** ✅ Concluído (v6.0-handoff)

- [x] HandoffTool nativo em Go (`internal/agent/handoff.go`)
- [x] MasterTeamService com delegação síncrona
- [x] Agent Loop com interrupção para handoff
- [x] Dashboard SSE com eventos de handoff
- [x] Simulação E2E aprovada

**Aceite:** Agentes Go se comunicam via handoff nativo sem Python runtime.

---

### 🟠 [P4] Memória Offline Qdrant

**Status:** 🔄 Parcialmente implementado

- [ ] Validar crons memory-sync-fiscal.sh nos 4 intervalos (5/15/60/1440 min)
- [ ] Confirmar indexação de Go code history no Qdrant
- [ ] Endpoint Telegram: `aurelia search "X"` retorna contexto semântico

**Aceite:** Aurelia busca no codebase sem Gemini API.

---

### 🟡 [P5] Browser Login Seguro

**Status:** 🔄 Em progresso (90%)

- [ ] Finalizar fluxo Playwright com confirmação humana antes de submit
- [ ] Kill-switch explícito + limite de passos

**Aceite:** Login em site-alvo com aprovação humana por passo.

---

### 🟡 [P6] Desktop Fallback Seguro

**Status:** 🔄 Em progresso (60%)

- [ ] Click seguro com limite de passos (xdotool/wmctrl)
- [ ] Kill-switch de desktop

**Aceite:** Ação de desktop com aprovação explícita.

---

### 🟡 [P7] KeePassXC Vault

**Status:** ⏳ Aguardando humano (deadline: 2026-03-27)

- [ ] `bash scripts/setup-keepassxc-vault.sh`
- [ ] Masterkey em hardware token ou TPM
- [ ] systemd service com `EnvironmentFile` correto

**Aceite:** Aurelia inicia sem plaintext secrets em disco.

---

### 🟢 [P8] Governance Polish

**Status:** 🔄 Em progresso

- [x] Secrets.env consolidado com systemd EnvironmentFile
- [x] Env Overrides implementados em `internal/config/config.go`
- [ ] Secret-audit no crontab semanal (`0 6 * * 1`)
- [ ] KeePassXC vault (deadline: 2026-03-27)

---

### ✅ [P9] Dashboard ULTRATRINK Premium (v6.0)
**Status:** ✅ Concluído
- [x] Refatoração modular React (componentes `shadcn`)
- [x] Sistema de tabs (Timeline, Squad, Brain, Roadmap)
- [x] Micro-animações com `framer-motion`

### ✅ [P10] Dashboard Real-Time Engine (v6.1-realtime)
**Status:** ✅ Concluído
- [x] EventBus + endpoint `/api/events` em Go
- [x] Instrumentação de `agent/loop.go` e `telegram/input_pipeline.go`
- [x] Frontend conectado via `EventSource` (React hooks)
- [x] Build de produção embarcado no binário Go

### 🟢 [P11] Cockpit de Comando V1
**Status:** 📋 Proposto
- [ ] Input Global (CMD+K) no Dashboard para envio de prompts
- [ ] Integração com barramento de comando da Aurélia
- [ ] Histórico de comandos locais (browser storage)

### 🟢 [P12] Observabilidade Hardware Pro
**Status:** 📋 Proposto
- [ ] Gráficos dinâmicos de GPU/VRAM (NVIDIA RTX 4090)
- [ ] Histórico de temperatura e clock
- [ ] Widget de carga de CPU por agente

### 🟢 [P13] Watchdog integration (Toasts)
**Status:** 📋 Proposto
- [ ] Notificações Push de Self-Healing (sonner/toast)
- [ ] Alertas de instabilidade em containers Docker
- [ ] Feedback visual de intervenções autônomas

---

### 🔴 [P14] Aurelia Autonomous Engineering — ULTRATRINK
**Status:** 🔄 Em progresso
**ADR:** [ADR-20260321-aurelia-autonomous-engineering.md](ADR-20260321-aurelia-autonomous-engineering.md)
**Taskmaster:** [JSON](taskmaster/ADR-20260321-aurelia-autonomous-engineering.json)

As 7 capacidades que transformam a Aurélia em engenheiro de software autônomo:

- [x] **Sub-1:** Tool Introspection System (`internal/agent/tool_catalog.go`) 🔄
- [x] **Sub-2:** Execution DNA System (`internal/persona/execution_dna.go`)
- [x] **Sub-3:** Planning Loop PREV (`internal/agent/planner.go` e `loop.go`)
- [x] **Sub-4:** Codebase Symbol Map (`internal/agent/codebase_map.go`)
- [x] **Sub-5:** Semantic Skill Router (`internal/skill/semantic_router.go`)
- [ ] **Sub-6:** Dashboard Cockpit / CMD+K (`frontend/src/components/dashboard/`)
- [x] **Sub-7:** Memory Context Assembler (`internal/memory/context_assembler.go`)

**Aceite:** A Aurélia recebe uma tarefa complexa, gera um plano estruturado, executa ferramenta por ferramenta com consciência do codebase, verifica o resultado e faz rollback se necessário — sem intervenção humana.

---

### ✅ Concluídos

| Slice | Evidência |
|-------|-----------|
| Autoridade AGENTS.md | 4 adaptadores atualizados (CLAUDE/CODEX/GEMINI/OPENCODE) |
| Gateway e Roteamento | enforcement + breaker + Prometheus |
| Pipeline de Voz | spool + heartbeat + STT Groq + CLI enqueue |
| Captura de Microfone | OpenWakeWord + VAD + buffer |
| TTS Telegram | Kokoro CPU, voz oficial, clonagem MiniMax |
| Memória e Health | bge-m3, SQLite persistido, /health real |
| Sync AI Context | crons ativos, regra de slice |
| Repositório Template | padronizado, 4 agentes, sudo=1 |
| Migração ADRs | plan.md + MODEL.md → ADR-20260320-* |
| Dashboard ULTRATRINK | V2 Modular React, Glassmorphism, Squad Grid (v6.0) |
| [P3] Agent-to-Agent Go Native | HandoffTool + MasterTeamService + SSE (v6.0-handoff) |
| [P9] Dashboard Premium | Componentes shadcn + tabs + framer-motion (v6.0) |
| [P10] Real-Time SSE Engine | EventBus + /api/events + React hooks (v6.1-realtime) |
| [P-] Unificação Cross-Model Skills | Symlinks centralizados em `.agents/skills/` |
| [P-] Env Overrides Config | Variáveis de ambiente sobrepondo config.json (21/03/2026) |

---

## Ordem de Execução

```
P1 → Voz E2E (prova humana)
P2 → Antigravity Handoff (round-trip)
P4 → Memória offline Qdrant
P5 → Browser Login Seguro
P6 → Desktop Fallback
P7 → KeePass Vault (humano)
P8 → Governance Polish (🔄 em progresso)
P11 → Cockpit de Comando
P12 → Observabilidade Hardware
P13 → Watchdog/Toasts
```

## Consequências

- ADRs 20260319 individuais removidos — este ADR é a fonte de verdade de slices.
- Blueprints 20260319 removidos — decisões absorvidas aqui.
- Manter: `ADR-20260320-politica-modelos-hardware-vram.md`, `ADR-20260320-plano-mestre-jarvis-local-first.md`.
- Novos slices: novas seções aqui ou `ADR-20260320-*` se detalhe extenso.
