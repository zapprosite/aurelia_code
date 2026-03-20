---
title: Operational Governance — Health Checks
description: Watchdog crons e health check requirements para container, disk, VRAM, services
owner: codex
updated: 2026-03-20
---

# Operational Governance — Health Checks

**Purpose:** Definir automated health monitoring + alerting para infraestrutura
**Authority:** ADR-20260319-Polish-Governance-All / OPERATIONAL_GOVERNANCE
**Updated:** 2026-03-20

---

## Cron Schedule

### Every 5 minutes (Watchdog)
```
*/5 * * * * /home/will/aurelia/scripts/health-check.sh --mode quick 2>&1 | logger
```

**Checks:**
- All critical containers running (aurelia.service 8 MCP servers)
- Disk usage < 80% (alert if >80%)
- VRAM usage < 90% on RTX 4090 (alert if >90%, logs if >85%)
- Gateway route healthy (checks Cloudflare Tunnel + Tailscale)

**Timeout:** 30 seconds
**Failure Action:** Log + retry in 1 min (3 attempts max)

---

### Every 15 minutes (Service Health)
```
*/15 * * * * /home/will/aurelia/scripts/health-check.sh --mode services 2>&1 | logger
```

**Checks:**
- n8n service responding (HTTP 200 on /api/v1/workflows)
- PostgreSQL all instances accessible (connect + SELECT 1)
- Qdrant collection accessible (GET /collections/repository_memory)
- Supabase auth available (if enabled)

**Timeout:** 60 seconds
**Failure Action:** Alert + attempt restart (if transient)

---

### Every hour (Deep Health)
```
0 * * * * /home/will/aurelia/scripts/health-check.sh --mode deep 2>&1 | logger
```

**Checks:**
- Database schema consistency (compare reality vs registry)
- Backup age (<24h for daily, <7d for weekly)
- Memory sync staleness (memory-sync-fiscal.sh last run <1h ago)
- Docker volume disk usage (data, backups)
- ZFS pool health (DEGRADED? OFFLINE? MISSING?)

**Timeout:** 120 seconds
**Failure Action:** Alert + collect diagnostics

---

### Daily at 6am (Smoke Tests)
```
0 6 * * * /home/will/aurelia/scripts/health-check.sh --mode smoke 2>&1 | logger
```

**Tests:**
- End-to-end: Write test record → SQLite → Postgres → Qdrant search
- Backup recovery simulation (restore latest snapshot to /tmp)
- API latency benchmarks (p99 < 500ms for critical paths)
- Telegram bot functionality (send test message, verify response)

**Timeout:** 300 seconds (5 min)
**Failure Action:** Full diagnostic dump to `/tmp/smoke-test-$(date +%Y%m%d-%H%M%S).txt`

---

## Alert Rules

### CRITICAL Alerts (Immediate Page)

| Alert | Condition | Action | Escalation |
|-------|-----------|--------|---|
| **Container Down** | aurelia.service not running | Telegram (Codex) | Restart service, manual check |
| **Disk Full** | / or /srv usage >95% | Telegram (Codex) | Manual cleanup required |
| **VRAM Critical** | RTX 4090 >95% | Telegram (Codex) | Kill memory-intensive process |
| **Qdrant Unreachable** | No response for 5min | Telegram (Codex) | Restart container, validate data |
| **Backup Missing** | No backup <24h old | Telegram (Codex+Humano) | Verify backup script running |

**Response SLA:** <5 minutes

---

### HIGH Alerts (Alert but not immediate)

| Alert | Condition | Action | Escalation |
|-------|-----------|--------|---|
| **Disk Warning** | / usage 80-95% | Telegram (Codex) | Review cleanup opportunities |
| **VRAM Warning** | RTX 4090 85-95% | Telegram (Codex) | Monitor, may cause slowdowns |
| **Service Degraded** | Response time > 1s | Telegram (Codex) | Check query load, indexes |
| **Memory Sync Stale** | No sync >1 hour | Telegram (Codex) | Check fiscal cron status |
| **Schema Mismatch** | Detected at hourly check | Telegram (Codex) | Review schema registry |

**Response SLA:** <30 minutes

---

### MEDIUM Alerts (Log for review)

| Alert | Condition | Action |
|-------|-----------|--------|
| **Slow API Response** | p99 latency > 500ms | Log to metrics, graph trend |
| **Container Restart** | Service restarted in last hour | Log event, check for crashes |
| **Backup Size Growing** | Daily backups >10GB | Log warning, review retention |

**Response SLA:** <4 hours

---

## Health Check Script Template

```bash
#!/usr/bin/env bash
# scripts/health-check.sh — Automated health monitoring

set -euo pipefail

MODE="${1:-quick}"
LOG_FILE="/home/will/.aurelia/logs/health-check.log"
ALERT_THRESHOLD=2  # Alert after N consecutive failures

function alert() {
    local level="$1"
    local msg="$2"
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] [$level] $msg" | tee -a "$LOG_FILE"

    if [ "$level" = "CRITICAL" ] || [ "$level" = "HIGH" ]; then
        # Send Telegram alert
        curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
            -d "chat_id=${TELEGRAM_CHAT_ID}" \
            -d "text=⚠️ Health Check [$level]: $msg"
    fi
}

function check_container() {
    if ! docker ps --filter "name=aurelia" --format "{{.Status}}" | grep -q "Up"; then
        alert "CRITICAL" "aurelia container not running"
        return 1
    fi
    echo "Container: OK"
    return 0
}

function check_disk() {
    local usage=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
    if [ "$usage" -gt 95 ]; then
        alert "CRITICAL" "Disk usage critical: ${usage}%"
        return 1
    elif [ "$usage" -gt 80 ]; then
        alert "HIGH" "Disk usage warning: ${usage}%"
        return 0
    fi
    echo "Disk: OK (${usage}%)"
    return 0
}

function check_vram() {
    # Parse nvidia-smi output for RTX 4090 usage
    local vram_used=$(nvidia-smi --query-gpu=memory.used --format=csv,noheader,nounits | awk '{print int($1)}')
    local vram_total=24576  # RTX 4090 = 24GB
    local percent=$((vram_used * 100 / vram_total))

    if [ "$percent" -gt 95 ]; then
        alert "CRITICAL" "VRAM critical: ${percent}% (${vram_used}MB / ${vram_total}MB)"
        return 1
    elif [ "$percent" -gt 85 ]; then
        alert "HIGH" "VRAM warning: ${percent}% (${vram_used}MB / ${vram_total}MB)"
        return 0
    fi
    echo "VRAM: OK (${percent}%)"
    return 0
}

function check_qdrant() {
    if ! curl -s http://localhost:6333/health >/dev/null 2>&1; then
        alert "CRITICAL" "Qdrant unreachable"
        return 1
    fi
    local count=$(curl -s http://localhost:6333/collections/repository_memory | jq '.result.points_count')
    echo "Qdrant: OK ($count vectors)"
    return 0
}

function check_postgres() {
    for db in n8n supabase litellm; do
        if ! docker exec postgres-$db pg_isready 2>/dev/null; then
            alert "HIGH" "PostgreSQL instance $db not ready"
            return 1
        fi
    done
    echo "PostgreSQL: OK (all instances)"
    return 0
}

function check_backup() {
    local latest=$(ls -t ~/.aurelia/backups/*.sql.gz 2>/dev/null | head -1)
    if [ -z "$latest" ]; then
        alert "CRITICAL" "No backups found"
        return 1
    fi

    local age_hours=$(($(date +%s) - $(date -d "$(stat -c %y "$latest")" +%s)) / 3600))
    if [ "$age_hours" -gt 24 ]; then
        alert "CRITICAL" "Backup stale: ${age_hours}h old"
        return 1
    fi
    echo "Backup: OK (${age_hours}h old)"
    return 0
}

function check_memory_sync() {
    local last_sync=$(stat -c %Y ~/.aurelia/metrics/memory-sync.prom 2>/dev/null || echo 0)
    local now=$(date +%s)
    local age_min=$(( (now - last_sync) / 60 ))

    if [ "$age_min" -gt 60 ]; then
        alert "HIGH" "Memory sync stale: ${age_min}m old"
        return 0
    fi
    echo "Memory Sync: OK (${age_min}m old)"
    return 0
}

case "$MODE" in
    quick)
        check_container && check_disk && check_vram && check_qdrant
        ;;
    services)
        check_postgres && check_qdrant && check_backup
        ;;
    deep)
        check_backup && check_memory_sync
        ;;
    smoke)
        echo "Running smoke tests..."
        # E2E test: write → search
        # Recovery simulation
        # Latency benchmarks
        ;;
    *)
        echo "Unknown mode: $MODE" >&2
        exit 1
        ;;
esac
```

---

## Metric Collection

Health checks write Prometheus metrics:

```
# Cron every 5min writes to ~/.aurelia/metrics/health.prom
aurelia_container_status{name="aurelia"} 1
aurelia_disk_usage_percent 45
aurelia_vram_usage_percent 62
aurelia_vram_mb 15265
aurelia_qdrant_available 1
aurelia_postgres_available{instance="n8n"} 1
aurelia_backup_age_hours 6
aurelia_memory_sync_age_minutes 3
```

**Scrape interval:** Prometheus pulls from `~/.aurelia/metrics/*.prom` every 1 minute

---

## Failure Recovery

### Container Crash Recovery
```bash
# health-check.sh detects container down
# Automatic action:
docker restart aurelia

# If restart fails 3 times:
# - Alert "CRITICAL: Container restart failed"
# - Wait for manual intervention
```

### Disk Full Recovery
```bash
# Check if old backups can be archived:
find ~/.aurelia/backups/ -name "*.sql.gz" -mtime +30 \
    -exec mv {} /srv/backups/archive/ \;

# If still full, alert for manual cleanup
# (Do not auto-delete without approval)
```

### VRAM Exceeded Recovery
```bash
# If >95%: Kill lowest-priority process
# (e.g., minor inference, not n8n workflows)

pkill -f "qwen3.5:9b" --oldest-first

# Re-alert after 5 minutes if still high
```

---

## Cron Installation

```bash
#!/bin/bash
# Install health check crons

mkdir -p ~/.aurelia/logs ~/.aurelia/metrics

# Add to user crontab
(crontab -l 2>/dev/null; cat <<'EOF'
*/5 * * * * /home/will/aurelia/scripts/health-check.sh --mode quick >> /home/will/.aurelia/logs/cron.log 2>&1
*/15 * * * * /home/will/aurelia/scripts/health-check.sh --mode services >> /home/will/.aurelia/logs/cron.log 2>&1
0 * * * * /home/will/aurelia/scripts/health-check.sh --mode deep >> /home/will/.aurelia/logs/cron.log 2>&1
0 6 * * * /home/will/aurelia/scripts/health-check.sh --mode smoke >> /home/will/.aurelia/logs/cron.log 2>&1
EOF
) | crontab -

# Verify installation
crontab -l | grep health-check
```

---

## Validation

```bash
# Test health check script manually
bash scripts/health-check.sh --mode quick
bash scripts/health-check.sh --mode services
bash scripts/health-check.sh --mode deep

# Check cron jobs are installed
crontab -l | grep "health-check"

# Monitor active checks
tail -f ~/.aurelia/logs/health-check.log

# View current metrics
cat ~/.aurelia/metrics/health.prom
```

---

## Links

- [Data Lifecycle Policy](./data-governance-lifecycle.md)
- [Incident Response Playbook](./operational-governance-incident-response.md)
- [Backup Verification](./operational-governance-backup-verification.md)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
