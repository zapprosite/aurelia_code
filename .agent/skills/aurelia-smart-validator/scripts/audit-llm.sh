#!/bin/bash
# Auditoria Soberana de Inferência (v2026.1) - REST Native
set -e

# Carregar LITELLM_MASTER_KEY se o arquivo .env existir
if [ -f .env ]; then
    export $(grep LITELLM_MASTER_KEY .env | xargs)
fi

LITELLM_URL=${LITELLM_LOCAL_URL:-"http://localhost:4000"}
API_KEY=${LITELLM_MASTER_KEY:-"sk-1234"}

echo "=== 🛡️ AURELIA SMART AUDIT: SOTA 2026.1 ==="
echo "Endpoint: $LITELLM_URL"

check_model() {
    local prompt=$1
    local desc=$2
    echo -n "--- Teste: $desc ... "
    
    RESPONSE=$(curl -s -X POST "$LITELLM_URL/chat/completions" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_KEY" \
        -d "{
            \"model\": \"aurelia-smart\",
            \"messages\": [{\"role\": \"user\", \"content\": \"$prompt\"}],
            \"max_tokens\": 50
        }")

    if echo "$RESPONSE" | grep -q "choices"; then
        MODEL=$(echo "$RESPONSE" | grep -oP '"model":"\K[^"]+')
        echo "✅ OK ($MODEL)"
    else
        echo "❌ FALHA"
        echo "Detalhes: $RESPONSE"
    fi
}

# --- T0: Soberania Local ---
check_model "Responda apenas LOCAL." "Tier 0 (Local)"

# --- T1/T2: Cascata ---
check_model "Diga o nome do maior planeta do sistema solar." "Tier 1/2 (Cascade)"

# --- EXECUTE EMBEDDINGS ---
echo -n "--- Teste: Embeddings (Nomic) ... "
E_RESP=$(curl -s -X POST "$LITELLM_URL/embeddings" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $API_KEY" \
    -d "{
        \"model\": \"openai/nomic-embed-text\",
        \"input\": [\"Aurelia Smart SOTA 2026\"]
    }")

if echo "$E_RESP" | grep -q "embedding"; then
    echo "✅ OK"
else
    echo "❌ FALHA"
    echo "$E_RESP"
fi

echo "=== ✅ Auditoria Finalizada ==="
