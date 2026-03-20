---
title: Schema Registry — SQLite
description: Registro completo de tabelas, colunas, índices e políticas de retenção do banco SQLite aurelia.db
owner: codex
updated: 2026-03-20
---

# Schema Registry — SQLite (`aurelia.db`)

**Source of Truth:** ✅ Primário para gateway, voice events, cron, tasks, memory
**Location:** `~/.aurelia/data/aurelia.db`
**Backup:** ZFS snapshot automático + `~/.aurelia/backups/`
**Retenção:** 30d hot, 90d warm, >90d archive gzipped

---

## Tabelas Planejadas

### 1. `gateway_route_states`

**Propósito:** Registrar estado de rotas do gateway + fallback status
**Criticidade:** HIGH
**Owner:** codex (operação), humano (policy)

**Schema:**
```sql
CREATE TABLE IF NOT EXISTS gateway_route_states (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    route_name TEXT NOT NULL UNIQUE,
    target_url TEXT NOT NULL,
    status TEXT CHECK(status IN ('up', 'down', 'degraded')),
    last_check_timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    health_check_interval_seconds INTEGER DEFAULT 60,
    fallback_url TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_gateway_route_states_status ON gateway_route_states(status);
CREATE INDEX idx_gateway_route_states_last_check ON gateway_route_states(last_check_timestamp);
```

**Índices:**
- `status` (para queries de rota "down")
- `last_check_timestamp` (para cleanup de records stale)

**Política de Retenção:** Forever (operacional)

---

### 2. `voice_events`

**Propósito:** Log de eventos de síntese de voz (TTS), qualidade, latência
**Criticidade:** MEDIUM
**Owner:** codex (operação)

**Schema:**
```sql
CREATE TABLE IF NOT EXISTS voice_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type TEXT NOT NULL CHECK(event_type IN ('tts_request', 'tts_complete', 'error')),
    provider TEXT NOT NULL CHECK(provider IN ('gemini', 'minimax', 'fallback')),
    text_length INTEGER,
    duration_ms INTEGER,
    status_code INTEGER,
    error_message TEXT,
    voice_id TEXT,
    sample_rate INTEGER DEFAULT 24000,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_voice_events_provider ON voice_events(provider);
CREATE INDEX idx_voice_events_created_at ON voice_events(created_at);
CREATE INDEX idx_voice_events_event_type ON voice_events(event_type);
```

**Política de Retenção:** 30d hot, 90d warm, >90d archive

---

### 3. `cron_tasks`

**Propósito:** Histórico de execução de crons (memory-sync-fiscal, health-check, etc)
**Criticidade:** MEDIUM
**Owner:** codex

**Schema:**
```sql
CREATE TABLE IF NOT EXISTS cron_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_name TEXT NOT NULL,
    mode TEXT,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    status TEXT CHECK(status IN ('pending', 'running', 'success', 'failure')),
    duration_ms INTEGER,
    error_log TEXT,
    metrics_json TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cron_tasks_task_name ON cron_tasks(task_name);
CREATE INDEX idx_cron_tasks_status ON cron_tasks(status);
CREATE INDEX idx_cron_tasks_start_time ON cron_tasks(start_time);
```

**Política de Retenção:** 30d hot (operacional), 90d warm (auditoria)

---

### 4. `memory_entries`

**Propósito:** Local memory cache antes de sync para Qdrant
**Criticidade:** HIGH
**Owner:** codex

**Schema:**
```sql
CREATE TABLE IF NOT EXISTS memory_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL UNIQUE,
    type TEXT CHECK(type IN ('project_memory', 'adr', 'runbook', 'plan')),
    owner TEXT,
    tags TEXT,
    file_hash TEXT NOT NULL,
    file_size INTEGER,
    modified_at DATETIME NOT NULL,
    synced_to_qdrant BOOLEAN DEFAULT 0,
    synced_to_postgres BOOLEAN DEFAULT 0,
    last_sync_attempt DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_memory_entries_type ON memory_entries(type);
CREATE INDEX idx_memory_entries_synced_to_qdrant ON memory_entries(synced_to_qdrant);
CREATE INDEX idx_memory_entries_modified_at ON memory_entries(modified_at);
CREATE INDEX idx_memory_entries_path ON memory_entries(path);
```

**Política de Retenção:** Forever (local truth antes de Qdrant)

---

### 5. `system_health`

**Propósito:** Métricas de health da infraestrutura (VRAM, disco, containers)
**Criticidade:** MEDIUM
**Owner:** codex (coleta automática)

**Schema:**
```sql
CREATE TABLE IF NOT EXISTS system_health (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    check_type TEXT NOT NULL CHECK(check_type IN ('vram', 'disk', 'container', 'service')),
    check_name TEXT NOT NULL,
    value REAL NOT NULL,
    threshold REAL,
    status TEXT CHECK(status IN ('ok', 'warning', 'critical')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_system_health_check_type ON system_health(check_type);
CREATE INDEX idx_system_health_created_at ON system_health(created_at);
```

**Política de Retenção:** 7d hot (métricas), 30d warm (trends)

---

## Tabelas Futuras

- `adr_execution_log` (quando ADR executar ações automáticas)
- `model_inference_metrics` (qwen3.5:9b latência + VRAM)
- `embedding_quality` (bge-m3 vector variance, dedup checks)

---

## Validação

```bash
# Listar tabelas reais
sqlite3 ~/.aurelia/data/aurelia.db ".tables"

# Exportar schema
sqlite3 ~/.aurelia/data/aurelia.db ".schema" > /tmp/schema.sql

# Diff contra registro
diff -u docs/schema-registry-sqlite.md <(sqlite3 ~/.aurelia/data/aurelia.db ".schema")
```

---

## Backup & Recovery

**Automatizado:**
- ZFS snapshots: hourly, daily, weekly
- WAL mode: ✅ (write-ahead logging para integridade)

**Manual:**
```bash
sqlite3 ~/.aurelia/data/aurelia.db ".dump" | gzip > backup-$(date +%Y%m%d).sql.gz
```

---

## Links

- [Schema Registry — PostgreSQL](./schema-registry-postgres.md)
- [Domain Ownership Table](./data-governance-domain-ownership.md)
- [Data Lifecycle Policy](./data-governance-lifecycle.md)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
