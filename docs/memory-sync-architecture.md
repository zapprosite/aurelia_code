---
description: Arquitetura de Memory Sync — Como a Aurelia bot acessa code history sem web, usando Qdrant + Postgres.
---

# Memory Sync Architecture — Aurelia Bot Code History

## Visão Geral

A **Aurelia bot** precisa trabalhar com **LLMs pequenas** (não Anthropic, não web) mantendo acesso completo ao **histórico de código e decisões** do repositório. Isso é possível através de **3 camadas de sincronização**:

```
┌──────────────────────────────────────────────────────────────┐
│ Camada 1: Memória Local (Markdown)                          │
│ ~/.claude/projects/.../memory/                              │
│ docs/adr/                                                   │
│ .context/runbooks/, plans/                                  │
└──────────────────┬───────────────────────────────────────────┘
                   │ [SYNC CRON — Fiscal]
                   ├─→ Detectar novos/modificados
                   ├─→ Gerar embeddings (bge-m3)
                   └─→ Indexar metadata
                   │
┌──────────────────▼───────────────────────────────────────────┐
│ Camada 2: Qdrant (Semantic Search)                          │
│ • Collection: repository_memory                             │
│ • Embedding: bge-m3 (384-dim, local)                        │
│ • Query: busca semântica (sem rede)                         │
└──────────────────┬───────────────────────────────────────────┘
                   │ [Dual-write]
                   │
┌──────────────────▼───────────────────────────────────────────┐
│ Camada 3: Postgres (Metadata + Full-Text)                   │
│ • Tabelas: memory_entries, adr_registry, sync_log           │
│ • Índices: path, type, owner, tags, full-text              │
│ • Query: busca estruturada + relacionamentos               │
└──────────────────┬───────────────────────────────────────────┘
                   │
┌──────────────────▼───────────────────────────────────────────┐
│ Aurelia Bot                                                  │
│ • Pergunta de usuário → embedding local                     │
│ • Busca Qdrant (semântica)                                  │
│ • Enriquece com Postgres (metadata)                         │
│ • Responde com LLM pequena + contexto local                │
└──────────────────────────────────────────────────────────────┘
```

---

## Camada 1: Memória Local

**Onde:** `~/.claude/projects/-home-will-aurelia/memory/`, `docs/adr/`, `.context/`

**Formato:** Markdown com frontmatter YAML

### Exemplo: Memory Entry

```markdown
---
name: ADR Polish Governance All
description: ADR mestra de governança industrial
type: project_memory
owner: codex
tags: [governance, homelab, critical]
---

# ADR-20260319-Polish-Governance-All

Conteúdo da memória...
```

**Propriedades Detectadas:**
- `path` — arquivo local
- `type` — project_memory, adr, runbook, plan
- `owner` — codex, humano, gemini
- `tags` — lista de tags
- `mtime` — modificado quando
- `hash` — SHA256 para detectar mudanças

---

## Camada 2: Qdrant (Busca Semântica)

**O quê:** Collection `repository_memory` armazena embeddings

**Model:** BGE-M3 (384 dimensões, local, multilíngue)

### Schema

```json
{
  "collection_name": "repository_memory",
  "vectors": {
    "size": 384,
    "distance": "Cosine"
  },
  "payload_schema": {
    "path": { "type": "text" },
    "type": { "type": "keyword" },
    "owner": { "type": "keyword" },
    "tags": { "type": "keyword" },
    "created_at": { "type": "datetime" },
    "modified_at": { "type": "datetime" }
  }
}
```

### Exemplo: Query Semântica

```bash
# Usuário pergunta: "Como foi decidido o vault de secrets?"
# Aurelia embedding: [0.123, 0.456, ...] (384-dim)

curl -X POST http://localhost:6333/collections/repository_memory/points/search \
  -H "Content-Type: application/json" \
  -d '{
    "vector": [0.123, 0.456, ...],
    "limit": 5,
    "with_payload": true,
    "filter": {
      "must": [
        {"key": "tags", "match": {"value": "secrets"}}
      ]
    }
  }'

# Resposta:
# [
#   {
#     "id": "mem_polish_governance",
#     "score": 0.92,
#     "payload": {
#       "path": "~/.claude/.../memory/adr_polish_governance_all.md",
#       "type": "project_memory",
#       "tags": ["governance", "secrets"],
#       "created_at": "2026-03-19T23:45:00Z"
#     }
#   },
#   ...
# ]
```

---

## Camada 3: Postgres (Indexação Estruturada)

**O quê:** Tabelas para metadata, relacionamentos, full-text search

### Tabelas Principais

#### `memory_entries`
```sql
CREATE TABLE memory_entries (
    id VARCHAR PRIMARY KEY,
    path TEXT UNIQUE,
    type VARCHAR,  -- project_memory, adr, runbook, plan
    owner VARCHAR,
    tags TEXT[],   -- ARRAY
    embedding_id VARCHAR,  -- referência Qdrant
    content_hash VARCHAR,  -- SHA256 para detectar mudanças
    size_bytes INT,
    created_at TIMESTAMP,
    modified_at TIMESTAMP,
    synced_at TIMESTAMP,
    INDEX idx_type (type),
    INDEX idx_owner (owner),
    FULLTEXT INDEX idx_path (path)
);
```

#### `adr_registry`
```sql
CREATE TABLE adr_registry (
    id VARCHAR PRIMARY KEY,
    slug VARCHAR UNIQUE,
    title TEXT,
    status VARCHAR,  -- proposed, in_progress, accepted
    description TEXT,
    owner VARCHAR,
    memory_entry_id VARCHAR REFERENCES memory_entries(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (memory_entry_id) REFERENCES memory_entries(id)
);
```

#### `sync_log`
```sql
CREATE TABLE sync_log (
    id SERIAL PRIMARY KEY,
    sync_type VARCHAR,  -- fast, postgres_index, validate, compact
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    files_processed INT,
    embeddings_created INT,
    errors INT,
    status VARCHAR,  -- success, failed, partial
    details JSON
);
```

### Exemplo: Query Estruturada

```sql
-- "Quais decisões (ADRs) foram tomadas por codex?"
SELECT
    ar.id, ar.title, ar.status,
    me.path, me.modified_at,
    me.tags
FROM adr_registry ar
JOIN memory_entries me ON ar.memory_entry_id = me.id
WHERE ar.owner = 'codex'
  AND ar.status = 'accepted'
ORDER BY ar.updated_at DESC;
```

---

## Fiscal: Sincronização Automática

**Responsabilidade:** Manter Qdrant + Postgres sincronizados com a memória local

### Cron Schedule

| Frequência | Script | Modo | O que faz |
|---|---|---|---|
| A cada 5 min | `memory-sync-fiscal.sh` | `fast` | Detectar novos/mod + embedding |
| A cada 15 min | `memory-sync-fiscal.sh` | `postgres-index` | Atualizar indices Postgres |
| Diariamente 6am | `memory-sync-fiscal.sh` | `validate` | Verificar integridade + relatório |
| Segunda 2am | `memory-sync-fiscal.sh` | `compact` | Cleanup + otimização |

### Script: `memory-sync-fiscal.sh`

```bash
#!/usr/bin/env bash
# ~/.aurelia/scripts/memory-sync-fiscal.sh

set -euo pipefail

MODE="${1:-fast}"
LOG_FILE="$HOME/.aurelia/logs/memory-sync-fiscal.log"

function log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" >> "$LOG_FILE"
}

function sync_fast() {
    log "Starting fast sync (new files + embedding)..."

    # 1. Varrer filesystem
    for file in ~/.claude/projects/-home-will-aurelia/memory/*.md; do
        if [ -f "$file" ]; then
            hash=$(sha256sum "$file" | cut -d' ' -f1)
            path=$(basename "$file")

            # 2. Gerar embedding (bge-m3)
            embedding=$(curl -s http://localhost:8000/embed \
                -H "Content-Type: application/json" \
                -d "{\"texts\": [\"$(cat $file)\"]}" | jq '.embeddings[0]')

            # 3. Upsert Qdrant
            curl -s -X POST http://localhost:6333/collections/repository_memory/points/upsert \
                -H "Content-Type: application/json" \
                -d "{
                  \"points\": [{
                    \"id\": \"mem_$(echo $path | sha256sum | cut -c1-16)\",
                    \"vector\": $embedding,
                    \"payload\": {
                      \"path\": \"$file\",
                      \"type\": \"project_memory\",
                      \"hash\": \"$hash\",
                      \"modified_at\": \"$(stat -c %y $file)\"
                    }
                  }]
                }"

            log "Embedded: $path (hash: $hash)"
        fi
    done
}

function sync_postgres_index() {
    log "Starting postgres index sync..."

    # Atualizar indices
    psql -d aurelia -c "REINDEX TABLE ai_context.memory_entries;"
    psql -d aurelia -c "REINDEX TABLE ai_context.adr_registry;"

    log "Postgres indices updated"
}

function validate() {
    log "Starting validation..."

    # Comparar counts Qdrant vs Postgres
    qdrant_count=$(curl -s http://localhost:6333/collections/repository_memory/points/count | jq '.result.count')
    postgres_count=$(psql -d aurelia -t -c "SELECT COUNT(*) FROM ai_context.memory_entries;")

    if [ "$qdrant_count" -eq "$postgres_count" ]; then
        log "✅ Validation passed: Qdrant=$qdrant_count, Postgres=$postgres_count"
    else
        log "⚠️  Validation warning: Count mismatch! Qdrant=$qdrant_count, Postgres=$postgres_count"
    fi
}

function compact() {
    log "Starting compaction..."

    # Cleanup entries antigas (>90 dias)
    psql -d aurelia -c "
        DELETE FROM ai_context.memory_entries
        WHERE modified_at < NOW() - INTERVAL '90 days'
          AND type IN ('plan', 'temp');
    "

    log "Compaction complete"
}

case "$MODE" in
    fast)
        sync_fast
        ;;
    postgres-index)
        sync_postgres_index
        ;;
    validate)
        validate
        ;;
    compact)
        compact
        ;;
    *)
        echo "Unknown mode: $MODE"
        exit 1
        ;;
esac

log "Sync complete: $MODE"
```

---

## Fluxo: Como Aurelia Usa Isso

```go
// 1. Usuário pergunta algo no Telegram/CLI
query := "Como foi decidido armazenar secrets?"

// 2. Aurelia embedda a pergunta (bge-m3, local)
queryEmbedding := embedModel.Embed(query)  // [0.123, 0.456, ...]

// 3. Busca semântica no Qdrant
results := qdrant.Search(queryEmbedding, limit: 5)
// Retorna: [mem_polish_governance, mem_keepassxc_tutorial, ...]

// 4. Enriquece com Postgres
for _, result := range results {
    metadata := postgres.Query(`
        SELECT owner, type, tags, created_at
        FROM memory_entries WHERE id = ?
    `, result.ID)
    result.Metadata = metadata
}

// 5. Monta contexto e passa para LLM pequena (e.g., Mistral 7B)
context := formatContext(results)
response := llmSmall.Generate(
    query,
    context,
    maxTokens: 200,
)

// 6. Responde ao usuário
fmt.Printf("🤖 Aurelia: %s\n", response)
```

---

## Benefícios

| Aspecto | Vantagem |
|---|---|
| **Sem Web** | Qdrant + Postgres são locais; nenhuma chamada externa |
| **LLMs Pequenas** | BGE-M3 é local; retrieval semântico funciona offline |
| **Code History** | Memória + ADR + runbooks sempre sincronizados |
| **Escalável** | Crons distribuem carga (5, 15, 60, 1440 min) |
| **Auditável** | Sync_log registra cada operação; fácil debug |
| **Seguro** | Sem API keys expostas; tudo em `/srv/data/` (ZFS) |

---

## Próximos Passos

1. ✅ Criar skill `/memory-sync-vector-db`
2. 📝 Documentar (este arquivo)
3. 🔧 Implementar `memory-sync-fiscal.sh` (script)
4. ⏰ Configurar crons em aurelia.service
5. 🧪 Testar sincronização end-to-end
6. 📊 Monitorar via `sync_log` no Postgres
