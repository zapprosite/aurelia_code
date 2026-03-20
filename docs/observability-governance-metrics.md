---
title: Observability Governance — Metrics Contract
description: Prometheus metrics obrigatórias, scrape intervals, e dashboard requirements
owner: codex
updated: 2026-03-20
---

# Observability Governance — Metrics Contract

**Purpose:** Definir quais métricas cada serviço DEVE expor para observabilidade
**Authority:** ADR-20260319-Polish-Governance-All / OBSERVABILITY_GOVERNANCE
**Updated:** 2026-03-20

---

## Mandatory Metrics by Service

### Core Infrastructure

**Every service MUST expose:**

```
service_up                                    # 1 if healthy, 0 if down
                                             # Type: gauge
                                             # Labels: service={name}

service_request_total{method,path,status}    # Total requests
                                             # Type: counter
                                             # Labels: method, path, status

service_request_duration_seconds{quantile}   # Request latency
                                             # Type: histogram with buckets
                                             # Buckets: [.001, .01, .1, .5, 1, 5]
                                             # Labels: path, method

service_errors_total{error_type}             # Total errors by type
                                             # Type: counter
                                             # Labels: error_type (timeout, 500, 403, etc)
```

**Scrape interval:** 30 seconds
**Timeout:** 10 seconds
**Port:** Service-specific (documented in NETWORK_MAP.md)

---

### Container Metrics (Required for all Docker services)

```
container_cpu_usage_seconds_total            # CPU time
container_memory_usage_bytes                 # Memory usage
container_network_receive_bytes_total        # Network inbound
container_network_transmit_bytes_total       # Network outbound
container_restart_count                      # Number of restarts
```

**Source:** Prometheus cAdvisor integration
**Scrape interval:** 30 seconds

---

### System Metrics (Infrastructure)

From health-check.sh crons → `~/.aurelia/metrics/health.prom`:

```
aurelia_disk_usage_percent                   # / filesystem usage
                                             # Type: gauge

aurelia_vram_usage_percent                   # RTX 4090 VRAM %
                                             # Type: gauge

aurelia_vram_mb                              # RTX 4090 VRAM used (MB)
                                             # Type: gauge

aurelia_container_status{name}               # 1 if up, 0 if down
                                             # Type: gauge

aurelia_postgres_available{instance}         # 1 if DB accessible
                                             # Type: gauge

aurelia_qdrant_available                     # 1 if collection accessible
                                             # Type: gauge

aurelia_backup_age_hours                     # Latest backup age
                                             # Type: gauge

aurelia_memory_sync_age_minutes               # Memory sync last run
                                             # Type: gauge
```

**Source:** `health-check.sh` every 5 minutes
**File:** `~/.aurelia/metrics/health.prom`
**Scrape interval:** 1 minute from file

---

### Database Metrics (PostgreSQL)

From postgres_exporter:

```
pg_stat_statements_calls                     # Query call count
pg_stat_statements_total_time_seconds        # Total query time
pg_stat_activity_count{state}                # Active connections by state
pg_database_size_bytes{datname}              # Database size
pg_table_size_bytes{relname}                 # Table size
pg_cache_hit_ratio{datname}                  # Cache hit ratio (%)
pg_stat_user_rows_inserted_total{relname}    # Rows inserted
pg_stat_user_rows_deleted_total{relname}     # Rows deleted
```

**Deployment:** postgres_exporter sidecar per instance
**Scrape interval:** 30 seconds
**Port:** 9187 (per instance, offset by +N)

---

### Memory Sync Metrics (Qdrant + Fiscal)

From memory-sync-fiscal.sh → `~/.aurelia/metrics/memory-sync.prom`:

```
aurelia_memory_files_total                   # Files scanned
                                             # Type: gauge

aurelia_memory_sync_duration_seconds{mode}   # Duration by mode (fast/postgres/validate/compact)
                                             # Type: gauge
                                             # Labels: mode

aurelia_qdrant_points_total                  # Vectors in repository_memory collection
                                             # Type: gauge

aurelia_qdrant_sync_latency_seconds          # Qdrant write latency
                                             # Type: histogram

aurelia_postgres_memory_entries_synced       # Rows synced to Postgres
                                             # Type: counter

aurelia_memory_cache_hit_ratio                # bge-m3 embedding cache hits
                                             # Type: gauge
```

**Source:** memory-sync-fiscal.sh every 5/15 minutes
**File:** `~/.aurelia/metrics/memory-sync.prom`
**Scrape interval:** 1 minute from file

---

### Application-Specific Metrics

#### n8n (Workflow Engine)
```
n8n_workflow_executions_total{status}        # Workflow runs (success/error/aborted)
n8n_workflow_execution_duration_seconds      # Execution time
n8n_webhook_requests_total                   # Webhook calls
n8n_error_logs_total{error_type}             # Errors by type
```

#### litellm (LLM Gateway)
```
litellm_requests_total{model,status}         # LLM requests
litellm_tokens_total{model,type}             # Tokens used (prompt/completion)
litellm_request_duration_seconds{model}      # Latency by model
litellm_fallback_triggered_total{reason}     # Fallback events
litellm_error_total{model,error_type}        # Errors by model
```

#### Supabase/PostgreSQL
```
supabase_user_sessions_total{status}         # Session metrics
supabase_auth_requests_total{status}         # Auth endpoint calls
supabase_storage_usage_bytes                 # Storage consumption
```

---

## Scrape Configuration

### Prometheus Config (`/etc/prometheus/prometheus.yml` snippet)

```yaml
global:
  scrape_interval: 30s
  evaluation_interval: 30s
  external_labels:
    cluster: aurelia

scrape_configs:
  # Core service metrics
  - job_name: 'aurelia'
    metrics_path: '/metrics'
    static_configs:
      - targets: ['localhost:5000']
    scrape_interval: 30s
    scrape_timeout: 10s

  # Health check metrics (file-based)
  - job_name: 'aurelia-health'
    static_configs:
      - targets: ['localhost']
        labels:
          __metrics_path__: '/var/aurelia/metrics/health.prom'
    scrape_interval: 1m
    scrape_timeout: 10s

  # Memory sync metrics (file-based)
  - job_name: 'aurelia-memory-sync'
    static_configs:
      - targets: ['localhost']
        labels:
          __metrics_path__: '/var/aurelia/metrics/memory-sync.prom'
    scrape_interval: 1m
    scrape_timeout: 10s

  # PostgreSQL exporters (per instance)
  - job_name: 'postgres-n8n'
    static_configs:
      - targets: ['localhost:9187']
    scrape_interval: 30s

  - job_name: 'postgres-supabase'
    static_configs:
      - targets: ['localhost:9188']
    scrape_interval: 30s

  - job_name: 'postgres-litellm'
    static_configs:
      - targets: ['localhost:9189']
    scrape_interval: 30s

  # Container metrics
  - job_name: 'cadvisor'
    static_configs:
      - targets: ['localhost:8080']
    scrape_interval: 30s

  # Optional: Node exporter (system metrics)
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
    scrape_interval: 30s
```

---

## Metric Publishing

### Health Check Metrics (File-based)

**Every 5 minutes** via `health-check.sh`:

```bash
#!/bin/bash
# Metrics file: ~/.aurelia/metrics/health.prom

cat > "$METRICS_DIR/health.prom.tmp" <<EOF
# HELP aurelia_disk_usage_percent Root filesystem usage percentage
# TYPE aurelia_disk_usage_percent gauge
aurelia_disk_usage_percent $(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')

# HELP aurelia_vram_usage_percent RTX 4090 VRAM usage percentage
# TYPE aurelia_vram_usage_percent gauge
aurelia_vram_usage_percent $(nvidia-smi --query-gpu=memory.used --format=csv,noheader,nounits)

# ... more metrics
EOF

mv "$METRICS_DIR/health.prom.tmp" "$METRICS_DIR/health.prom"
```

---

## Alert Rules

### Prometheus AlertManager Config

```yaml
groups:
  - name: aurelia-alerts
    interval: 30s
    rules:
      # CRITICAL alerts (page immediately)
      - alert: ContainerDown
        expr: aurelia_container_status{name="aurelia"} == 0
        for: 1m
        severity: critical
        annotations:
          summary: "Aurelia container is down"

      - alert: DiskSpaceCritical
        expr: aurelia_disk_usage_percent > 95
        for: 2m
        severity: critical
        annotations:
          summary: "Disk usage critical: {{ $value }}%"

      - alert: VRAMCritical
        expr: aurelia_vram_usage_percent > 95
        for: 2m
        severity: critical
        annotations:
          summary: "VRAM critical: {{ $value }}%"

      - alert: BackupMissing
        expr: aurelia_backup_age_hours > 24
        for: 1h
        severity: critical
        annotations:
          summary: "No backup in last 24 hours"

      # HIGH alerts (notify, but not immediate page)
      - alert: DiskSpaceWarning
        expr: aurelia_disk_usage_percent > 80
        for: 5m
        severity: warning
        annotations:
          summary: "Disk usage warning: {{ $value }}%"

      - alert: VRAMWarning
        expr: aurelia_vram_usage_percent > 85
        for: 5m
        severity: warning
        annotations:
          summary: "VRAM warning: {{ $value }}%"

      - alert: MemorySyncStale
        expr: aurelia_memory_sync_age_minutes > 60
        for: 5m
        severity: warning
        annotations:
          summary: "Memory sync hasn't run in {{ $value }} minutes"

      # MEDIUM alerts (log)
      - alert: SlowAPIResponse
        expr: service_request_duration_seconds{quantile="0.99"} > 1
        for: 5m
        severity: info
        annotations:
          summary: "API p99 latency high: {{ $value }}s"

      - alert: DatabaseQuerySlow
        expr: pg_stat_statements_total_time_seconds{quantile="0.99"} > 5
        for: 5m
        severity: info
        annotations:
          summary: "Slow database query detected"
```

---

## Dashboard Requirements

### Required Dashboards (Grafana)

#### 1. Overview Dashboard (`aurelia-overview`)
**Tiles required:**
- Container status (UP/DOWN)
- Disk usage (%)
- VRAM usage (% + MB)
- Backup age (hours)
- Last memory sync (minutes ago)
- Active requests (rate)

**Refresh rate:** 30 seconds
**Audience:** Operations team (quick health check)

#### 2. System Dashboard (`aurelia-system`)
**Panels required:**
- CPU usage by container
- Memory usage by container
- Network I/O (bytes/sec)
- Disk I/O (IOPS)
- ZFS pool health (if applicable)
- Container restart count

**Refresh rate:** 1 minute

#### 3. Database Dashboard (`aurelia-databases`)
**Per-instance panels:**
- Connection count
- Transaction rate
- Cache hit ratio (%)
- Query latency (p50/p99)
- Table sizes
- Replication lag (if applicable)

**Refresh rate:** 30 seconds
**Instances:** n8n, supabase, litellm, dev

#### 4. Memory Sync Dashboard (`aurelia-memory-sync`)
**Panels:**
- Files scanned (trend)
- Vectors in Qdrant (trend)
- Rows synced to Postgres (counter)
- Sync duration by mode (boxplot)
- Embedding latency (p50/p99)
- Qdrant query latency (p99)

**Refresh rate:** 5 minutes

#### 5. LLM Metrics Dashboard (`aurelia-llm`)
**Panels:**
- LLM request rate by model (stacked area)
- Token usage by model (stacked bar)
- Model latency distribution (heatmap)
- Fallback triggers (bar chart)
- Error rate by model (gauge)

**Refresh rate:** 1 minute

---

## Validation

### Metrics Availability Check

```bash
# Every service should expose these metrics
curl http://localhost:5000/metrics | grep -c "aurelia_"  # Should be >10

# Health metrics should be published every 5 minutes
ls -l ~/.aurelia/metrics/health.prom
# Age should be < 6 minutes

# Memory sync metrics should be published every 5-15 minutes
ls -l ~/.aurelia/metrics/memory-sync.prom
# Age should be < 20 minutes

# Prometheus should scrape all targets
curl -s http://localhost:9090/api/v1/targets | jq '.data.activeTargets | length'  # Should be >8

# Alert rules should be loaded
curl -s http://localhost:9090/api/v1/rules | jq '.data.groups[].rules | length'  # Should be >5
```

---

## Links

- [Health Checks](./operational-governance-health-checks.md)
- [Alert Rules & Incident Response](./operational-governance-incident-response.md)
- [Qdrant Collection Contract](./qdrant-collection-contract.md)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
