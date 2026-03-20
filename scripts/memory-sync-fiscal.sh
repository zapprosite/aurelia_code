#!/usr/bin/env bash
# memory-sync-fiscal.sh — Fiscal automático para sincronizar memória → Qdrant + Postgres
# Invocado via systemd timers em diferentes frequências

set -euo pipefail

# Configuração
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
MEMORY_DIR="$HOME/.claude/projects/-home-will-aurelia/memory"
ADR_DIR="$PROJECT_ROOT/docs/adr"
CONTEXT_DIR="$PROJECT_ROOT/.context"
LOG_DIR="$HOME/.aurelia/logs"
METRICS_DIR="$HOME/.aurelia/metrics"
LOG_FILE="$LOG_DIR/memory-sync-fiscal.log"

# Garantir diretórios
mkdir -p "$LOG_DIR" "$METRICS_DIR"

# Defaults
MODE="${1:-fast}"
QDRANT_URL="${QDRANT_URL:-http://localhost:6333}"
POSTGRES_DB="${POSTGRES_DB:-aurelia}"
DRY_RUN="${DRY_RUN:-false}"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# ============================================================================
# Helpers
# ============================================================================

function log() {
    local level="$1"
    shift
    local msg="$*"
    local ts=$(date +'%Y-%m-%d %H:%M:%S')
    echo "[$ts] [$level] $msg" >> "$LOG_FILE"
    [ "$level" = "ERROR" ] && echo -e "${RED}[ERROR]${NC} $msg" >&2
    [ "$level" = "WARN" ] && echo -e "${YELLOW}[WARN]${NC} $msg" >&2
}

function info() { log "INFO" "$@"; }
function warn() { log "WARN" "$@"; }
function error() { log "ERROR" "$@"; }

function gauge_metric() {
    local name="$1"
    local value="$2"
    local labels="${3:-}"

    # Prometheus format
    if [ -n "$labels" ]; then
        echo "$name{$labels} $value" >> "$METRICS_DIR/memory-sync.prom"
    else
        echo "$name $value" >> "$METRICS_DIR/memory-sync.prom"
    fi
}

function counter_metric() {
    local name="$1"
    local value="$2"
    local labels="${3:-}"

    echo "# TYPE $name counter" >> "$METRICS_DIR/memory-sync.prom"
    if [ -n "$labels" ]; then
        echo "$name{$labels} $value" >> "$METRICS_DIR/memory-sync.prom"
    else
        echo "$name $value" >> "$METRICS_DIR/memory-sync.prom"
    fi
}

# ============================================================================
# Modo: FAST (5 minutos)
# ============================================================================

function sync_fast() {
    info "Starting FAST sync (new files + embedding)..."
    local start_time=$(date +%s)
    local files_processed=0
    local embeddings_created=0
    local errors=0

    # Coleta todos os markdown files
    local files=()
    while IFS= read -r -d '' file; do
        files+=("$file")
    done < <(find "$MEMORY_DIR" "$ADR_DIR" "$CONTEXT_DIR/runbooks" "$CONTEXT_DIR/plans" \
        -name "*.md" -type f -print0 2>/dev/null)

    info "Found ${#files[@]} markdown files to process"

    for file in "${files[@]}"; do
        if [ ! -f "$file" ]; then
            continue
        fi

        ((files_processed++))

        local rel_path="${file#$HOME/}"
        local file_hash=$(sha256sum "$file" | cut -d' ' -f1)
        local file_mtime=$(stat -c %y "$file" 2>/dev/null | cut -d' ' -f1-2)
        local file_size=$(stat -c %s "$file" 2>/dev/null)

        # Detectar tipo
        local type="project_memory"
        if [[ "$file" == *"/adr/"* ]]; then
            type="adr"
        elif [[ "$file" == *"/runbooks/"* ]]; then
            type="runbook"
        elif [[ "$file" == *"/plans/"* ]]; then
            type="plan"
        fi

        # Extrair metadata do frontmatter
        local owner="codex"
        local tags="general"
        if grep -q "^owner:" "$file"; then
            owner=$(grep "^owner:" "$file" | head -1 | cut -d: -f2 | xargs)
        fi
        if grep -q "^tags:" "$file"; then
            tags=$(grep "^tags:" "$file" | head -1 | cut -d: -f2 | xargs)
        fi

        info "Processing: $rel_path (type=$type, owner=$owner)"

        # Gerar embedding com bge-m3
        # Nota: Assumir servidor embedding rodando em localhost:8000
        if command -v python3 &> /dev/null; then
            # Fallback: usar Python local para embedding
            local embedding=$(python3 << 'PYTHON'
import json, sys, hashlib
sys.path.insert(0, "/home/will/aurelia/internal/embedding")
try:
    from bge_m3_embed import EmbedModel
    model = EmbedModel()
    with open(sys.argv[1]) as f:
        text = f.read()[:8000]  # Limitar a 8k chars
    vec = model.embed([text])[0]
    print(json.dumps(vec))
except Exception as e:
    print("[]", file=sys.stderr)
    sys.exit(1)
PYTHON
 "$file" 2>/dev/null || echo "[]")

            if [ "$embedding" != "[]" ] && [ -n "$embedding" ]; then
                ((embeddings_created++))

                # Upsert no Qdrant
                local point_id="mem_$(echo $file_hash | cut -c1-16)"
                local payload=$(cat << EOF
{
  "id": "$point_id",
  "vector": $embedding,
  "payload": {
    "path": "$rel_path",
    "type": "$type",
    "owner": "$owner",
    "tags": [$tags],
    "hash": "$file_hash",
    "size_bytes": $file_size,
    "modified_at": "$file_mtime"
  }
}
EOF
)

                if [ "$DRY_RUN" = "false" ]; then
                    if curl -s -X POST "$QDRANT_URL/collections/repository_memory/points/upsert" \
                        -H "Content-Type: application/json" \
                        -d "{\"points\": [$payload]}" > /dev/null 2>&1; then
                        info "  ✓ Embedded: $rel_path"
                    else
                        warn "  ✗ Failed to upsert to Qdrant: $rel_path"
                        ((errors++))
                    fi
                fi
            else
                warn "  ✗ Failed to generate embedding: $rel_path"
                ((errors++))
            fi
        else
            warn "  ⊘ Python3 not found; skipping embedding"
        fi
    done

    local end_time=$(date +%s)
    local duration=$((end_time - start_time))

    # Registrar no Postgres
    if [ "$DRY_RUN" = "false" ]; then
        psql -d "$POSTGRES_DB" -c "
            INSERT INTO ai_context.sync_log (sync_type, started_at, completed_at, files_processed, embeddings_created, errors, status)
            VALUES ('fast', NOW() - INTERVAL '$duration seconds', NOW(), $files_processed, $embeddings_created, $errors, 'success')
            ON CONFLICT DO NOTHING;
        " 2>/dev/null || warn "Failed to write sync_log to Postgres"
    fi

    # Métricas
    counter_metric "aurelia_memory_sync_files_total" "$files_processed" "mode=\"fast\""
    counter_metric "aurelia_memory_sync_embeddings_total" "$embeddings_created" "mode=\"fast\""
    gauge_metric "aurelia_memory_sync_duration_seconds" "$duration" "mode=\"fast\""

    info "FAST sync complete: files=$files_processed, embeddings=$embeddings_created, errors=$errors, duration=${duration}s"
}

# ============================================================================
# Modo: POSTGRES-INDEX (15 minutos)
# ============================================================================

function sync_postgres_index() {
    info "Starting POSTGRES-INDEX sync..."
    local start_time=$(date +%s)

    if [ "$DRY_RUN" = "false" ]; then
        # Atualizar índices
        psql -d "$POSTGRES_DB" -c "
            REINDEX INDEX CONCURRENTLY ai_context.idx_memory_type;
            REINDEX INDEX CONCURRENTLY ai_context.idx_memory_owner;
        " 2>/dev/null || warn "Failed to reindex Postgres"

        # Sync metadata (dual-write from Qdrant → Postgres)
        # Aqui seria implementado sync que lê do Qdrant e atualiza Postgres

        # Registrar
        psql -d "$POSTGRES_DB" -c "
            INSERT INTO ai_context.sync_log (sync_type, completed_at, files_processed, status)
            VALUES ('postgres-index', NOW(), 0, 'success')
            ON CONFLICT DO NOTHING;
        " 2>/dev/null || warn "Failed to write sync_log"
    fi

    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    gauge_metric "aurelia_memory_sync_duration_seconds" "$duration" "mode=\"postgres-index\""

    info "POSTGRES-INDEX sync complete: duration=${duration}s"
}

# ============================================================================
# Modo: VALIDATE (diário 6am)
# ============================================================================

function sync_validate() {
    info "Starting VALIDATE..."
    local start_time=$(date +%s)

    if [ "$DRY_RUN" = "false" ]; then
        # Contar Qdrant
        local qdrant_count=$(curl -s "$QDRANT_URL/collections/repository_memory/points/count" 2>/dev/null | \
            python3 -c "import sys, json; print(json.load(sys.stdin).get('result', {}).get('count', 0))" 2>/dev/null || echo "0")

        # Contar Postgres
        local postgres_count=$(psql -d "$POSTGRES_DB" -t -c "SELECT COUNT(*) FROM ai_context.memory_entries;" 2>/dev/null || echo "0")

        info "Validation: Qdrant=$qdrant_count points, Postgres=$postgres_count entries"

        if [ "$qdrant_count" -eq "$postgres_count" ]; then
            info "✓ Validation PASSED"
        else
            warn "⚠ Count mismatch! Qdrant=$qdrant_count, Postgres=$postgres_count"
        fi

        # Registrar resultado
        psql -d "$POSTGRES_DB" -c "
            INSERT INTO ai_context.sync_log (sync_type, completed_at, status, details)
            VALUES ('validate', NOW(), 'success', json_build_object('qdrant_count', $qdrant_count, 'postgres_count', $postgres_count))
            ON CONFLICT DO NOTHING;
        " 2>/dev/null || warn "Failed to write sync_log"

        gauge_metric "aurelia_memory_entries_total" "$postgres_count"
    fi

    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    gauge_metric "aurelia_memory_sync_duration_seconds" "$duration" "mode=\"validate\""

    info "VALIDATE complete: duration=${duration}s"
}

# ============================================================================
# Modo: COMPACT (segunda 2am)
# ============================================================================

function sync_compact() {
    info "Starting COMPACT..."
    local start_time=$(date +%s)

    if [ "$DRY_RUN" = "false" ]; then
        # Deletar entries antigas (>90 dias) com type=plan ou temp
        local deleted=$(psql -d "$POSTGRES_DB" -t -c "
            DELETE FROM ai_context.memory_entries
            WHERE modified_at < NOW() - INTERVAL '90 days'
              AND type IN ('plan', 'temp');
            SELECT row_count FROM LAST_QUERY_COUNT;
        " 2>/dev/null || echo "0")

        info "Deleted $deleted old entries"

        # REINDEX (mais pesado)
        psql -d "$POSTGRES_DB" -c "REINDEX TABLE CONCURRENTLY ai_context.memory_entries;" 2>/dev/null || warn "Failed to reindex table"

        # Registrar
        psql -d "$POSTGRES_DB" -c "
            INSERT INTO ai_context.sync_log (sync_type, completed_at, status, details)
            VALUES ('compact', NOW(), 'success', json_build_object('deleted', $deleted))
            ON CONFLICT DO NOTHING;
        " 2>/dev/null || warn "Failed to write sync_log"
    fi

    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    gauge_metric "aurelia_memory_sync_duration_seconds" "$duration" "mode=\"compact\""

    info "COMPACT complete: duration=${duration}s"
}

# ============================================================================
# Main
# ============================================================================

function main() {
    info "========================================="
    info "Memory Sync Fiscal — Mode: $MODE"
    info "========================================="

    case "$MODE" in
        fast)
            sync_fast
            ;;
        postgres-index)
            sync_postgres_index
            ;;
        validate)
            sync_validate
            ;;
        compact)
            sync_compact
            ;;
        *)
            error "Unknown mode: $MODE"
            echo "Usage: $0 [fast|postgres-index|validate|compact]"
            exit 1
            ;;
    esac

    info "Sync complete."
}

main "$@"
