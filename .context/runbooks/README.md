# Runbooks

Procedimentos operacionais e playbooks para automação, monitoramento e resposta a incidentes.

> **Diferença:** Skills são para AI agents. Runbooks são para humanos + automação.

---

## Runbooks Disponíveis

### 1. Governance Polish Phases
**Arquivo:** [`governance-polish-phases.md`](./governance-polish-phases.md)

**Propósito:** Executar ADR-20260319-Polish-Governance-All em 4 fases.

**Fases:**
- **Fase 1 (CRITICAL):** Humano — criar vault KeePassXC, migração secrets, shred plaintext
- **Fase 2 (HIGH):** Codex — deletar backups, refatorar MCP, schema registry
- **Fase 3 (MEDIUM):** Codex — health checks, backup verification, incident playbook
- **Fase 4 (LOW):** Codex — cleanup, UFW, compliance matrix

**Invocação:** `/governance-polish --phase 1 --show-checklist`

**Status:** 🟡 Proposto (Fase 1 pendente)

---

### 2. Memory Sync Fiscal — Crons
**Arquivo:** [`memory-sync-fiscal-cron.md`](./memory-sync-fiscal-cron.md)

**Propósito:** Monitorar e sincronizar memória local → Qdrant + Postgres (automático).

**Crons (4 frequências):**
- **5 min** — Fast sync (embedding incremental)
- **15 min** — Postgres index (atualizar indices)
- **6am diária** — Validate (integrity check + relatório)
- **2am segunda** — Compact (cleanup + otimização)

**Script:** `bash scripts/memory-sync-fiscal.sh --mode [fast|postgres-index|validate|compact]`

**Invocação:** Configurar systemd timers (no runbook)

**Status:** 🟡 Proposto (systemd timer config pendente)

---

## Como Usar Runbooks

### 1. Ler o Runbook
```bash
cat .context/runbooks/governance-polish-phases.md
```

### 2. Identificar a Fase/Seção
Cada runbook tem seções claras com checklists ou scripts.

### 3. Executar Ações
Seguir o checklist ou comandos listados.

### 4. Monitorar
Verificar logs (geralmente em `~/.aurelia/logs/`) para validar sucesso.

---

## Integração com Skills

- **Skill** [`/governance-polish`](../skills/governance-polish/SKILL.md) ↔ **Runbook** `governance-polish-phases.md`
- **Skill** [`/memory-sync-vector-db`](../skills/memory-sync-vector-db/SKILL.md) ↔ **Runbook** `memory-sync-fiscal-cron.md`

Skills **orquestram**, runbooks **executam**.

---

## Criando Novos Runbooks

1. Criar arquivo: `.context/runbooks/my-runbook.md`
2. Estrutura:
   ```markdown
   # Runbook: My Runbook

   **Status:** 🟡 Proposto
   **Responsável:** [humano|codex|gemini|aurelia]
   **Logs:** ~/.aurelia/logs/my-runbook.log

   ---

   ## O que é?
   [Descrição breve]

   ---

   ## Checklist / Passos
   [Instruções passo a passo]

   ---

   ## Troubleshooting
   [Casos de erro + solução]
   ```

3. Registrar em `README.md` (este arquivo)

---

## Monitoramento

### Logs Centralizados
```bash
tail -f ~/.aurelia/logs/*.log
```

### Métricas
```bash
cat ~/.aurelia/metrics/memory-sync.prom  # Prometheus format
```

### Dashboard (Grafana)
- URL: http://localhost:3000 (se configurado)
- Dashboards: 1 por stack
