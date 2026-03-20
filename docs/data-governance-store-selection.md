---
title: Store Selection Matrix
description: Árvore de decisão para escolher entre SQLite, PostgreSQL, Qdrant, Supabase
owner: codex
updated: 2026-03-20
---

# Store Selection Matrix — Where Does Data Live?

**Purpose:** Decisão estruturada sobre qual store usar para cada tipo de dado
**Authority:** ADR-20260319-Polish-Governance-All / DATA_GOVERNANCE
**Updated:** 2026-03-20

---

## Decision Tree

```
┌─ Data at Rest?
│
├─YES ─┬─ Operacional/Stateful (gateway, task history, metrics)?
│      │
│      ├─YES ─┬─ Hot path (<100ms latency)?
│      │      │
│      │      ├─YES → SQLite (~/.aurelia/data/aurelia.db)
│      │      │       [gateway_route_states, voice_events, cron_tasks, system_health]
│      │      │
│      │      └─NO ──→ PostgreSQL (appropriate instance)
│      │               [ai_context.*, workflow logs if >30d hot needed]
│      │
│      └─NO ──┬─ Secretos/Credenciais?
│             │
│             ├─YES → KeePassXC Vault (+ encrypted Postgres columns)
│             │       [API keys, Postgres passwords, Cloudflare tokens]
│             │
│             └─NO ──┬─ Semanticamente Pesquisável (ML/embedding)?
│                    │
│                    ├─YES → Qdrant (bge-m3 vectors)
│                    │       [repository_memory, conversation_context]
│                    │
│                    └─NO ──→ PostgreSQL (default for relational)
│                             [users, sessions, configs, audit logs]
│
└─NO ──→ Transient/Session State
         [memory_entries cache before Qdrant sync]
         └─→ SQLite (local truth) → Postgres (sync target)
```

---

## Store Comparison Matrix

| Critério | SQLite | PostgreSQL | Qdrant | Supabase | KeePassXC |
|---|---|---|---|---|---|
| **Latência** | <5ms | 10-50ms | 20-100ms | 50-200ms | N/A (offline) |
| **Durabilidade** | ZFS snapshots | Hourly snapshots | Regenerable | Cloud managed | Encrypted file |
| **Escala** | GB (single file) | TB (cluster) | 100M+ vectors | SaaS limit | MB (passwords) |
| **Searchabilidade** | SQL WHERE | SQL + indexes | Semantic (cosine) | SQL + PostGIS | N/A |
| **Sincronia** | N/A (local truth) | Multi-instance | Derived from Postgres | Optional mirror | Manual backup |
| **Segurança** | File-level (ZFS) | RBAC + encryption | Network isolation | OAuth2 | Master password |
| **Operacional** | 1 process | Docker + replicação | Docker + collection | SaaS | Native app |

---

## Store Assignment Rules

### SQLite (`~/.aurelia/data/aurelia.db`)
**Use when:**
- ✅ Local-first (não need replication)
- ✅ Latency crítica (<10ms)
- ✅ Operacional/stateful (gateway, voice)
- ✅ Crash-safe é essencial (ZFS-backed)

**Tables:**
- `gateway_route_states` — Route health check status
- `voice_events` — TTS synthesis log
- `cron_tasks` — Task execution history
- `memory_entries` — Local cache antes sync
- `system_health` — Infrastructure metrics

**Lifecycle:** Operational (hot), never deleted unless explicitly cleared

---

### PostgreSQL (4 Instâncias)

#### n8n Instance
**Use when:**
- ✅ Automação workflows (n8n native store)
- ✅ Forever data (workflows, credentials, webhooks)
- ✅ Multi-user access (credentials shared)

**Tables:**
- `workflows` — Automation definitions (CRITICAL)
- `workflow_history` — Execution audit trail
- `credentials` — Encrypted API keys
- `webhooks` — Registered endpoints

---

#### supabase-db Instance
**Use when:**
- ✅ User sessions (Telegram bot)
- ✅ Temporary data (7d-30d retention)
- ✅ Optional mirror of SQLite

**Tables:**
- `sessions` — Active user sessions
- `messages` — Chat history (30d hot)
- `users` — User registry
- `auth.users` — Supabase auth (managed)

---

#### litellm-db Instance
**Use when:**
- ✅ LLM router state (model configs, API keys)
- ✅ Usage tracking (30d retention)
- ✅ Error logs (troubleshooting)

**Tables:**
- `api_keys` — LLM provider credentials
- `usage_logs` — Token usage by model
- `model_configs` — Router configuration
- `error_logs` — Fallback decisions

---

#### dev Instance
**Use when:**
- ✅ Experimentation/sandboxing
- ✅ Testing schema changes
- ✅ No production constraints

**Tables:**
- Any (no restrictions)

---

### Qdrant (`localhost:6333`)

**Use when:**
- ✅ Semantic search required (embedding-based)
- ✅ Vector dimension <= 384 (bge-m3 constraint)
- ✅ Derived data (regenerable from Postgres)
- ✅ Offline access for Aurelia bot

**Collection:**
- `repository_memory` — bge-m3 384-dim embeddings
  - Payload: file_path, memory_type, owner, synced_at, text_content
  - Updated: Every 5min from Postgres via fiscal cron

**Lifecycle:** Derived (regenerable), never source of truth

---

### Supabase (Optional Mirror)

**Use when:**
- ✅ Replication desired (SQLite → Postgres → Supabase)
- ✅ Cloud backup wanted
- ✅ Mobile app needs online access

**Note:** Currently optional, not required by ADR-20260319

**Sync Flow:** SQLite (local) → supabase-db (Postgres) → Supabase Cloud (mirror)

---

### KeePassXC Vault (`/srv/data/vault/aurelia.kdbx`)

**Use for ONLY:**
- ✅ Passwords/API keys/secrets
- ✅ Credentials never in plaintext files
- ✅ Master password required to access

**Items:**
- Postgres credentials (n8n, supabase, litellm instances)
- Cloudflare API tokens
- LLM provider API keys
- SSH keys
- GPG keys

**Lifecycle:** Forever (manual rotation quarterly per ADR-20260319 / SECRETS_GOVERNANCE)

---

## Sync Architecture

```
SQLite (Local Truth)
    ↓
    ├─→ cron:every 5min: memory-sync-fiscal.sh --mode fast
    │   └─→ Find changed *.md files in ~/.claude/projects/-home-will-aurelia/memory/
    │
    ├─→ cron:every 15min: memory-sync-fiscal.sh --mode postgres-index
    │   └─→ Push memory_entries → ai_context.memory_entries (Postgres)
    │
    ├─→ cron:every 6am: memory-sync-fiscal.sh --mode validate
    │   └─→ Verify Qdrant vs Postgres consistency
    │
    └─→ cron:Mon 2am: memory-sync-fiscal.sh --mode compact
        └─→ Cleanup old metrics, optimize storage

Postgres (Metadata)
    ↓
    ├─→ Periodic: Stream to Qdrant via bge-m3 embeddings
    │
    └─→ Optional: Mirror to Supabase Cloud (if needed)

Qdrant (Vector Index)
    ↓
    └─→ Used by: Aurelia bot (offline semantic search)
```

---

## Decision Examples

### Example 1: New metric for GPU VRAM usage
**Q:** Where should `system_health` metric live?
- Operacional? ✅ YES (need real-time alerts)
- Hot path? ✅ YES (<10ms for health checks)
- **Decision:** SQLite `system_health` table
- Sync: Via cron to Prometheus metrics

### Example 2: Store voice synthesis latencies
**Q:** Where should `voice_events` live?
- Operacional? ✅ YES (TTS latency tracking)
- Hot path? ✅ YES (5min granularity)
- Analytics needed? Maybe (30d retention)
- **Decision:** SQLite `voice_events` table (primary)
- Extended: Archive to Postgres after 30d hot retention

### Example 3: Conversation memory for Aurelia bot
**Q:** Where should multi-turn conversation context live?
- Semantically searchable? ✅ YES (embedding needed)
- Offline access? ✅ YES (no internet at query time)
- Size? ~100M tokens possible
- **Decision:** Qdrant `repository_memory` collection (bge-m3 vectors)
- Metadata: Postgres `ai_context.memory_entries`
- Source: Markdown files in memory sync directory

### Example 4: API key for litellm
**Q:** Where should OpenAI API key live?
- Is a secret? ✅ YES (never plaintext)
- Frequently accessed? ✅ YES (on every inference)
- **Decision:** KeePassXC vault (accessed at startup)
- Reference: litellm-db stores config only (key_id, not actual key)

---

## Validation

```bash
# Check SQLite tables match decision matrix
sqlite3 ~/.aurelia/data/aurelia.db ".schema" | grep "CREATE TABLE"

# Check Postgres instances have expected schemas
for db in n8n supabase litellm; do
    docker exec postgres-$db psql -c "\dt"
done

# Check Qdrant collection exists
curl -s http://localhost:6333/collections | jq '.result[] | .name'

# Check KeePassXC vault exists and locked
ls -lah /srv/data/vault/aurelia.kdbx
```

---

## Links

- [Domain Ownership Table](./data-governance-domain-ownership.md)
- [Data Lifecycle Policy](./data-governance-lifecycle.md)
- [Schema Registry — SQLite](./schema-registry-sqlite.md)
- [Schema Registry — PostgreSQL](./schema-registry-postgres.md)
- [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
