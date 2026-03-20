---
title: Schema Registry — PostgreSQL
description: Registro completo de tabelas, colunas, índices e políticas de retenção dos 4 PostgreSQL instâncias
owner: codex
updated: 2026-03-20
---

# Schema Registry — PostgreSQL

**Instances:** 4 (n8n, supabase-db, litellm-db, dev)
**High Availability:** Replicação por serviço
**Retenção:** 30d hot, 90d warm, >90d archive

---

## Instance 1: `n8n` (Automations DB)

**Container:** `n8n` (n8n.io)
**Port:** 5432 (internal)
**Backup:** Docker volume snapshots

### Tabelas Críticas

| Tabela | Propósito | Owner | Retenção |
|--------|-----------|-------|----------|
| `workflows` | Definições de automações | codex | Forever |
| `workflow_history` | Execução de workflows (audit) | codex | 90d |
| `credentials` | API keys para serviços (encrypted) | humano | Forever |
| `webhooks` | Endpoints de webhook registrados | codex | Forever |

### Validação
```bash
docker exec n8n psql -U root -d n8n -c "\dt"
```

---

## Instance 2: `supabase-db` (Mirror + Sessions)

**Container:** `supabase-postgres`
**Port:** 5432 (internal)
**Replicação:** Mirror de SQLite (futuro) + sessions diretas

### Tabelas

| Tabela | Propósito | Owner | Retenção |
|--------|-----------|-------|----------|
| `sessions` | User sessions (Telegram bot) | codex | 7d active, 30d archive |
| `messages` | Chat history | codex | 30d hot, 90d warm |
| `users` | User registry (if needed) | humano | Forever |
| `auth.users` | Supabase auth (managed) | supabase | Forever |

### Validação
```bash
docker exec supabase-postgres psql -U postgres -d postgres -c "\dt public.*"
```

---

## Instance 3: `litellm-db` (LLM Gateway)

**Container:** `litellm`
**Port:** 5432 (internal)
**Purpose:** Proxy de modelos LLM (qwen3.5:9b fallback)

### Tabelas

| Tabela | Propósito | Owner | Retenção |
|--------|-----------|-------|----------|
| `api_keys` | LLM provider keys | humano | Forever |
| `usage_logs` | Token usage by model | codex | 30d |
| `model_configs` | Router config (qwen, gemini, etc) | codex | Forever |
| `error_logs` | Fallback decisions | codex | 7d |

### Validação
```bash
docker exec litellm psql -U postgres -d litellm -c "\dt"
```

---

## Instance 4: `dev` (Development/Testing)

**Port:** 5432 (optional, for dev workflows)
**Propósito:** Sandbox para experimentação

### Tabelas
Sem restrição (ambiente de desenvolvimento)

---

## Tabelas Futuras (Planejadas)

### `ai_context.memory_entries` (Shared across DBs)

```sql
CREATE TABLE IF NOT EXISTS ai_context.memory_entries (
    id SERIAL PRIMARY KEY,
    path TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    owner TEXT,
    tags TEXT[],
    file_hash VARCHAR(64),
    synced_at TIMESTAMP DEFAULT NOW(),
    modified_at TIMESTAMP
);

CREATE INDEX idx_memory_entries_type ON ai_context.memory_entries(type);
CREATE INDEX idx_memory_entries_synced_at ON ai_context.memory_entries(synced_at);
```

### `ai_context.adr_registry` (ADR Tracking)

```sql
CREATE TABLE IF NOT EXISTS ai_context.adr_registry (
    id SERIAL PRIMARY KEY,
    adr_slug VARCHAR(255) NOT NULL UNIQUE,
    status TEXT CHECK(status IN ('proposed', 'in_progress', 'accepted', 'blocked')),
    owner TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### `ai_context.sync_log` (Memory Sync Audit)

```sql
CREATE TABLE IF NOT EXISTS ai_context.sync_log (
    id SERIAL PRIMARY KEY,
    sync_mode TEXT,
    files_processed INTEGER,
    embeddings_created INTEGER,
    duration_ms INTEGER,
    status TEXT CHECK(status IN ('success', 'partial', 'failure')),
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## Validação Global

```bash
# Conectar a cada instance e listar schemas
for db in n8n supabase-postgres litellm; do
    echo "=== $db ==="
    docker exec $db psql -U postgres -l
    docker exec $db psql -U postgres -c "\dt" || true
done
```

---

## Links

- [Schema Registry — SQLite](./schema-registry-sqlite.md)
- [Domain Ownership Table](./data-governance-domain-ownership.md)
- [Qdrant Collection Contract](./qdrant-collection-contract.md)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
