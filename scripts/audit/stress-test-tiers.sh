#!/bin/bash
# 🧪 STRESS TEST: AURELIA SMART TIERS (v2026.1)
# Valida a transição local -> cascade -> expert

LITELLM_URL="http://localhost:4000/chat/completions"
API_KEY="sk-1234"

function run_test() {
    local label=$1
    local prompt=$2
    echo -e "\n🚀 [TESTE] $label"
    
    START=$(date +%s.%N)
    RESPONSE=$(curl -s -X POST "$LITELLM_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_KEY" \
        -d "{
            \"model\": \"aurelia-smart\",
            \"messages\": [{\"role\": \"user\", \"content\": \"$prompt\"}],
            \"max_tokens\": 100
        }")
    END=$(date +%s.%N)
    
    DURATION=$(echo "$END - $START" | bc)
    MODEL=$(echo "$RESPONSE" | grep -oP '"model":"\K[^"]+')
    CONTENT=$(echo "$RESPONSE" | grep -oP '"content":"\K[^"]+' | head -c 100)

    echo "--- Modelo: $MODEL"
    echo "--- Tempo:  ${DURATION}s"
    echo "--- Resposta: $CONTENT..."
}

echo "=== 🛡️ INICIANDO STRESS TEST TIERS SOTA 2026.1 ==="

# 1. Teste de Soberania (Tier 0)
# Deve responder rápido via Gemma 3 Local
run_test "Tier 0: Soberania Local" "Diga apenas 'Soberania Local Ativa'."

# 2. Teste de Velocidade (Groq Turbo)
# Vamos pedir algo que gere tokens para testar a velocidade (Groq Scout)
run_test "Tier 1: Groq Turbo Speed" "Escreva um poema curto sobre a velocidade da luz e Groq LPU."

# 3. Teste de Complexidade/Failover (Forçar Transição)
# Nota: Para testar T1/T2, podemos forçar o roteador a pular o T0 
# enviando um prompt muito longo ou simulando carga.
# Aqui vamos validar o retorno do modelo para ver onde ele caiu.
run_test "Tier 1/2: Cascade/Expert" "Analise a estrutura de um banco de dados vetorial Qdrant e compare com Postgres pgvector."

echo -e "\n=== ✅ TESTE DE ESTRESSE FINALIZADO ==="
