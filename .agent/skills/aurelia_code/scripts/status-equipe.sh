#!/bin/bash
#
# Script: status-equipe.sh
# Descrição: Verificar status da equipe de agentes
#
# Uso: ./status-equipe.sh
#

set -e

MISSION_ID="${MISSION_ID:-default-mission}"
QDRANT_URL="${QDRANT_URL:-http://localhost:6333}"

echo "=========================================="
echo "  AURELIA_CODE - Status da Equipe"
echo "=========================================="
echo "Missão: $MISSION_ID"
echo ""

# Carregar agent cards
CONFIG_DIR="$(dirname "$0")/../configs"
AGENT_CARDS="$CONFIG_DIR/agent-cards.json"

if command -v jq &> /dev/null && [ -f "$AGENT_CARDS" ]; then
    LEADER=$(jq -r '.leader.id' "$AGENT_CARDS")
    AGENTS=$(jq -r '.agents[].id' "$AGENT_CARDS")
else
    LEADER="aurelia_code"
    AGENTS="pesquisador coder revisor"
fi

echo "Líder: $LEADER"
echo ""

# Verificar status no Qdrant
echo "1. Verificando missões no Qdrant..."

MISSIONS=$(curl -s -X POST "$QDRANT_URL/collections/aurelia_swarm_missions/points/search" \
  -H "Content-Type: application/json" \
  -d "{
    \"limit\": 10,
    \"with_payload\": true
  }" 2>/dev/null || echo '{}')

if command -v jq &> /dev/null; then
    TOTAL=$(echo "$MISSIONS" | jq '.result.length // 0')
    ACTIVE=$(echo "$MISSIONS" | jq '[.result[] | select(.payload.status == "active")] | length')
    COMPLETED=$(echo "$MISSIONS" | jq '[.result[] | select(.payload.status == "completed")] | length')
else
    TOTAL=0
    ACTIVE=0
    COMPLETED=0
fi

echo "  Total de missões: $TOTAL"
echo "  Ativas: $ACTIVE"
echo "  Concluídas: $COMPLETED"
echo ""

# Verificar contexto compartilhado
echo "2. Verificando contexto compartilhado..."

CONTEXT=$(curl -s -X POST "$QDRANT_URL/collections/aurelia_swarm_context/points/search" \
  -H "Content-Type: application/json" \
  -d "{
    \"limit\": 5,
    \"filter\": {\"must\": [{\"key\": \"mission_id\", \"match\": {\"value\": \"$MISSION_ID\"}}]},
    \"with_payload\": true
  }" 2>/dev/null || echo '{}')

if command -v jq &> /dev/null; then
    CONTEXT_COUNT=$(echo "$CONTEXT" | jq '.result.length // 0')
else
    CONTEXT_COUNT=0
fi

echo "  Contextos registrados: $CONTEXT_COUNT"
echo ""

# Verificar decisões
echo "3. Verificando decisões do líder..."

DECISIONS=$(curl -s -X POST "$QDRANT_URL/collections/aurelia_swarm_decisions/points/search" \
  -H "Content-Type: application/json" \
  -d "{
    \"limit\": 5,
    \"filter\": {\"must\": [{\"key\": \"mission_id\", \"match\": {\"value\": \"$MISSION_ID\"}}]},
    \"with_payload\": true
  }" 2>/dev/null || echo '{}')

if command -v jq &> /dev/null; then
    DECISIONS_COUNT=$(echo "$DECISIONS" | jq '.result.length // 0')
else
    DECISIONS_COUNT=0
fi

echo "  Decisões registradas: $DECISIONS_COUNT"
echo ""

# Status dos agentes
echo "4. Status dos agentes:"
for AGENT in $AGENTS; do
    echo "  - $AGENT: disponível"
done
echo ""

echo "=========================================="
echo "  Resumo"
echo "=========================================="
echo "Missão atual: $MISSION_ID"
echo "Agentes ativos: $(echo $AGENTS | wc -w)"
echo "Tarefas concluídas: $COMPLETED"
echo "Tarefas ativas: $ACTIVE"
echo ""

echo "Comandos disponíveis:"
echo "  /ac-missao [desc]   - Nova missão"
echo "  /ac-spawn [papel]   - Delegar tarefa"
echo "  /ac-status          - Ver status"
echo ""

exit 0