#!/bin/bash
#
# Script: telegram-quality-stress.sh
# Descrição: Stress test + Quality Assessment do Telegram bot
# Versão flexível - foca em verificar se o bot responde e se a resposta é relevante
# Uso: ./telegram-quality-stress.sh [test_count]
#

set -e

COUNT=${1:-5}
BOT_TOKEN="TELEGRAM_BOT_TOKEN_PLACEHOLDER"
CHAT_ID="7220607041"

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_ok() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_err() { echo -e "${RED}[ERR]${NC} $1"; }

echo "=========================================="
echo "  TELEGRAM BOT QUALITY ASSESSMENT"
echo "  Dev Senior Pro Standards (Flex)"
echo "=========================================="
echo "Test count: $COUNT"
echo ""

# Test prompts - variety of complexity
declare -a TEST_PROMPTS=(
    "/status"
    "Qual o status do sistema?"
    "Me mostre a GPU"
    "Lista containers"
    "Temp CPU?"
    "RAM livre?"
    "Ollama online?"
    "LiteLLM status"
    "Ver métricas"
    "Dashboard Grafana"
)

# Wait for bot response via getUpdates
wait_for_bot_response() {
    local timeout=$1
    local attempt=0
    
    while [ $attempt -lt $timeout ]; do
        updates=$(curl -s "https://api.telegram.org/bot$BOT_TOKEN/getUpdates?timeout=1" 2>/dev/null)
        count=$(echo "$updates" | jq '.result | length')
        
        if [ "$count" -gt 0 ]; then
            # Get last message
            last_msg=$(echo "$updates" | jq -r '.result[-1].message.text // empty')
            last_from=$(echo "$updates" | jq -r '.result[-1].message.from.id // 0')
            
            # Get bot ID
            bot_id=$(curl -s "https://api.telegram.org/bot$BOT_TOKEN/getMe" | jq -r '.result.id')
            
            # Check if message is from bot
            if [ "$last_from" = "$bot_id" ] && [ -n "$last_msg" ] && [ "$last_msg" != "null" ]; then
                echo "$last_msg"
                return 0
            fi
        fi
        
        sleep 1
        attempt=$((attempt + 1))
    done
    
    echo ""
    return 1
}

# Quality check - flexible criteria
check_response_quality() {
    local response="$1"
    local score=0
    local max_score=100
    
    # 1. Response not empty (30 pts)
    if [ -n "$response" ] && [ "$response" != "null" ]; then
        score=$((score + 30))
    else
        echo "❌ Resposta vazia"
        return 1
    fi
    
    # 2. Contains relevant info (30 pts)
    if echo "$response" | grep -qiE "cpu|gpu|ram|container|ollama|litellm|qdrant|status|online|offline|ok|°c|%|w|gb|mb|tb|gi"; then
        score=$((score + 30))
    fi
    
    # 3. Has formatting (20 pts)
    if echo "$response" | grep -qE '\*\*|\`\`\`|\||\[\]'; then
        score=$((score + 20))
    fi
    
    # 4. Has structure/tables (10 pts)
    if echo "$response" | grep -qE '\|.*\|'; then
        score=$((score + 10))
    fi
    
    # 5. Has icons/emojis (10 pts)
    if echo "$response" | grep -qE '🟢|🔴|🟡|⚠️|✅|❌|🎛️|🖥️|📊|💾|🐳|🔥|💽|📈'; then
        score=$((score + 10))
    fi
    
    echo "Score: $score/100"
    [ $score -ge 50 ]
}

# Single test
run_test() {
    local msg="$1"
    local test_num=$2
    
    log_info "Test $test_num: $msg"
    
    # Get update offset before sending
    offset_before=$(curl -s "https://api.telegram.org/bot$BOT_TOKEN/getUpdates?limit=1" | jq -r '.result[-1].update_id // 0')
    
    # Send message
    send_resp=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/sendMessage" \
        -H "Content-Type: application/json" \
        -d "{\"chat_id\": $CHAT_ID, \"text\": \"$msg\"}")
    
    msg_id=$(echo "$send_resp" | jq -r '.result.message_id // empty')
    
    if [ -z "$msg_id" ] || [ "$msg_id" = "null" ]; then
        log_err "Falha ao enviar"
        return 1
    fi
    
    log_ok "Enviado (ID: $msg_id)"
    log_info "Aguardando resposta..."
    
    # Wait for response
    response=$(wait_for_bot_response 20) || true
    
    if [ -z "$response" ]; then
        log_warn "Timeout - sem resposta em 20s"
        return 1
    fi
    
    echo ""
    echo "┌─────────────────────────────────────────"
    echo "│ 📥 RESPOSTA:"
    echo "├─────────────────────────────────────────"
    echo "$response" | head -20 | sed 's/^/│ /'
    echo "└─────────────────────────────────────────"
    echo ""
    
    # Quality check
    if check_response_quality "$response"; then
        log_ok "Qualidade OK"
        return 0
    else
        log_warn "Qualidade abaixo do esperado"
        return 1
    fi
}

echo "=========================================="
echo "  FASE 1: SMOKE TEST"
echo "=========================================="

# Bot info
BOT_NAME=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/getMe" | jq -r '.result.first_name')
log_ok "Bot: $BOT_NAME"

# Chat
CHAT_NAME=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/getChat" \
    -H "Content-Type: application/json" \
    -d "{\"chat_id\": $CHAT_ID}" | jq -r '.result.first_name // "unknown"')
log_ok "Chat: $CHAT_NAME"

echo ""
echo "=========================================="
echo "  FASE 2: QUALITY TESTS"
echo "=========================================="

PASSED=0
FAILED=0

for i in $(seq 1 $COUNT); do
    idx=$(( (i-1) % ${#TEST_PROMPTS[@]} ))
    prompt="${TEST_PROMPTS[$idx]}"
    
    echo ""
    echo "═══════════════════════════════════════"
    
    if run_test "$prompt" "$i"; then
        PASSED=$((PASSED + 1))
    else
        FAILED=$((FAILED + 1))
    fi
    
    sleep 2
done

echo ""
echo "=========================================="
echo "  📊 RESULTADOS"
echo "=========================================="
echo "Tests: $COUNT | ✅ $PASSED | ❌ $FAILED"

pass_rate=$((PASSED * 100 / COUNT))
if [ $pass_rate -eq 100 ]; then
    echo -e "${GREEN}✅ QUALIDADE DEV SENIOR PRO${NC}"
elif [ $pass_rate -ge 60 ]; then
    echo -e "${YELLOW}⚠️ QUALIDADE BOA${NC}"
else
    echo -e "${RED}❌ PRECISA REFINAR${NC}"
fi
echo ""
