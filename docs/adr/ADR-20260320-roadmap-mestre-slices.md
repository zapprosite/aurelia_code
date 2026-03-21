---
title: Aurelia — Roadmap Mestre de Slices
status: active
date: 2026-03-20
decision-makers: [humano, aurelia]
supersedes: todos-ADRs-20260319, todos-blueprints-20260319
---

# ADR-20260320: Aurelia — Roadmap Mestre de Slices

## Contexto

Hardware alvo: AMD Ryzen 9 7900X + RX 7900 XTX (24 GiB VRAM, ROCm), 32 GiB DDR5.
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

### 🟠 [P3] Swarm Hierárquico Go — Onda 3

**Status:** 📋 Proposto — branch `agent-to-agent` pronta para cherry-pick

Arquivos para cherry-pick:
- `internal/voice/kokoro_client.go` + `kokoro_test.go`
- `scripts/simulate_swarm_2026.go`
- `supabase/migrations/20260320_swarm_tables.sql`

Tarefas:
- [ ] Cherry-pick joias de `agent-to-agent`
- [ ] Implementar `agent_bus` Go em Postgres (schema pronto no SQL)
- [ ] Memória compartilhada via Qdrant entre agentes
- [ ] Traduzir padrões LangGraph → Go nativo (channels + goroutines)

**Aceite:** Múltiplos agentes Go se comunicam via bus sem Python runtime.

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

**Status:** 📋 Proposto

- [ ] Fases 1-4 de Polish-Governance-All
- [ ] Secret-audit no crontab semanal (`0 6 * * 1`)

---

### 🟢 [P10] Dashboard Real-Time Engine (Onda 4)
**Status:** 📋 Proposto
- [ ] Implementar camada de streaming via WebSockets/SSE em Go
- [ ] Substituir `MOCK_FEED` por dados reais do barramento de eventos
- [ ] Status de rede e latência de agentes live

### 🟢 [P11] Cockpit de Comando V1
**Status:** 📋 Proposto
- [ ] Input Global (CMD+K) no Dashboard para envio de prompts
- [ ] Integração com barramento de comando da Aurélia
- [ ] Histórico de comandos locais (browser storage)

### 🟢 [P12] Observabilidade Hardware Pro
**Status:** 📋 Proposto
- [ ] Gráficos dinâmicos de GPU/VRAM (ROCm/NVIDIA)
- [ ] Histórico de temperatura e clock
- [ ] Widget de carga de CPU por agente

### 🟢 [P13] Watchdog integration (Toasts)
**Status:** 📋 Proposto
- [ ] Notificações Push de Self-Healing (sonner/toast)
- [ ] Alertas de instabilidade em containers Docker
- [ ] Feedback visual de intervenções autônomas

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

---

## Ordem de Execução

```
P1 → Voz E2E (prova humana)
P2 → Antigravity Handoff (round-trip)
P3 → Swarm Go (cherry-pick agent-to-agent)
P4 → Memória offline Qdrant
P5 → Browser Login Seguro
P6 → Desktop Fallback
P7 → KeePass Vault (humano)
P8 → Governance Polish
P10 → Dashboard Real-Time Engine
P11 → Cockpit de Comando
P12 → Observabilidade Hardware
P13 → Watchdog/Toasts
```

## Consequências

- ADRs 20260319 individuais removidos — este ADR é a fonte de verdade de slices.
- Blueprints 20260319 removidos — decisões absorvidas aqui.
- Manter: `ADR-20260320-politica-modelos-hardware-vram.md`, `ADR-20260320-plano-mestre-jarvis-local-first.md`.
- Novos slices: novas seções aqui ou `ADR-20260320-*` se detalhe extenso.
