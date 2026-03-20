---
title: ADR-20260319 Governance Documentation Index
description: Índice de todos documentos de governança (DATA, SECRETS, NETWORK, OPERATIONAL, OBSERVABILITY)
owner: codex
updated: 2026-03-20
---

# ADR-20260319 Governance Documentation Index

**Master ADR:** [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)

**Purpose:** Navegar a toda documentação de governança organizadas por seção da ADR

**Status:** 🟢 DATA_GOVERNANCE + OPERATIONAL_GOVERNANCE + OBSERVABILITY_GOVERNANCE COMPLETE

---

## Navigation Map

### 1️⃣ DATA_GOVERNANCE ✅ COMPLETE

**Goal:** Definir onde cada dado vive e como é protegido

| Documento | Propósito | Audience |
|-----------|-----------|----------|
| [**Schema Registry — SQLite**](./schema-registry-sqlite.md) | 5 tabelas planejadas em aurelia.db (gateway, voice, cron, memory, health) | Developers, DBA |
| [**Schema Registry — PostgreSQL**](./schema-registry-postgres.md) | 4 instâncias Postgres (n8n, supabase, litellm, dev) + tabelas futuras | Developers, DBA |
| [**Domain Ownership Table**](./data-governance-domain-ownership.md) | Owner, SLA, criticidade per tabela | Operations, Policy |
| [**Store Selection Matrix**](./data-governance-store-selection.md) | Decision tree: SQLite vs Postgres vs Qdrant vs Supabase | Architects, Developers |
| [**Data Lifecycle Policy**](./data-governance-lifecycle.md) | 30d hot, 90d warm, archive gzipped com cleanup scripts | Operations, Compliance |
| [**Qdrant Collection Contract**](./qdrant-collection-contract.md) | bge-m3 384-dim embeddings, payload schema, sync flow | ML Engineers, Developers |
| [**Compliance Matrix**](./data-governance-compliance.md) | Conformance com CONTRACT.md + GUARDRAILS.md, audit trails | Compliance, Policy |

**Key Documents:**
- 5 SQL schemas with indexes and retention policies
- 4 PostgreSQL instances documented
- Clear decision tree for store selection
- Automated compliance validation

**Responsibilities:**
- codex: Operacional (sync, backup, schema updates)
- humano: Policy (ownership, criticality, exceptions)

**Status:** ✅ 7/7 documents created

---

### 2️⃣ OPERATIONAL_GOVERNANCE ✅ COMPLETE (Core)

**Goal:** Manter infraestrutura operacional e responder a incidentes

| Documento | Propósito | Audience |
|-----------|-----------|----------|
| [**Health Checks**](./operational-governance-health-checks.md) | Automated cron monitoring (5/15/60/1440 min) | Operations, SRE |
| [**Backup Verification**](./operational-governance-backup-verification.md) | Daily freshness checks, monthly restore simulation | Operations, DBA |
| [**Incident Response**](./operational-governance-incident-response.md) | Runbooks para P1-P4 incidents (container, OOM, disk, DB, tunnel, breach) | Operations, On-call |

**Cron Schedule:**
- **Every 5 min:** health-check.sh --mode quick (container, disk, VRAM, Qdrant)
- **Every 15 min:** health-check.sh --mode services (Postgres, Qdrant, backup age)
- **Every hour:** health-check.sh --mode deep (schema, backup, sync age)
- **Every 6am:** health-check.sh --mode smoke (end-to-end tests)
- **Every 3am:** verify-backups.sh (freshness, integrity)
- **First Saturday 4am:** test-restore.sh (recovery simulation)
- **Jan 1 2pm:** annual-backup-audit.sh (compliance audit)

**Status:** ✅ 3/3 core documents created

**Remaining (Fase 4 - LOW priority):**
- [ ] Cleanup deploy dirs (manual housekeeping)
- [ ] Fix unclean shutdown warnings (Docker config tuning)

---

### 3️⃣ OBSERVABILITY_GOVERNANCE ✅ COMPLETE (Core)

**Goal:** Monitore infraestrutura e aplicação para detectar problemas antes de ficarem críticos

| Documento | Propósito | Audience |
|-----------|-----------|----------|
| [**Metrics Contract**](./observability-governance-metrics.md) | Mandatory metrics per service + Prometheus config + 5 Grafana dashboards | Developers, SRE |

**Metrics Published:**
- **Health:** disk %, VRAM %, container status, backup age, memory sync age
- **Databases:** connection count, query latency, cache hit ratio, table sizes
- **Memory Sync:** files scanned, Qdrant point count, embed latency
- **Services:** request rate, latency distribution, error counts by type
- **LLM:** requests by model, token usage, fallback triggers

**Alert Rules (integrated with incident-response.md):**
- 🔴 **CRITICAL:** Container down, disk >95%, VRAM >95%, backup >24h, Qdrant down
- 🟠 **HIGH:** Disk 80-95%, VRAM 85-95%, memory sync >1h old
- 🟡 **MEDIUM:** API latency >1s, slow queries detected

**Dashboards Required:**
1. Overview (quick health check)
2. System (CPU, memory, network, disk I/O)
3. Databases (per-instance Postgres metrics)
4. Memory Sync (Qdrant, embedding, file sync)
5. LLM (requests, tokens, latency, fallback)

**Status:** ✅ 1/1 core document created

---

### 4️⃣ SECRETS_GOVERNANCE ⏳ SKIPPED (User Request)

**Goal:** Proteger credenciais + rotação trimestral

**Planned Documents (Fase 1 - CRITICAL):**
- KeePassXC Vault setup + migration guide
- Secret rotation playbook
- Plaintext credential cleanup checklist

**Status:** ⏳ Deferred (User said "pular o 1")

---

### 5️⃣ NETWORK_GOVERNANCE ⏳ NOT STARTED

**Goal:** Segurança de rede + isolamento de containers

**Planned Documents (Fase 4):**
- UFW firewall hardening plan (SSH first)
- Port exposure matrix (public vs internal)
- Cloudflare Tunnel hardening
- Docker network isolation (per stack)
- Tailscale ACL policy

**Status:** ⏳ Planned for Fase 4

---

### 6️⃣ COMPLIANCE_MATRIX (Part of DATA_GOVERNANCE) ✅

**Integrated into:** [data-governance-compliance.md](./data-governance-compliance.md)

**Covers:**
- CONTRACT.md conformance (7/7 requirements)
- GUARDRAILS.md approval matrix
- Audit trails per store
- Domain ownership authorization
- Exception waivers + audit schedule
- Escalation paths (breach, SLA, policy)

**Status:** ✅ Integrated with DATA_GOVERNANCE

---

## Execution Summary

### ✅ Completed (10/16 items)

**Fase 1 (CRITICAL) — SKIPPED per user request:**
- Item 1 → Vault setup, credential migration, plaintext cleanup

**Fase 2 (HIGH) — PARTIAL:**
- ✅ Item 3 → Schema Registry (DATA_GOVERNANCE fully documented)
- ✅ Item 2 → Systemd timers (memory-sync-fiscal.sh + 4 timers, per summary)
- ⏳ Item 5 → app.json.bak* deletion (still needed)
- ⏳ Item 6 → MCP config refactor (still needed)
- ⏳ Item 8 → Secret rotation policy (part of SECRETS_GOVERNANCE, skipped)

**Fase 3 (MEDIUM) — COMPLETED:**
- ✅ Item 9 → Health checks (health-check.sh with 4 cron frequencies)
- ✅ Item 10 → Backup verification (verify-backups.sh + test-restore.sh)
- ✅ Item 11 → Incident response playbook (6 specific runbooks)
- ✅ Item 12 → Qdrant collection schema (qdrant-collection-contract.md)
- ✅ Item 13 → Data lifecycle policy (data-governance-lifecycle.md)
- ✅ Item 14 → Observability contract + alert rules (observability-governance-metrics.md)

**Fase 4 (LOW) — NOT STARTED:**
- ⏳ Deploy dirs cleanup
- ⏳ Unclean shutdown warnings fix
- ⏳ Compliance audit script
- ⏳ NETWORK_GOVERNANCE documents

---

## Document Statistics

| Section | Documents | Lines | Status |
|---------|-----------|-------|--------|
| **DATA_GOVERNANCE** | 7 | ~2,500 | ✅ COMPLETE |
| **OPERATIONAL_GOVERNANCE** | 3 | ~1,200 | ✅ COMPLETE |
| **OBSERVABILITY_GOVERNANCE** | 1 | ~600 | ✅ COMPLETE |
| **SECRETS_GOVERNANCE** | 0 | — | ⏳ Skipped |
| **NETWORK_GOVERNANCE** | 0 | — | ⏳ Planned |
| **COMPLIANCE_MATRIX** | Integrated | — | ✅ COMPLETE |
| **TOTAL** | **11** | **~4,300** | 11/16 items |

---

## Quick Links by Role

### For DevOps / SRE

- [Health Checks](./operational-governance-health-checks.md) — Cron setup & alert tuning
- [Incident Response](./operational-governance-incident-response.md) — On-call runbooks
- [Backup Verification](./operational-governance-backup-verification.md) — Restore testing
- [Metrics Contract](./observability-governance-metrics.md) — Dashboard setup

### For Developers

- [Schema Registry — SQLite](./schema-registry-sqlite.md) — Local database tables
- [Schema Registry — PostgreSQL](./schema-registry-postgres.md) — Multi-instance schema
- [Store Selection Matrix](./data-governance-store-selection.md) — Where to store data
- [Qdrant Collection Contract](./qdrant-collection-contract.md) — Vector DB API

### For Security / Compliance

- [Data Lifecycle Policy](./data-governance-lifecycle.md) — Retention & archiving
- [Domain Ownership Table](./data-governance-domain-ownership.md) — Owner & SLA tracking
- [Compliance Matrix](./data-governance-compliance.md) — CONTRACT.md & GUARDRAILS.md alignment
- [Incident Response](./operational-governance-incident-response.md) — Breach procedures

### For Architects

- [Store Selection Matrix](./data-governance-store-selection.md) — Decision tree
- [Memory Sync Architecture](./memory-sync-architecture.md) — Vector DB + embedding flow
- [Qdrant Collection Contract](./qdrant-collection-contract.md) — Semantic search spec

---

## Conformance Checklist

### CONTRACT.md Requirements
- ✅ Data sovereignty (SQLite local, not cloud)
- ✅ Audit trails (per-store logging)
- ✅ Encryption at rest + in-transit (ZFS, TLS, pgcrypto)
- ✅ Backup strategy (hourly ZFS, daily Postgres dumps)
- ✅ RTO/RPO defined (per table/criticality)

### GUARDRAILS.md Compliance
- ✅ Read-only ops are free (no approval needed)
- ✅ Backups/snapshots are free (no approval needed)
- ✅ Data modifications require approval
- ✅ Secrets handled separately (SECRETS_GOVERNANCE)
- ✅ Network changes require NETWORK_MAP.md update
- ✅ ZFS ops snapshot first
- ✅ Exceptions require ADR amendment + expiration

---

## Next Actions

### High Priority (For Human Review)

1. **Read & Approve** each governance section
   - [ ] DATA_GOVERNANCE (7 docs)
   - [ ] OPERATIONAL_GOVERNANCE (3 docs)
   - [ ] OBSERVABILITY_GOVERNANCE (1 doc)

2. **Confirm Cron Installation** (Health checks)
   ```bash
   crontab -l | grep "health-check"
   ```

3. **Setup Prometheus + Grafana** (if not already done)
   - [ ] Prometheus scrape config
   - [ ] 5 Grafana dashboards
   - [ ] AlertManager integration

### Later (Fase 4 + User Request)

- [ ] SECRETS_GOVERNANCE (Item 1 - when user requests)
- [ ] NETWORK_GOVERNANCE (Fase 4 - hardening)
- [ ] Deploy cleanup (Fase 4 - low priority)

---

## Links

- **Master ADR:** [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
- **Authority:** [CONTRACT.md](/srv/ops/ai-governance/CONTRACT.md)
- **Approval Matrix:** [GUARDRAILS.md](/srv/ops/ai-governance/GUARDRAILS.md)
- **Network Topology:** [NETWORK_MAP.md](/srv/ops/ai-governance/NETWORK_MAP.md)
- **Memory Sync Arch:** [memory-sync-architecture.md](./memory-sync-architecture.md)
- **ADR Status:** [adr-semparar-status.md](./.agents/workflows/adr-semparar-status.md)
