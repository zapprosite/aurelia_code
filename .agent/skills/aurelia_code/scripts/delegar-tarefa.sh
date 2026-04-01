#!/bin/bash
#
# Script: delegar-tarefa.sh
# Descrição: Delegar tarefa para sub-agent via A2A
#
# Uso: ./delegar-tarefa.sh [papel] [tarefa]
#

set -e

PAPEL="${1:-}"
TAREFA="${2:-}"

if [ -z "$PAPEL" ] || [ -z "$TAREFA" ]; then
    echo "Erro: Papel e tarefa são obrigatórios"
    echo "Uso: $0 [papel] [tarefa]"
    echo ""
    echo "Papeis disponíveis:"
    echo "  Pesquisador  - Pesquisa web e docs"
    echo "  Coder        - Implementação de código"
    echo "  Revisor      - Code review e testes"
    exit 1
fi

# Mapear papel para agent ID
case "$PAPEL" in
    Pesquisador|pesquisador)
        AGENT_ID="pesquisador"
        AGENT_URL="http://localhost:8081/agent"
        ;;
    Coder|coder)
        AGENT_ID="coder"
        AGENT_URL="http://localhost:8082/agent"
        ;;
    Revisor|revisor)
        AGENT_ID="revisor"
        AGENT_URL="http://localhost:8083/agent"
        ;;
    *)
        echo "Erro: Papel desconhecido '$PAPEL'"
        exit 1
        ;;
esac

TASK_ID=$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid)
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
MISSION_ID="${MISSION_ID:-default-mission}"

echo "=========================================="
echo "  Delegando para $AGENT_ID"
echo "=========================================="
echo "Tarefa: $TAREFA"
echo "Task ID: $TASK_ID"
echo ""

# Enviar task via A2A (JSON-RPC 2.0)
echo "1. Enviando task via A2A..."

A2A_REQUEST=$(cat <<EOF
{
  "jsonrpc": "2.0",
  "id": "$TASK_ID",
  "method": "tasks/send",
  "params": {
    "task": {
      "id": "$TASK_ID",
      "sessionId": "$MISSION_ID",
      "message": {
        "role": "user",
        "parts": [{"type": "text", "text": "$TAREFA"}]
      }
    },
    "agentId": "$AGENT_ID"
  }
}
EOF
)

# Tentar enviar (fallback se não disponível)
if command -v curl &> /dev/null; then
    RESPONSE=$(curl -s -X POST "$AGENT_URL/rpc" \
        -H "Content-Type: application/json" \
        -d "$A2A_REQUEST" 2>/dev/null) || RESPONSE='{"error": {"message": "Agent não disponível"}}'
else
    RESPONSE='{"error": {"message": "curl não disponível"}}'
fi

echo "Resposta: $RESPONSE"
echo ""

# Registrar no Qdrant
QDRANT_URL="${QDRANT_URL:-http://localhost:6333}"
echo "2. Registrando delegação..."

curl -s -X POST "$QDRANT_URL/collections/aurelia_swarm_missions/points/upsert" \
  -H "Content-Type: application/json" \
  -d "{
    \"points\": [{
      \"id\": \"$TASK_ID\",
      \"vector\": [],
      \"payload\": {
        \"task_id\": \"$TASK_ID\",
        \"mission_id\": \"$MISSION_ID\",
        \"role\": \"$AGENT_ID\",
        \"description\": \"$TAREFA\",
        \"status\": \"delegated\",
        \"delegated_at\": \"$TIMESTAMP\"
      }
    }]
  }" 2>/dev/null || echo "Aviso: Qdrant não disponível"

echo ""
echo "✓ Task delegada para $AGENT_ID"
echo ""
echo "Para verificar status: /ac-status"
echo ""

exit 0