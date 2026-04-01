#!/bin/bash
#
# Script: iniciar-missao.sh
# Descrição: Criar nova missão e iniciar squad de agentes
#
# Uso: ./iniciar-missao.sh "Descrição da missão"
#

set -e

MISSION_DESCRIPTION="${1:-}"
if [ -z "$MISSION_DESCRIPTION" ]; then
    echo "Erro: Descrição da missão é obrigatória"
    echo "Uso: $0 \"Descrição da missão\""
    exit 1
fi

MISSION_ID=$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid)
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "=========================================="
echo "  AURELIA_CODE - Nova Missão"
echo "=========================================="
echo "Missão: $MISSION_DESCRIPTION"
echo "ID: $MISSION_ID"
echo "Timestamp: $TIMESTAMP"
echo ""

# Carregar configs
CONFIG_DIR="$(dirname "$0")/../configs"
AGENT_CARDS="$CONFIG_DIR/agent-cards.json"

if [ ! -f "$AGENT_CARDS" ]; then
    echo "Erro: Agent cards não encontrado em $AGENT_CARDS"
    exit 1
fi

# Extrair agentes do JSON (usando jq se disponível, ou fallback)
if command -v jq &> /dev/null; then
    AGENTS=$(jq -r '.agents[].id' "$AGENT_CARDS")
    LEADER=$(jq -r '.leader.id' "$AGENT_CARDS")
else
    echo "Aviso: jq não instalado. Usando fallback manual."
    AGENTS="pesquisador coder revisor"
    LEADER="aurelia_code"
fi

echo "Líder: $LEADER"
echo "Agentes disponíveis: $AGENTS"
echo ""

# Criar registro inicial da missão no Qdrant
QDRANT_URL="${QDRANT_URL:-http://localhost:6333}"

echo "1. Registrando missão no Qdrant..."
curl -s -X POST "$QDRANT_URL/collections/aurelia_swarm_missions/points" \
  -H "Content-Type: application/json" \
  -d "{
    \"points\": [{
      \"id\": \"$MISSION_ID\",
      \"vector\": [],
      \"payload\": {
        \"mission_id\": \"$MISSION_ID\",
        \"description\": \"$MISSION_DESCRIPTION\",
        \"leader\": \"$LEADER\",
        \"status\": \"active\",
        \"created_at\": \"$TIMESTAMP\",
        \"tasks\": []
      }
    }]
  }" || echo "Aviso: Qdrant não disponível, continuando..."

echo ""
echo "2. Squad pronto para missão!"
echo ""
echo "Para delegar tarefas, use:"
echo "  /ac-spawn Pesquisador \"pesquisar X\""
echo "  /ac-spawn Coder \"implementar Y\""
echo "  /ac-spawn Revisor \"revisar Z\""
echo ""
echo "Para ver status: /ac-status"
echo ""

exit 0