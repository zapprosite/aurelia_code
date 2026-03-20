---
title: Domain Ownership Table
description: Mapeamento de proprietário, SLA, criticidade por tabela em SQLite, PostgreSQL, Qdrant
owner: codex
updated: 2026-03-20
---

# Domain Ownership Table

**Purpose:** Propriedade clara, SLA e criticidade de cada store de dados
**Authority:** ADR-20260319-Polish-Governance-All / DATA_GOVERNANCE
**Updated:** 2026-03-20

---

## SQLite (`~/.aurelia/data/aurelia.db`) — Source of Truth

| Tabela | Owner (Operação) | Owner (Policy) | Criticidade | SLA | RTO | RPO |
|--------|---|---|---|---|---|---|
| `gateway_route_states` | codex | humano | **HIGH** | 99.5% uptime | 5 min | 1 min |
| `voice_events` | codex | codex | MEDIUM | 95% | 30 min | 5 min |
| `cron_tasks` | codex | codex | MEDIUM | 95% | 1 hour | 10 min |
| `memory_entries` | codex | codex | **HIGH** | 99.5% | 5 min | 0 (must sync) |
| `system_health` | codex (auto) | codex | MEDIUM | 90% | 2 hours | 15 min |

---

## PostgreSQL (4 Instâncias)

### n8n (Automations DB)

| Tabela | Owner (Operação) | Owner (Policy) | Criticidade | SLA | Retenção |
|--------|---|---|---|---|---|
| `workflows` | codex | humano | **CRITICAL** | 99.9% | Forever |
| `workflow_history` | codex | codex | **HIGH** | 99.5% | 90d |
| `credentials` | humano | humano | **CRITICAL** | 99.9% | Forever (encrypted) |
| `webhooks` | codex | codex | **HIGH** | 99.5% | Forever |

**RTO:** 15 min | **RPO:** 5 min (hourly snapshots)

---

### supabase-db (Sessions + Mirror)

| Tabela | Owner (Operação) | Owner (Policy) | Criticidade | SLA | Retenção |
|--------|---|---|---|---|---|
| `sessions` | codex | codex | MEDIUM | 95% | 7d active, 30d archive |
| `messages` | codex | codex | MEDIUM | 95% | 30d hot, 90d warm |
| `users` | humano | humano | **HIGH** | 99.5% | Forever |
| `auth.users` | supabase | supabase | **CRITICAL** | 99.9% | Forever (managed) |

**RTO:** 30 min | **RPO:** 10 min

---

### litellm-db (LLM Gateway)

| Tabela | Owner (Operação) | Owner (Policy) | Criticidade | SLA | Retenção |
|--------|---|---|---|---|---|
| `api_keys` | humano | humano | **CRITICAL** | 99.9% | Forever |
| `usage_logs` | codex | codex | MEDIUM | 90% | 30d |
| `model_configs` | codex | codex | **HIGH** | 99.5% | Forever |
| `error_logs` | codex | codex | LOW | 90% | 7d |

**RTO:** 20 min | **RPO:** 5 min

---

### dev (Development)

| Tabela | Owner (Operação) | Owner (Policy) | Criticidade | SLA | Retenção |
|--------|---|---|---|---|---|
| (all) | will | will | LOW | None | None (experimental) |

**RTO:** N/A | **RPO:** N/A (no backup required)

---

## Qdrant (Vector Index)

| Collection | Owner | Criticidade | SLA | Retenção | Embedding Model |
|---|---|---|---|---|---|
| `repository_memory` | codex | **HIGH** | 99% | Forever (derived from Postgres) | bge-m3 (384-dim) |

**Purpose:** Semantic search over memory entries (code history for Aurelia bot)
**RTO:** 1 hour (can regenerate from Postgres) | **RPO:** 0 (derived)

---

## Escalation Matrix

### Severidade de Outage

| Severity | Exemplo | Owner Notificar | Tempo Resposta | Tempo Resolução |
|---|---|---|---|---|
| **CRITICAL** | workflows/credentials indisponíveis | humano (24/7) | 5 min | 15 min |
| **HIGH** | gateway_route_states ou memory_entries down | codex | 15 min | 1 hour |
| **MEDIUM** | voice_events/sessions degraded | codex | 1 hour | 4 hours |
| **LOW** | system_health metrics stale | codex | 4 hours | 24 hours |

**Comunicação:** Telegram bot + Alert rules em Prometheus

---

## Handoff Rules

### Quando escalate para humano?
- Credenciais comprometidas (litellm-db, n8n)
- Autorização policy change necessária
- RTO/RPO não alcançável via automation
- Conflito entre DATA_GOVERNANCE e outro departamento

### Quando escalate para codex?
- Operação de rotina (sync, backup, validation)
- Métricas degradadas mas não critical
- Performance tuning
- Script/automation improvements

---

## Validação

**Command to verify ownership:**
```bash
# Check SQLite actual owners via git blame on schema-registry-sqlite.md
git blame docs/schema-registry-sqlite.md | grep "Tabela\|gateway_route"

# Check Postgres instances via docker ps
docker ps --filter "name=n8n|supabase|litellm"

# Verify Qdrant collections
curl -s http://localhost:6333/collections | jq '.result'
```

---

## Links

- [Schema Registry — SQLite](./schema-registry-sqlite.md)
- [Schema Registry — PostgreSQL](./schema-registry-postgres.md)
- [Data Lifecycle Policy](./data-governance-lifecycle.md)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
