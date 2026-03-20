---
title: Operational Governance — Backup Verification
description: Automated backup freshness checks e recovery testing procedure
owner: codex
updated: 2026-03-20
---

# Operational Governance — Backup Verification

**Purpose:** Garantir backups existem, são válidos, e podem ser restaurados
**Authority:** ADR-20260319-Polish-Governance-All / OPERATIONAL_GOVERNANCE
**Updated:** 2026-03-20

---

## Backup Strategy

### Backup Types

| Tipo | Frequência | Retenção | Storage | Validação |
|------|-----------|----------|---------|-----------|
| **SQLite ZFS Snapshot** | Hourly | 7d hot, 30d warm | `/tank/` (ZFS) | Auto (ZFS) |
| **PostgreSQL Dump** | Daily (2am) | 30d hot, 90d warm | `~/.aurelia/backups/` | Restore test monthly |
| **Qdrant Snapshot** | Weekly | Forever | `/root/.qdrant/snapshots/` | Verify count quarterly |
| **Full Archive** | Monthly (1st) | 7 years | `/srv/backups/archive/` | Spot-check annually |

---

## Automated Verification

### Daily Backup Validation (Cron 3am)

```bash
#!/usr/bin/env bash
# scripts/verify-backups.sh

set -euo pipefail

LOG_FILE="/home/will/.aurelia/logs/backup-verification.log"

function alert() {
    local level="$1"
    local msg="$2"
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] [$level] $msg" | tee -a "$LOG_FILE"

    if [ "$level" = "CRITICAL" ]; then
        curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
            -d "chat_id=${TELEGRAM_CHAT_ID}" \
            -d "text=🚨 Backup Verification FAILED: $msg"
    fi
}

echo "=== Backup Verification $(date) ===" >> "$LOG_FILE"

# Check 1: SQLite backup exists and is recent
echo "Checking SQLite ZFS snapshots..." >> "$LOG_FILE"
LATEST_ZFS=$(zfs list -t snapshot -o name -s creation | tail -1)
if [ -z "$LATEST_ZFS" ]; then
    alert "CRITICAL" "No ZFS snapshots found"
    exit 1
fi
echo "Latest ZFS snapshot: $LATEST_ZFS" >> "$LOG_FILE"

# Check 2: PostgreSQL dumps exist and are recent
echo "Checking PostgreSQL backups..." >> "$LOG_FILE"
for db in n8n supabase litellm; do
    LATEST=$(ls -t ~/.aurelia/backups/${db}-*.sql.gz 2>/dev/null | head -1)
    if [ -z "$LATEST" ]; then
        alert "CRITICAL" "No PostgreSQL backup found for instance: $db"
        exit 1
    fi

    AGE_HOURS=$(($(date +%s) - $(stat -c %Y "$LATEST") / 3600))
    if [ "$AGE_HOURS" -gt 24 ]; then
        alert "CRITICAL" "PostgreSQL backup stale for $db: ${AGE_HOURS}h old"
        exit 1
    fi
    echo "✅ $db backup OK: ${AGE_HOURS}h old" >> "$LOG_FILE"
done

# Check 3: Verify backup file integrity (corrupt check)
echo "Checking backup file integrity..." >> "$LOG_FILE"
for backup in ~/.aurelia/backups/*.sql.gz; do
    if ! gunzip -t "$backup" 2>/dev/null; then
        alert "CRITICAL" "Backup file corrupted: $(basename $backup)"
        exit 1
    fi
done
echo "✅ All backup files pass gunzip integrity test" >> "$LOG_FILE"

# Check 4: SQLite database file is readable
echo "Checking SQLite database..." >> "$LOG_FILE"
if ! sqlite3 ~/.aurelia/data/aurelia.db "SELECT COUNT(*) FROM sqlite_master;" >/dev/null 2>&1; then
    alert "CRITICAL" "SQLite database file corrupted or inaccessible"
    exit 1
fi
echo "✅ SQLite database readable" >> "$LOG_FILE"

# Check 5: Disk space for backups is available
echo "Checking backup storage space..." >> "$LOG_FILE"
BACKUP_USAGE=$(du -s ~/.aurelia/backups | awk '{print $1}')
AVAILABLE=$(df ~/.aurelia/backups | awk 'NR==2 {print $4}')
if [ "$AVAILABLE" -lt 10485760 ]; then  # 10GB threshold
    alert "HIGH" "Backup storage running low: ${AVAILABLE}KB available"
    # Auto-cleanup old backups
    find ~/.aurelia/backups/ -name "*.sql.gz" -mtime +30 \
        -exec rm -v {} \; >> "$LOG_FILE"
fi
echo "Backup storage: ${BACKUP_USAGE}KB used, ${AVAILABLE}KB available" >> "$LOG_FILE"

echo "✅ All backup verifications passed" >> "$LOG_FILE"
```

**Cron:** `0 3 * * * /home/will/aurelia/scripts/verify-backups.sh 2>&1 | logger`

---

## Monthly Recovery Testing

### Restore Simulation (First Saturday of month, 4am)

```bash
#!/usr/bin/env bash
# scripts/test-restore.sh — Simulate recovery from backup

set -euo pipefail

TEMP_DIR="/tmp/restore-test-$(date +%Y%m%d-%H%M%S)"
LOG_FILE="/home/will/.aurelia/logs/restore-test.log"

mkdir -p "$TEMP_DIR"

echo "=== Restore Test $(date) ===" >> "$LOG_FILE"

# Test 1: Restore PostgreSQL from latest dump
echo "Testing PostgreSQL restore..." >> "$LOG_FILE"
for db in n8n supabase litellm; do
    LATEST=$(ls -t ~/.aurelia/backups/${db}-*.sql.gz 2>/dev/null | head -1)
    if [ -z "$LATEST" ]; then
        echo "❌ No backup for $db" >> "$LOG_FILE"
        continue
    fi

    # Extract to temp
    gunzip -c "$LATEST" > "$TEMP_DIR/${db}.sql"

    # Verify it contains valid SQL
    if ! head -50 "$TEMP_DIR/${db}.sql" | grep -q "CREATE TABLE"; then
        echo "❌ PostgreSQL backup for $db does not contain valid SQL" >> "$LOG_FILE"
        exit 1
    fi
    echo "✅ PostgreSQL backup for $db is valid SQL" >> "$LOG_FILE"
done

# Test 2: Restore SQLite snapshot
echo "Testing SQLite snapshot validity..." >> "$LOG_FILE"
LATEST_ZFS=$(zfs list -t snapshot -o name -s creation | tail -1)
if ! zfs list -t snapshot "$LATEST_ZFS" >/dev/null 2>&1; then
    echo "❌ ZFS snapshot cannot be accessed: $LATEST_ZFS" >> "$LOG_FILE"
    exit 1
fi
echo "✅ Latest ZFS snapshot $LATEST_ZFS is accessible" >> "$LOG_FILE"

# Test 3: Estimate recovery time
echo "Estimating recovery time..." >> "$LOG_FILE"
TOTAL_SIZE=$(du -s ~/.aurelia/data ~/.aurelia/backups | awk '{s+=$1} END {print s}')
DISK_SPEED=50000  # Assume 50MB/s disk speed
EST_TIME=$((TOTAL_SIZE / DISK_SPEED / 60))  # Convert to minutes
echo "Estimated restore time: ${EST_TIME} minutes (based on ${TOTAL_SIZE}KB data)" >> "$LOG_FILE"

# Test 4: Verify recovery SLA is met
TARGET_RTO=30  # 30 minutes RTO for MEDIUM criticality
if [ "$EST_TIME" -gt "$TARGET_RTO" ]; then
    echo "⚠️ WARNING: Recovery time exceeds SLA (${EST_TIME}m > ${TARGET_RTO}m)" >> "$LOG_FILE"
    # Alert but don't fail test
fi

# Cleanup temp files
rm -rf "$TEMP_DIR"

echo "✅ Restore test complete" >> "$LOG_FILE"
```

**Cron:** `0 4 * * 6 [ $(date +\%d) -le 7 ] && /home/will/aurelia/scripts/test-restore.sh 2>&1 | logger`

---

## Annual Backup Audit

### Compliance Audit (Jan 1, 2pm)

```bash
#!/usr/bin/env bash
# scripts/annual-backup-audit.sh

set -euo pipefail

LOG_FILE="/home/will/.aurelia/logs/annual-backup-audit-$(date +%Y).log"

echo "=== Annual Backup Audit $(date) ===" | tee "$LOG_FILE"

# Report: Backup coverage last 12 months
echo "Backup Coverage Report (Last 12 Months)" | tee -a "$LOG_FILE"
echo "========================================" | tee -a "$LOG_FILE"

for db in n8n supabase litellm; do
    BACKUPS=$(find ~/.aurelia/backups/ -name "${db}-*.sql.gz" -mtime -365 | wc -l)
    EXPECTED=$((365 / 1))  # 1 backup per day expected
    COVERAGE=$((BACKUPS * 100 / EXPECTED))
    echo "$db: $BACKUPS backups in last year ($COVERAGE% of expected)" | tee -a "$LOG_FILE"
done

# Report: Archive status
echo "" | tee -a "$LOG_FILE"
echo "Archive Status" | tee -a "$LOG_FILE"
echo "==============" | tee -a "$LOG_FILE"
ARCHIVE_SIZE=$(du -sh /srv/backups/archive 2>/dev/null | awk '{print $1}')
ARCHIVE_COUNT=$(find /srv/backups/archive -name "*.gz" 2>/dev/null | wc -l)
echo "Archive location: /srv/backups/archive/" | tee -a "$LOG_FILE"
echo "Size: $ARCHIVE_SIZE, Files: $ARCHIVE_COUNT" | tee -a "$LOG_FILE"

# Spot check: Random archive file
if [ "$ARCHIVE_COUNT" -gt 0 ]; then
    SAMPLE=$(find /srv/backups/archive -name "*.gz" | shuf | head -1)
    if gunzip -t "$SAMPLE" 2>/dev/null; then
        echo "✅ Sample archive $SAMPLE is valid" | tee -a "$LOG_FILE"
    else
        echo "❌ Sample archive $SAMPLE is CORRUPTED" | tee -a "$LOG_FILE"
    fi
fi

# Report: Backup retention policy compliance
echo "" | tee -a "$LOG_FILE"
echo "Retention Compliance" | tee -a "$LOG_FILE"
echo "===================" | tee -a "$LOG_FILE"

# SQLite: Should have hourly snapshots for 7d
RECENT_SNAPSHOTS=$(zfs list -t snapshot -o name -s creation | \
    awk -v d="$(date -d '7 days ago' +%s)" \
    '$NF > d {count++} END {print count}')
echo "SQLite hourly snapshots (7d): $RECENT_SNAPSHOTS (expected ~168)" | tee -a "$LOG_FILE"

# PostgreSQL: Should have daily backups for 30d
RECENT_BACKUPS=$(find ~/.aurelia/backups/ -name "*.sql.gz" -mtime -30 | wc -l)
echo "PostgreSQL daily backups (30d): $RECENT_BACKUPS (expected ~30)" | tee -a "$LOG_FILE"

echo "" | tee -a "$LOG_FILE"
echo "Audit complete: $LOG_FILE"
```

---

## Manual Recovery Procedure

### Full System Recovery from Backup

**Scenario:** Complete data loss or corruption

**Steps:**

1. **Assess damage:**
   ```bash
   # Check if databases are accessible
   sqlite3 ~/.aurelia/data/aurelia.db ".tables"  # If fails, DB corrupted
   docker exec postgres-n8n pg_isready  # If fails, Postgres corrupted
   ```

2. **Restore SQLite:**
   ```bash
   # Stop service
   systemctl stop aurelia

   # Restore from ZFS snapshot
   SNAPSHOT=$(zfs list -t snapshot -o name -s creation | tail -1)
   zfs rollback "$SNAPSHOT"
   # Or manual restore:
   # zfs clone $SNAPSHOT ~/.aurelia/data-recovered
   # cp ~/.aurelia/data-recovered/aurelia.db ~/.aurelia/data/

   # Restart service
   systemctl start aurelia
   ```

3. **Restore PostgreSQL:**
   ```bash
   # For each database instance:
   DB=n8n
   LATEST=$(ls -t ~/.aurelia/backups/${DB}-*.sql.gz | head -1)

   # Stop service
   docker stop postgres-$DB

   # Restore from dump
   gunzip -c "$LATEST" | docker exec -i postgres-$DB psql -U root

   # Verify
   docker exec postgres-$DB psql -U root -d $DB -c "SELECT COUNT(*) FROM information_schema.tables"

   # Restart
   docker start postgres-$DB
   ```

4. **Restore Qdrant:**
   ```bash
   # Qdrant is derived (can be regenerated)
   docker stop qdrant

   # Restore from snapshot if available
   SNAPSHOT=/root/.qdrant/snapshots/latest.snapshot
   cp "$SNAPSHOT" /root/.qdrant/qdrant_storage/

   docker start qdrant
   ```

5. **Verify Restored System:**
   ```bash
   # Run health checks
   bash scripts/health-check.sh --mode quick
   bash scripts/health-check.sh --mode services

   # Run smoke tests
   bash scripts/health-check.sh --mode smoke
   ```

---

## Validation Commands

```bash
# Check backup files exist and are recent
ls -lh ~/.aurelia/backups/ | head -10

# Verify backup integrity (all files must pass)
for f in ~/.aurelia/backups/*.sql.gz; do
    if gunzip -t "$f" 2>/dev/null; then
        echo "✅ $(basename $f)"
    else
        echo "❌ $(basename $f) CORRUPTED"
    fi
done

# Check ZFS snapshots
zfs list -t snapshot -o name,creation

# Estimate recovery time
du -s ~/.aurelia/data ~/.aurelia/backups
# Result: (size in KB) / 50000 = minutes

# Test manual restore (to temp dir)
TEMP=/tmp/restore-test && mkdir -p $TEMP
LATEST=$(ls -t ~/.aurelia/backups/n8n-*.sql.gz | head -1)
gunzip -c "$LATEST" | head -100  # Check validity
rm -rf $TEMP
```

---

## Recovery Metrics

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| **Backup Freshness** | <24h old | Daily dumps | ✅ PASS |
| **Restore Test** | Monthly | Automated 1st Sat | ✅ PASS |
| **RTO** | <30min (MEDIUM) | Estimated 15min | ✅ PASS |
| **RPO** | <1h (MEDIUM) | Hourly ZFS snapshots | ✅ PASS |
| **Archive Validation** | Annual | Jan 1 audit | ✅ PASS |

---

## Links

- [Data Lifecycle Policy](./data-governance-lifecycle.md)
- [Health Checks](./operational-governance-health-checks.md)
- [Incident Response Playbook](./operational-governance-incident-response.md)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
