---
name: Taskmaster Index
description: Mapa operacional centralizado de todas as pendências críticas, organizado por onda de prioridade e status
type: operational-index
last_updated: 2026-03-19T23:59:00-03:00
---

# 🎯 Taskmaster Index — Aurelia Slice Execution Map

Este índice consolida todas as **12 ADRs nonstop** do plano Aurelia em uma visão única, organizada por:
- **Onda de execução** (Wave: 1 crítico → 4 suportivo)
- **Status real** (Proposto, Em execução, Aceito)
- **Bloqueadores e evidência**
- **Comandos de resumo** para cada motor (Codex, Claude, Antigravity)

---

## 📊 Resumo Executivo

| Onda | Slice | Crítico | Status | Bloqueado | Handoff |
| :---: | :---: | :---: | :---: | :---: | :---: |
| **Onda 1** | Slice 7 (Wake E2E) | 🔴 | 🔵 Em execução | Prova humana | Codex |
| **Onda 1** | Slice 10 (Voz oficial) | 🔴 | 🔵 Em execução | Voice ID | Codex |
| **Onda 2** | Slice 4 (Antigravity handoff) | 🔴 | 🔵 Em execução | E2E validation | Claude |
| **Onda 2** | Slice 2 (Safe login) | 🔴 | 🔵 Em execução | Smoke tests | Claude |
| **Onda 3** | Slice 11 (Agent bus) | 🟡 | 🟡 Proposto | Design + Postgres | Claude |
| **Onda 4** | Slice 3 (Desktop fallback) | 🟡 | 🔵 Em execução | Integration | Claude |
| **Support** | Voice capture runtime | 🟢 | 🔵 Em execução | Integration | Codex |
| **Support** | State/memory runtime | 🟢 | 🔵 Em execução | Sync | Codex |
| **Support** | Deploy gateway voice | 🟢 | 🔵 Em execução | Rollout | Codex |
| **Support** | Extensions governance | 🟢 | ✅ Aceito | — | — |
| **Support** | Media/voice transcripts | 🟢 | 🔵 Em execução | Integration | Claude |
| **Support** | Offline homelab Qdrant | 🟢 | 🟡 Proposto | Manual + SQLite | Claude |

---

## 🌊 ONDA 1 — Voz e Experiência de Resposta (Crítico)

Todas as três ações de voz executam em paralelo. Conclusão dessa onda **libera a prova humana live** e viabiliza Telegram bidirecional.

### ADR-20260319-voice-e2e-proof-live

**Objetivo:** Provar no environment live o caminho completo **wake word → STT → resposta**.

| Campo | Valor |
| :--- | :--- |
| **Status** | 🔵 Em execução (Phase: validation, 25% completo) |
| **Slice** | Slice 7 |
| **Owner** | Codex |
| **Bloqueador crítico** | Prova humana positiva no headset real |
| **Próximas ações** | 1. Validar device ALSA definitivo 2. Fazer prova humana 3. Coletar evidência |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-voice-e2e-proof-live.json` |

**Smoke tests:**
```bash
go test ./internal/voice ./cmd/aurelia -count=1
curl -fsS http://127.0.0.1:8484/v1/voice/capture/status
```

**Resumo para Codex:** Continue a prova humana positiva no headset. Capture evidência em `health`, `voice_status` e `voice_events`.

---

### ADR-20260319-aurelia-media-voice

**Objetivo:** Implementar **transcript de mídia** e **voz oficial da Aurelia** com identidade consistente.

| Campo | Valor |
| :--- | :--- |
| **Status** | 🔵 Em execução |
| **Slice** | Slice 10 (Parte 1/2) |
| **Owner** | Claude |
| **Bloqueador crítico** | Voice ID stable (TTS model selection) |
| **Próximas ações** | 1. Selecionar voice model (MiniMax vs Gemini) 2. Ficar com TTS determinístico 3. Integrar com media_processor |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-aurelia-media-voice.json` |

**Comandos de teste:**
```bash
go test ./internal/media ./internal/tts -count=1
curl -fsS http://127.0.0.1:8484/v1/voice/synthesize -d '{"text":"test"}'
```

---

### ADR-20260319-aurelia-authorized-voice-clone

**Objetivo:** Executar **clonagem autorizada** da voz oficial a partir de arquivo local, com **consentimento registrado** e **rollback pronto**.

| Campo | Valor |
| :--- | :--- |
| **Status** | 🔵 Em execução |
| **Slice** | Slice 10 (Parte 2/2) |
| **Owner** | Claude |
| **Bloqueador crítico** | Arquivo de áudio de referência + consentimento legal |
| **Próximas ações** | 1. Obter áudio de referência assinado 2. Integrar com voice clone worker 3. Validar smoke TTS |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-aurelia-authorized-voice-clone.json` |

**Evidência mínima:** Arquivo de consentimento + prova de síntese em logs.

---

## 🌊 ONDA 2 — Orquestração Segura de Browser e Antigravity

Executa **sequencialmente** após Onda 1 estar estabilizada. Envolve browser, login seguro e roteamento leve.

### ADR-20260319-antigravity-handoff-e2e

**Objetivo:** Handoff **fim a fim** com Antigravity: prompt → CLI → resposta → sinal de conclusão (sem retrabalho).

| Campo | Valor |
| :--- | :--- |
| **Status** | 🔵 Em execução |
| **Slice** | Slice 4 |
| **Owner** | Claude |
| **Bloqueador crítico** | E2E validation com prompt real no CLI |
| **Próximas ações** | 1. Testar handoff com resumo de contexto 2. Validar signal de conclusão 3. Executar E2E Antigravity |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-antigravity-handoff-e2e.json` |

**Comandos de teste:**
```bash
./aurelia handoff-test --prompt "Hello" --timeout 30s
./aurelia smoke-antigravity-integration
```

---

### ADR-20260319-browser-safe-login

**Objetivo:** Fluxo de **login guiado seguro** em browser com validação de campos e proteção contra injeção.

| Campo | Valor |
| :--- | :--- |
| **Status** | 🔵 Em execução |
| **Slice** | Slice 2 |
| **Owner** | Claude |
| **Bloqueador crítico** | Smoke browser tests com Playwright |
| **Próximas ações** | 1. Implementar safe_fill para campos de password 2. Validar protocolo de login 3. Testar smoke browser |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-browser-safe-login.json` |

**Teste mínimo:**
```bash
npm test -- browser-safe-login.spec.ts
```

---

## 🌊 ONDA 3 — Swarm Hierárquico da Aurélia (Estrutural)

Depende da conclusão de Ondas 1 e 2. Implementa o **agent bus**, **dashboard** e **assistência entre agentes**.

### ADR-20260319-hierarchical-agent-bus

**Objetivo:** Implementar **agent bus com PostgreSQL** para coordenação de swarm, incluindo:
- Worker claim e lease (com timeout)
- Fila de tarefas
- Memória semântica em Qdrant

| Campo | Valor |
| :--- | :--- |
| **Status** | 🟡 Proposto |
| **Slice** | Slice 11 (Parte 1/3) |
| **Owner** | Claude |
| **Bloqueador crítico** | Schema PostgreSQL + LangGraph reference |
| **Próximas ações** | 1. Design schema (workers, tasks, leases) 2. Implementar worker discovery 3. Validar worker claim + release |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-hierarchical-agent-bus.json` |

**Referência:** `open-agent-supervisor`, `langgraph-supervisor`

**Schema mínimo:**
```sql
CREATE TABLE workers (id UUID, name TEXT, capacity INT, leased_at TIMESTAMP);
CREATE TABLE tasks (id UUID, worker_id UUID, status TEXT, created_at TIMESTAMP);
CREATE TABLE leases (worker_id UUID, task_id UUID, expires_at TIMESTAMP);
```

---

## 🌊 ONDA 4 — Desktop Fallback Seguro

Última onda. Implementa **click**, **typing** e **kill-switch** com validação de cada ação.

### ADR-20260319-desktop-safe-fallback

**Objetivo:** Implementar fallback seguro para desktop com:
- **Safe click:** Validação de posição antes de clicar
- **Safe typing:** Digitação reversível com timeout
- **Kill-switch:** Abort em `N` passos ou `T` segundos

| Campo | Valor |
| :--- | :--- |
| **Status** | 🔵 Em execução |
| **Slice** | Slice 3 |
| **Owner** | Claude |
| **Bloqueador crítico** | Integration tests com Playwright |
| **Próximas ações** | 1. Implementar safe_click 2. Adicionar safe_type reversível 3. Validar kill-switch com timeout |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-desktop-safe-fallback.json` |

**Teste mínimo:**
```bash
npm test -- desktop-safe-*.spec.ts --timeout 10000
```

---

## 🔧 SUPORTE (Enablers — Em Paralelo)

Estas slices executam **em paralelo** com as ondas críticas. Fornecem runtime, state management e deployment.

### ADR-20260319-voice-capture-runtime

**Objetivo:** Integrar **voice capture worker** ao runtime nonstop com persistência de estado.

| Status | 🔵 Em execução |
| :--- | :--- |
| **Owner** | Codex |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-voice-capture-runtime.json` |

---

### ADR-20260319-state-memory-runtime

**Objetivo:** Persistência de **gateway state** e **transcripts locais** em SQLite com sync.

| Status | 🔵 Em execução |
| :--- | :--- |
| **Owner** | Codex |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-state-memory-runtime.json` |

---

### ADR-20260319-deploy-gateway-voice

**Objetivo:** Rollout contínuo em `/home/will/aurelia-24x7` com health checks.

| Status | 🔵 Em execução |
| :--- | :--- |
| **Owner** | Codex |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-deploy-gateway-voice.json` |

---

### ADR-20260319-extensions-governance

**Objetivo:** Política final de extensões para Aurelia (plugin architecture).

| Status | ✅ Aceito |
| :--- | :--- |
| **Owner** | — |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-extensions-governance.json` |

---

### ADR-20260319-offline-homelab-manual-qdrant

**Objetivo:** Manual offline do homelab para recuperação semântica com Qdrant 9B e SQLite (sem web).

| Status | 🟡 Proposto |
| :--- | :--- |
| **Owner** | Claude |
| **Bloqueador crítico** | Manual + SQLite schema |
| **Arquivo JSON** | `docs/adr/taskmaster/ADR-20260319-offline-homelab-manual-qdrant.json` |

**Teste mínimo:**
```bash
sqlite3 :memory: "SELECT COUNT(*) FROM vector_cache"
curl -fsS http://127.0.0.1:6333/health  # Qdrant
```

---

## 📋 Matriz de Handoff Rápido

Use esta tabela para **retomar trabalho** em qualquer slice:

| ADR ID | Motor | Owner Engine | Resume Prompt | Last Updated |
| :--- | :--- | :--- | :--- | :--- |
| voice-e2e-proof-live | Codex | codex | Faça a prova humana no headset | 2026-03-19 |
| aurelia-media-voice | Claude | claude | Selecione voice model e integre | 2026-03-19 |
| aurelia-authorized-voice-clone | Claude | claude | Obtenha áudio de referência + consentimento | 2026-03-19 |
| antigravity-handoff-e2e | Claude | claude | Teste handoff E2E com resumo de contexto | 2026-03-19 |
| browser-safe-login | Claude | claude | Implemente safe_fill para password fields | 2026-03-19 |
| hierarchical-agent-bus | Claude | claude | Design schema PostgreSQL + worker claim | 2026-03-19 |
| desktop-safe-fallback | Claude | claude | Implemente safe_click + safe_type reversível | 2026-03-19 |
| voice-capture-runtime | Codex | codex | Integre voice worker ao runtime | 2026-03-19 |
| state-memory-runtime | Codex | codex | Persista state em SQLite com sync | 2026-03-19 |
| deploy-gateway-voice | Codex | codex | Rollout para aurelia-24x7 | 2026-03-19 |
| extensions-governance | — | — | ✅ Completo | 2026-03-19 |
| offline-homelab-manual-qdrant | Claude | claude | Crie manual offline + SQLite schema | 2026-03-19 |

---

## 🎬 Como Usar Este Índice

1. **Para retomar uma slice:**
   - Procure o ADR ID na tabela
   - Use a "Resume Prompt" como base para o prompt de retomada
   - Abra o JSON correspondente em `docs/adr/taskmaster/`

2. **Para entender prioridades:**
   - Leia de cima para baixo: Onda 1 → Onda 2 → Onda 3 → Onda 4
   - Onda 1 + 2 são **críticas para MVP**
   - Onda 3 + 4 são **estruturais e de segurança**

3. **Para sincronizar contexto:**
   - Após conclusão de cada slice, execute: `sync-ai-context`
   - Atualize este índice com novo status e timestamp

4. **Para escalar bloqueadores:**
   - Procure a coluna "Bloqueador crítico"
   - Reporte ao proprietário do motor (Codex/Claude/Antigravity)

---

## 📍 Referências

- **Autoridade:** `/home/will/aurelia/AGENTS.md`
- **Backlog oficial:** `/home/will/aurelia/docs/adr/PENDING-SLICES-20260319.md`
- **Plano executivo:** `/home/will/aurelia/plan.md`
- **Arquivos JSON:** `/home/will/aurelia/docs/adr/taskmaster/`

---

**Última atualização:** 2026-03-19
**Próxima revisão:** Após conclusão de Onda 1 (Prova humana live)
