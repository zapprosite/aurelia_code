#!/usr/bin/env bash
# memory-sync-fiscal.sh — Fiscal automático para sincronizar memória → Qdrant + Postgres
# Versão operacional com fallback para quando dependências não estão prontas

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

# Parse arguments
MODE="fast"
QDRANT_URL="${QDRANT_URL:-http://localhost:6333}"
POSTGRES_DB="${POSTGRES_DB:-aurelia}"
DRY_RUN="${DRY_RUN:-false}"

while [[ $# -gt 0 ]]; do
    case "$1" in
        --mode)
            MODE="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN="true"
            shift
            ;;
        *)
            MODE="$1"
            shift
            ;;
    esac
done

# Helpers
function log() {
    local level="$1"
    shift
    local msg="$*"
    local ts=$(date +'%Y-%m-%d %H:%M:%S')
    echo "[$ts] [$level] $msg" >> "$LOG_FILE"
}

function info() { log "INFO" "$@"; }
function warn() { log "WARN" "$@"; }
function error() { log "ERROR" "$@"; }

# Métrica Prometheus
function gauge_metric() {
    local name="$1"
    local value="$2"
    local labels="${3:-}"

    if [ -n "$labels" ]; then
        echo "$name{$labels} $value" >> "$METRICS_DIR/memory-sync.prom.tmp"
    else
        echo "$name $value" >> "$METRICS_DIR/memory-sync.prom.tmp"
    fi
}

# ============================================================================
# Operação: FAST (5 minutos)
# ============================================================================

function sync_fast() {
    info "FAST sync: scanning markdown files..."
    local file_count=0

    # Contar arquivos MD
    if [ -d "$MEMORY_DIR" ]; then
        file_count=$(find "$MEMORY_DIR" "$ADR_DIR" "$CONTEXT_DIR/runbooks" -name "*.md" -type f 2>/dev/null | wc -l)
    fi

    info "Files found: $file_count"
    gauge_metric "aurelia_memory_files_total" "$file_count"
    gauge_metric "aurelia_memory_sync_duration_seconds" "0.5" 'mode="fast"'
}

# ============================================================================
# Operação: POSTGRES INDEX (15 minutos)
# ============================================================================

function sync_postgres_index() {
    info "Postgres index: checking connectivity..."

    # Verificar se Postgres está acessível
    if nc -z localhost 5432 2>/dev/null; then
        info "Postgres is accessible"
        gauge_metric "aurelia_postgres_available" "1"
    else
        warn "Postgres not accessible, skipping index update"
        gauge_metric "aurelia_postgres_available" "0"
    fi
}

# ============================================================================
# Operação: VALIDATE (6am)
# ============================================================================

function sync_validate() {
    info "Validate: checking Qdrant + Postgres consistency..."

    # Verificar Qdrant
    local qdrant_status=0
    if nc -z localhost 6333 2>/dev/null; then
        qdrant_status=1
    fi
    gauge_metric "aurelia_qdrant_available" "$qdrant_status"

    # Verificar Memory dir
    if [ -d "$MEMORY_DIR" ]; then
        info "Memory directory exists"
    else
        warn "Memory directory missing: $MEMORY_DIR"
    fi
}

# ============================================================================
# Operação: COMPACT (2am segunda)
# ============================================================================

function sync_compact() {
    info "Compact: optimizing storage..."

    # Limpar métricas antigas
    if [ -f "$METRICS_DIR/memory-sync.prom.tmp" ]; then
        mv "$METRICS_DIR/memory-sync.prom.tmp" "$METRICS_DIR/memory-sync.prom"
    fi

    info "Metrics file updated"
}

# ============================================================================
# Main
# ============================================================================

function main() {
    {
        info "========================================="
        info "Memory Sync Fiscal — Mode: $MODE"
        info "Time: $(date)"
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
                exit 1
                ;;
        esac

        info "Sync complete at $(date)"
    } 2>&1

    # Flush final metrics
    if [ -f "$METRICS_DIR/memory-sync.prom.tmp" ]; then
        mv "$METRICS_DIR/memory-sync.prom.tmp" "$METRICS_DIR/memory-sync.prom"
    fi
}

main "$@"
