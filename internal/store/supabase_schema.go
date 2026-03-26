package store

// schemaSQL creates the minimal Supabase tables needed by Aurelia.
// Only long-lived, queryable data lives here — runtime state stays in SQLite.
// All statements use CREATE TABLE IF NOT EXISTS so Migrate() is idempotent.
const schemaSQL = `
-- knowledge_items: curated, long-lived facts indexed into Qdrant.
-- Source of truth for semantic search; SQLite holds only runtime cache.
CREATE TABLE IF NOT EXISTS knowledge_items (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id          TEXT        NOT NULL DEFAULT 'aurelia',
    canonical_bot_id TEXT       NOT NULL,
    domain          TEXT        NOT NULL,
    kind            TEXT        NOT NULL CHECK (kind IN (
                        'message','note','decision','knowledge',
                        'task_summary','event','obsidian_note')),
    content         TEXT        NOT NULL,
    metadata        JSONB       NOT NULL DEFAULT '{}',
    source_system   TEXT        NOT NULL,
    source_id       TEXT        NOT NULL,
    version         INTEGER     NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (source_system, source_id)
);

CREATE INDEX IF NOT EXISTS ki_bot_domain ON knowledge_items (canonical_bot_id, domain);
CREATE INDEX IF NOT EXISTS ki_source ON knowledge_items (source_system, source_id);
CREATE INDEX IF NOT EXISTS ki_kind ON knowledge_items (kind);

-- component_status: historical snapshots of system health.
-- Queryable history; in-memory / SQLite hold only the latest snapshot.
CREATE TABLE IF NOT EXISTS component_status (
    id          BIGSERIAL   PRIMARY KEY,
    component   TEXT        NOT NULL,
    status      TEXT        NOT NULL CHECK (status IN ('healthy','degraded','offline','stale')),
    summary     TEXT,
    latency_ms  INTEGER,
    last_ok_at  TIMESTAMPTZ,
    last_error  TEXT,
    source      TEXT,
    details     JSONB       NOT NULL DEFAULT '{}',
    checked_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS cs_component_checked ON component_status (component, checked_at DESC);
`
