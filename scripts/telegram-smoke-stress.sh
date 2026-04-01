#!/bin/bash
#
# Script: telegram-smoke-stress.sh
# Descrição: Smoke test + Stress test do Telegram bot
# Uso: ./telegram-smoke-stress.sh [count] [delay_ms]
#

set -e

COUNT=${1:-50}
DELAY_MS=${2:-500}
BOT_TOKEN="8793928549:AAEr3tjaarijUWxu-iru0Vcm6N6DkwjndL4"
CHAT_ID="7220607041"

echo "=========================================="
echo "  TELEGRAM SMOKE + STRESS TEST"
echo "=========================================="
echo "Messages: $COUNT"
echo "Delay: ${DELAY_MS}ms"
echo ""

# Função para enviar mensagem e medir tempo
send_and_measure() {
    local msg="$1"
    local start=$(date +%s%N)
    
    # Usar getUpdates para ver se o bot está respondendo
    # Primeiro envia uma mensagem ao bot (via usuário)
    # Luego verifica as respostas
    
    echo "   → Enviando: $msg"
}

# Test 1: Verificar se bot está online
echo "1. Smoke test - getMe"
RESPONSE=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/getMe" | jq -r '.result.first_name // .error_description')
echo "   Bot: $RESPONSE"

# Test 2: Verificar status do bot
echo ""
echo "2. Smoke test - getWebhookInfo"
RESPONSE=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/getWebhookInfo" | jq '.result')
echo "   Webhook: ${RESPONSE:0:200}..."

# Test 3: Contagem de mensagens no chat
echo ""
echo "3. Smoke test - getChat"
RESPONSE=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/getChat" \
  -H "Content-Type: application/json" \
  -d "{\"chat_id\": $CHAT_ID}" | jq -r '.result.first_name // .error_description')
echo "   Chat: $RESPONSE"

# Stress: Loop de mensagens de teste (enviadas via API, não são respondidas pelo bot)
echo ""
echo "4. Stress test - $COUNT mensagens (API throughput)"
echo ""

START_TIME=$(date +%s)
SUCCESS=0
FAIL=0
TIMES=()

for i in $(seq 1 $COUNT); do
    MSG="Stress test $i/$COUNT"
    
    MSG_START=$(date +%s%N)
    RESPONSE=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/sendMessage" \
      -H "Content-Type: application/json" \
      -d "{\"chat_id\": $CHAT_ID, \"text\": \"$MSG\"}" 2>/dev/null)
    MSG_END=$(date +%s%N)
    
    ELAPSED=$(( (MSG_END - MSG_START) / 1000000 ))
    TIMES+=($ELAPSED)
    
    if echo "$RESPONSE" | jq -e '.ok' > /dev/null 2>&1; then
        SUCCESS=$((SUCCESS + 1))
        if [ $((i % 10)) -eq 0 ]; then
            echo "   [$i/$COUNT] OK - ${ELAPSED}ms"
        fi
    else
        FAIL=$((FAIL + 1))
        echo "   [$i/$COUNT] FAIL: $(echo "$RESPONSE" | jq -r '.error_description // "unknown"')"
    fi
    
    sleep "0.${DELAY_MS}"
done

END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))

echo ""
echo "=========================================="
echo "  RESULTADOS (API Throughput)"
echo "=========================================="
echo "Total time: ${TOTAL_TIME}s"
echo "Success: $SUCCESS/$COUNT"
echo "Fail: $FAIL/$COUNT"

# Calcular média
SUM=0
for t in "${TIMES[@]}"; do
    SUM=$((SUM + t))
done
AVG=$((SUM / COUNT))
echo "Avg response: ${AVG}ms"

# Calcular p95
sorted=($(printf '%s\n' "${TIMES[@]}" | sort -n))
p95_idx=$((COUNT * 95 / 100))
p95=${sorted[$p95_idx]}
echo "P95 response: ${p95}ms"

# Métricas
echo ""
echo "Métricas:"
if [ $FAIL -eq 0 ]; then
    echo "  ✅ 100% API success rate"
else
    FAIL_RATE=$((FAIL * 100 / COUNT))
    echo "  ⚠️  $FAIL_RATE% API failure rate"
fi

if [ $AVG -lt 1000 ]; then
    echo "  ✅ Avg < 1s (fast)"
else
    echo "  ⚠️  Avg > 1s (slow)"
fi

if [ $p95 -lt 2000 ]; then
    echo "  ✅ P95 < 2s"
else
    echo "  ⚠️  P95 > 2s (slow)"
fi

echo ""
echo "=========================================="
echo "  NOTA"
echo "=========================================="
echo "Este teste mede API throughput, não latência real do bot."
echo "O bot responde apenas a mensagens recebidas do usuário."
echo "Para testar o bot: envie mensagens no Telegram Desktop."
echo ""
echo "=========================================="
echo "  FIM DO TESTE"
echo "=========================================="
