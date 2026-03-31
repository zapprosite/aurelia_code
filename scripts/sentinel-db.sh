#!/bin/bash
# 🛠️ sentinel-db.sh (SOTA 2026) — Especialista Supabase
# Autoridade: Aurélia Sentinel Swarm
# Versão: 1.2.0 (Industrializado)

set -e

WORKSPACE_DIR="/home/will/aurelia"
LOG_FILE="$WORKSPACE_DIR/logs/sentinel-db-$(date +%Y%m%d).log"
SENTINEL_STATE_FILE="$WORKSPACE_DIR/data/.sentinel_state"

# Função para logging estruturado (JSON-like para fácil parse por agentes)
log_event() {
    local level=$1
    local scope=$2
    local message=$3
    echo "[$(date -Iseconds)] [$level] [$scope] $message" | tee -a "$LOG_FILE"
}

log_event "INFO" "SENTINEL-DB" "Iniciando Auditoria Industrial do Banco de Dados..."

# 1. Verificar Outliers de Consultas (Performance)
log_event "DEBUG" "QUERY_ANALYSER" "Analisando Slow Queries via Supabase CLI..."
if supabase db outliers -n 5 >> "$LOG_FILE" 2>&1; then
    log_event "INFO" "QUERY_ANALYSER" "Analise de performance concluída."
else
    log_event "WARN" "QUERY_ANALYSER" "Falha ao obter outliers. Supabase CLI configurado?"
fi

# 2. Integridade de Migrações
log_event "DEBUG" "MIGRATION_GUARD" "Verificando paridade de migrações..."
OUT_MIGRATION=$(supabase migration list 2>/dev/null || echo "error")
if [[ "$OUT_MIGRATION" == *"false"* ]] || [[ "$OUT_MIGRATION" == "error" ]]; then
    log_event "CRITICAL" "MIGRATION_GUARD" "Discrepância detectada em migrações! Bloqueando deploy."
    # Registro de estado para o Orquestrador
    echo "STATUS=ERROR SCOPE=MIGRATIONS" > "$SENTINEL_STATE_FILE"
else
    log_event "OK" "MIGRATION_GUARD" "Migrações em sync."
fi

# 3. Resiliência do Serviço REST
log_event "DEBUG" "HEALTH_CHECK" "Validando API REST Supabase..."
if curl -s -f "http://localhost:54321/rest/v1/" > /dev/null; then
    log_event "OK" "HEALTH_CHECK" "API REST Online."
else
    log_event "ERROR" "HEALTH_CHECK" "API REST Offline! Disparando Hard Restart..."
    docker restart supabase_rest
fi

log_event "INFO" "SENTINEL-DB" "Auditoria Finalizada com sucesso."
