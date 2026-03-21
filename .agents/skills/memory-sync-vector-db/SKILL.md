---
type: skill
name: Memory Sync — Vector DB Architecture
description: Sincronizar memória do repositório (markdown) → Qdrant (embeddings) + Postgres (metadata). Permite Aurelia bot trabalhar com LLMs pequenas sem web.
skillSlug: memory-sync-vector-db
phases: [E, V]
generated: 2026-03-19
status: unfilled
scaffoldVersion: "2.0.0"
---

# Skill: memory-sync-vector-db

## Propósito

**Arquitetura de Memory Sync para Aurelia Bot**

```
┌─────────────────────────────────────────────────────────────┐
│ Repositório Aurelia — Code History + Context               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ~/.claude/projects/.../memory/    (markdown files)        │
│  docs/adr/                         (decisões arquiteto)    │
│  .context/runbooks/                (procedimentos)         │
│         │                                                  │
│         ├──→ [SYNC CRON] ────────────────────────┐        │
│         │                                         │        │
│         ├─→ Qdrant (conversation_memory)         │        │
│         │   ├─ collection: "repository_memory"   │        │
│         │   ├─ embedding: bge-m3                 │        │
│         │   └─ payload: {memory_id, path, ...}  │        │
│         │                                        ↓        │
│         └─→ Postgres (ai_context schema)         │        │
│             ├─ table: memory_entries             │        │
│             ├─ table: adr_registry               │        │
│             └─ table: code_history               │        │
│                                                  ↓        │
│         Aurelia Bot ← Acesso local, LLM pequena           │
│         (sem web, sem Anthropic API)                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Fluxo de Sincronização

1. **Coleta (Fiscal - 5min):**
   - Varrer `~/.claude/projects/-home-will-aurelia/memory/`
   - Varrer `docs/adr/`
   - Varrer `.context/` (runbooks, plans)
   - Detectar novos/modificados (mtime, hash)

2. **Embedding (10min):**
   - Enviar textos para bge-m3 (local)
   - Gerar embeddings 384-dim
   - Armazenar no Qdrant

3. **Indexação (Postgres - 15min):**
   - Registrar metadata (caminho, tipo, owner, tags)
   - Criar referências cruzadas
   - Atualizar full-text search index

4. **Validação (diária - 6am):**
   - Verificar integridade Qdrant ↔ Postgres
   - Limpar duplicatas
   - Gerar relatório de cobertura

## Invocações

```bash
# Trigger manual de sync
/memory-sync-vector-db --sync-now

# Status da sincronização
/memory-sync-vector-db --status

# Validar integridade
/memory-sync-vector-db --validate

# Listar embeddings recentes
/memory-sync-vector-db --list-recent --limit 10
```

## Crons Aurelia (Fiscal Automático)

```bash
# A cada 5 min: coleta + embedding (incremental)
*/5 * * * * /home/will/aurelia/scripts/memory-sync-fiscal.sh --mode fast

# A cada 15 min: indexação Postgres (atualizar metadata)
*/15 * * * * /home/will/aurelia/scripts/memory-sync-fiscal.sh --mode postgres-index

# Diariamente 6am: validação + relatório
0 6 * * * /home/will/aurelia/scripts/memory-sync-fiscal.sh --mode validate

# Semanalmente (segunda): compactação + cleanup
0 2 * * 1 /home/will/aurelia/scripts/memory-sync-fiscal.sh --mode compact
```

## Integração com Aurelia Bot

```go
// aurelia/internal/vector/memory_sync.go (pseudocódigo)

type MemorySyncFiscal struct {
    qdrantClient *qdrant.Client
    postgresDB   *sql.DB
    watchDir     string
    lastSync     time.Time
}

func (m *MemorySyncFiscal) SyncMemory() error {
    // 1. Varrer filesystem
    files := m.scanMemoryFiles()

    // 2. Para cada arquivo: gerar embedding
    for _, file := range files {
        content := readFile(file)
        embedding := m.embedWithBgeM3(content)

        // 3. Salvar no Qdrant
        m.qdrant.Upsert("repository_memory", &Point{
            ID:      hashID(file),
            Vector:  embedding,
            Payload: map[string]any{
                "path":      file,
                "type":      getType(file),
                "modified":  file.ModTime(),
                "owner":     extractOwner(content),
            },
        })

        // 4. Registrar no Postgres
        m.postgres.Exec(`
            INSERT INTO memory_entries (id, path, type, embedding_id, synced_at)
            VALUES ($1, $2, $3, $4, NOW())
            ON CONFLICT(path) DO UPDATE SET synced_at = NOW()
        `)
    }

    m.lastSync = time.Now()
    return nil
}

// Aurelia bot pode consultar sem web:
func (a *AureliaBot) QueryMemory(ctx context.Context, question string) ([]string, error) {
    // 1. Embedar pergunta localmente
    questionEmbedding := a.embedWithBgeM3(question)

    // 2. Buscar no Qdrant (semantic search)
    results := a.qdrant.Search("repository_memory", questionEmbedding, limit: 5)

    // 3. Enriquecer com Postgres metadata
    for _, result := range results {
        meta := a.postgres.QueryRow(`
            SELECT owner, type, created_at FROM memory_entries WHERE id = $1
        `, result.ID).Scan(...)
        result.Metadata = meta
    }

    return results, nil
}
```

## Dados Armazenados

### Qdrant Collection: `repository_memory`

```json
{
  "id": "mem_adr_polish_governance",
  "vector": [0.123, 0.456, ...],  // 384-dim bge-m3
  "payload": {
    "path": "~/.claude/projects/-home-will-aurelia/memory/adr_polish_governance_all.md",
    "type": "project_memory",
    "owner": "codex",
    "tags": ["governance", "homelab", "critical"],
    "created_at": "2026-03-19T23:45:00Z",
    "modified_at": "2026-03-19T23:50:00Z",
    "size_bytes": 2048
  }
}
```

### Postgres Schema

```sql
-- Tabela: ai_context.memory_entries
CREATE TABLE memory_entries (
    id VARCHAR PRIMARY KEY,
    path TEXT UNIQUE,
    type VARCHAR,  -- project_memory, adr, runbook, plan
    owner VARCHAR,
    tags TEXT[],
    embedding_id VARCHAR,
    content_hash VARCHAR,
    created_at TIMESTAMP,
    modified_at TIMESTAMP,
    synced_at TIMESTAMP,
    size_bytes INT
);

-- Tabela: ai_context.adr_registry
CREATE TABLE adr_registry (
    id VARCHAR PRIMARY KEY,
    slug VARCHAR UNIQUE,
    title TEXT,
    status VARCHAR,  -- proposed, in_progress, accepted
    owner VARCHAR,
    memory_entry_id VARCHAR REFERENCES memory_entries(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Tabela: ai_context.sync_log
CREATE TABLE sync_log (
    id SERIAL PRIMARY KEY,
    sync_type VARCHAR,  -- fast, postgres_index, validate, compact
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    files_processed INT,
    embeddings_created INT,
    errors INT,
    status VARCHAR  -- success, failed, partial
);
```

## Fiscal: Monitoramento Contínuo

Cron executa `memory-sync-fiscal.sh` em diferentes frequências:

- **`--mode fast` (5min):** Detecção rápida de novos arquivos + embedding
- **`--mode postgres-index` (15min):** Atualizar índices no Postgres
- **`--mode validate` (diária 6am):** Verificar integridade + relatório
- **`--mode compact` (semanal segunda):** Limpeza + otimização

Logs em `~/.aurelia/logs/memory-sync-fiscal.log`

## Benefícios para Aurelia

1. **Sem Web:** Aurelia consulta Qdrant localmente (não precisa Anthropic API)
2. **LLMs Pequenas:** Embedding local (bge-m3) + retrieval semântico
3. **Code History:** Contexto completo do repositório sempre disponível
4. **Histórico de Decisões:** ADRs e memória sincronizados automaticamente
5. **Escalável:** Crons distribuem carga (5min, 15min, diária, semanal)

## Referências

- [ADR-20260319-Polish-Governance-All](../../docs/adr/ADR-20260319-Polish-Governance-All.md)
- [Memory Architecture Doc](../../docs/memory-sync-architecture.md) (será criado)
- [Aurelia Internal: vector](../../internal/vector/)
- [Qdrant API](https://qdrant.tech/documentation/)
- BGE-M3: Local embedding model
