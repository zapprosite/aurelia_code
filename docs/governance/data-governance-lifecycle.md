---
title: Data Lifecycle Policy
description: Política de 30d hot, 90d warm, archive gzipped para dados operacionais
owner: codex
updated: 2026-03-20
---

# Data Lifecycle Policy

**Purpose:** Definir retenção, arquiving, e cleanup de dados em SQLite, PostgreSQL, Qdrant, Supabase
**Authority:** ADR-20260319-Polish-Governance-All / DATA_GOVERNANCE
**Updated:** 2026-03-20

---

## Standard Lifecycle — 30/90/Archive

```
Day 0 — Day 29: HOT (30 days)
├─ Location: SQLite or Postgres main table
├─ Purpose: Operational (alerts, queries, analytics)
├─ Backup: Daily snapshots
├─ Access: Fast path, indexable
│
Day 30 — Day 89: WARM (90 days total, 60 days warm)
├─ Location: Postgres archive table or compressed SQLite export
├─ Purpose: Audit, compliance, historical queries
├─ Backup: Weekly snapshots
├─ Access: Slower, scan-based
│
Day 90+: ARCHIVE (Forever, gzipped)
├─ Location: ~/aurelia/backups/ (gzipped SQL dumps)
├─ Purpose: Compliance, legal hold, incident forensics
├─ Backup: Monthly archives to ZFS snapshots
├─ Access: Manual extraction, uncompressed on demand
│
Deletion Exceptions:
├─ CRITICAL tables: Never deleted (workflows, credentials, users)
├─ Compliance-sensitive: 7 year retention (audit logs)
└─ Explicitly labeled "Forever": Retained indefinitely
```

---

## Table-by-Table Retention Policy

### SQLite (`~/.aurelia/data/aurelia.db`)

| Tabela | Hot (days) | Warm (days) | Archive | Notes |
|--------|---|---|---|---|
| `gateway_route_states` | Forever | N/A | N/A | Operational, no rotation |
| `voice_events` | 30 | 90 | gzip | TTS logs, compliance optional |
| `cron_tasks` | 30 | 90 | gzip | Cron execution audit trail |
| `memory_entries` | Forever | N/A | N/A | Local truth before sync |
| `system_health` | 7 | 30 | discard | Metrics, auto-cleanup daily |

**Cleanup Scripts:**
```bash
# Daily cleanup of old system_health (keep 7d)
sqlite3 ~/.aurelia/data/aurelia.db \
  "DELETE FROM system_health WHERE created_at < datetime('now', '-7 days')"

# Monthly archive voice_events (30d → 90d → gzip)
sqlite3 ~/.aurelia/data/aurelia.db \
  "SELECT * FROM voice_events WHERE created_at BETWEEN datetime('now', '-90 days') AND datetime('now', '-30 days')" \
  | gzip > ~/.aurelia/backups/voice_events-$(date +%Y%m%d).sql.gz

# Delete archived voice_events from hot
sqlite3 ~/.aurelia/data/aurelia.db \
  "DELETE FROM voice_events WHERE created_at < datetime('now', '-30 days')"
```

---

### PostgreSQL — n8n Instance

| Tabela | Hot (days) | Warm (days) | Archive | Notes |
|--------|---|---|---|---|
| `workflows` | Forever | N/A | N/A | CRITICAL, never delete |
| `workflow_history` | 30 | 90 | gzip | Execution audit, rotated monthly |
| `credentials` | Forever | N/A | N/A | CRITICAL, encrypted, never delete |
| `webhooks` | Forever | N/A | N/A | Integration config, never delete |

**Lifecycle Job (Monthly):**
```bash
#!/bin/bash
# Move workflow_history to warm (Postgres archive table)
docker exec postgres-n8n psql -U root -d n8n <<EOF
-- Create archive table if not exists
CREATE TABLE IF NOT EXISTS workflow_history_archive AS
  SELECT * FROM workflow_history WHERE 1=0; -- schema only

-- Move 30-90d old rows
INSERT INTO workflow_history_archive
SELECT * FROM workflow_history
  WHERE created_at < NOW() - INTERVAL '30 days';

DELETE FROM workflow_history
  WHERE created_at < NOW() - INTERVAL '30 days';

-- Compress and backup
COPY (SELECT * FROM workflow_history_archive)
  TO PROGRAM 'gzip > /backup/workflow_history-archive-$(date +%Y%m%d).sql.gz';
EOF
```

---

### PostgreSQL — supabase-db Instance

| Tabela | Hot (days) | Warm (days) | Archive | Notes |
|--------|---|---|---|---|
| `sessions` | 7 | 30 | discard | Active sessions, auto-cleanup |
| `messages` | 30 | 90 | gzip | Chat history, user requested |
| `users` | Forever | N/A | N/A | User registry, never delete |
| `auth.users` | Forever | N/A | N/A | Managed by Supabase, never touch |

**Session Cleanup (Auto):**
```sql
-- Trigger: auto-delete inactive sessions after 7d
CREATE OR REPLACE FUNCTION cleanup_old_sessions() RETURNS void AS $$
BEGIN
  DELETE FROM sessions WHERE last_activity < NOW() - INTERVAL '7 days';
END;
$$ LANGUAGE plpgsql;

-- Schedule via pg_cron (if available)
SELECT cron.schedule('cleanup-sessions', '0 2 * * *', 'SELECT cleanup_old_sessions()');
```

---

### PostgreSQL — litellm-db Instance

| Tabela | Hot (days) | Warm (days) | Archive | Notes |
|--------|---|---|---|---|
| `api_keys` | Forever | N/A | N/A | Credentials, never delete |
| `usage_logs` | 30 | 60 | gzip | Token billing history |
| `model_configs` | Forever | N/A | N/A | Router config, never delete |
| `error_logs` | 7 | 30 | discard | Troubleshooting logs only |

**Cleanup:**
```bash
# Daily: Delete error logs > 7d
docker exec postgres-litellm psql -U postgres -d litellm \
  "DELETE FROM error_logs WHERE created_at < NOW() - INTERVAL '7 days'"

# Monthly: Archive usage_logs > 30d
docker exec postgres-litellm psql -U postgres -d litellm \
  "SELECT * FROM usage_logs WHERE created_at BETWEEN NOW() - INTERVAL '60 days' AND NOW() - INTERVAL '30 days'" \
  | gzip > ~/.aurelia/backups/usage_logs-$(date +%Y%m%d).sql.gz
```

---

### PostgreSQL — dev Instance

| Tabela | Hot (days) | Warm (days) | Archive | Notes |
|--------|---|---|---|---|
| (all) | N/A | N/A | N/A | Experimental, cleanup at will |

**No lifecycle policy (sandbox)**

---

### Qdrant (`repository_memory` Collection)

| Dado | Retenção | Regeneration | Notes |
|---|---|---|---|
| `repository_memory` vectors | Forever (Derived) | From Postgres monthly | Bge-m3 embeddings |
| Stale vectors (>90d old metadata) | Purge monthly | Regenerate on access | Keep metadata fresh |

**Regeneration Script:**
```bash
#!/bin/bash
# Monthly: Re-embed stale entries from Postgres
# (vectors older than 90d are regenerated from latest metadata)

# Query Postgres for entries with stale embeddings
STALE_IDS=$(psql $POSTGRES_DB -c \
  "SELECT id FROM ai_context.memory_entries
   WHERE synced_to_qdrant AND modified_at > synced_at + INTERVAL '90 days'")

# Re-embed via bge-m3 (local)
for id in $STALE_IDS; do
  # Fetch markdown content from $MEMORY_DIR
  # Run through bge-m3 encoder
  # Update Qdrant vector
done
```

---

## Backup & Recovery

### Hot Backup (Daily)
```bash
# SQLite: ZFS snapshot (automatic, hourly)
zfs snapshot storage/aurelia@daily-$(date +%Y%m%d)

# PostgreSQL: Docker volume snapshot (automatic via backup-sql.sh cron)
docker exec postgres-n8n pg_dump -U root n8n | gzip > \
  ~/.aurelia/backups/n8n-$(date +%Y%m%d).sql.gz
```

### Warm Archive (Monthly)
```bash
# Export 30-90d old data to gzip
sqlite3 ~/.aurelia/data/aurelia.db ".mode csv" \
  "SELECT * FROM voice_events WHERE created_at BETWEEN datetime('now', '-90 days') AND datetime('now', '-30 days')" \
  | gzip > ~/.aurelia/backups/voice_events-warm-$(date +%Y%m%d).csv.gz

# Also backup schema
sqlite3 ~/.aurelia/data/aurelia.db ".schema" > /tmp/schema.sql
gzip /tmp/schema.sql
mv /tmp/schema.sql.gz ~/.aurelia/backups/schema-$(date +%Y%m%d).sql.gz
```

### Cold Archive (Quarterly)
```bash
# All archived data (>90d old) to offsite or deep storage
cd ~/.aurelia/backups/
tar czf archive-$(date +%Y%m%d).tar.gz *-warm-*.gz *.schema.sql.gz
# Store in: /srv/backups/archive/ (ZFS long-term retention)
```

---

## Compliance & Legal Hold

### Data Subject Request (DSR)
**Scenario:** User requests export of all personal data
```bash
# Query all tables with user_id for data subject
psql supabase-db -c "
  SELECT * FROM users WHERE id = $user_id;
  SELECT * FROM sessions WHERE user_id = $user_id;
  SELECT * FROM messages WHERE user_id = $user_id;
"

# Export to encrypted ZIP for user
zip -e user_export_$user_id.zip *.csv
```

### Incident Forensics (Legal Hold)
**Scenario:** Security incident, preserve all logs
```bash
# Don't rotate/delete logs for 90 days
# Lock archive directory
chmod 555 ~/.aurelia/backups/incident-*

# Flag in audit:
# Event: LEGAL_HOLD
# Reason: Security incident #123
# Until: 2026-06-20 (90d)
```

### GDPR Right to Deletion
**Scenario:** User requests deletion
```bash
# ONLY delete if no legal hold, no compliance requirement
# Verify: no pending audits, no active investigation

# Delete from hot
sqlite3 ~/.aurelia/data/aurelia.db "DELETE FROM users WHERE id = $user_id"
psql supabase-db -c "DELETE FROM users WHERE id = $user_id"

# Scrub warm (if exists)
# Scrub archive (if exists)
# Log deletion timestamp in audit_log table
```

---

## Validation Commands

```bash
# Check hot data volume (last 30 days)
sqlite3 ~/.aurelia/data/aurelia.db \
  "SELECT COUNT(*) FROM voice_events WHERE created_at > datetime('now', '-30 days')"

# Check warm data volume (30-90 days old)
sqlite3 ~/.aurelia/data/aurelia.db \
  "SELECT COUNT(*) FROM voice_events WHERE created_at BETWEEN datetime('now', '-90 days') AND datetime('now', '-30 days')"

# Check archive backup size
du -sh ~/.aurelia/backups/*.gz

# Verify Qdrant can regenerate from Postgres
curl -s http://localhost:6333/collections/repository_memory | jq '.result.points_count'

# List all backup files with timestamps
ls -lhS ~/.aurelia/backups/ | grep -E "(hot|warm|archive)"
```

---

## Exceptions & Waivers

### Exception: system_health (7d instead of 30d)
- **Reason:** Metrics are high-volume, low-value after 7d
- **Approved:** codex
- **Expires:** N/A (permanent exception)
- **Validator:** `SELECT COUNT(*) FROM system_health WHERE created_at < datetime('now', '-7 days')` should return 0

### Exception: None currently
(Add exceptions via ADR-20260319 updates as needed)

---

## Links

- [Data Governance Domain Ownership](./data-governance-domain-ownership.md)
- [Store Selection Matrix](./data-governance-store-selection.md)
- [Schema Registry — SQLite](./schema-registry-sqlite.md)
- [Schema Registry — PostgreSQL](./schema-registry-postgres.md)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
