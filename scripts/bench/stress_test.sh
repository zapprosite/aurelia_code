#!/bin/bash
# scripts/bench/stress_test.sh - SOTA 2026.1
# Teste de estresse mesclado: 25 TTS + 25 LLM (Gemma 3 27B)

CONCURRENCY=5
TOTAL_REQ=25
LOG_DIR="scripts/bench/logs"
mkdir -p $LOG_DIR

echo "--- INICIANDO STRESS TEST SOTA 2026.1 ---"
echo "Alvo LLM: gemma3:27b"
echo "Alvo TTS: localhost:8012 (Kodoro)"
echo "-----------------------------------------"

# Função para teste de TTS
stress_tts() {
    local id=$1
    local start=$(date +%s%N)
    curl -s -X POST "http://localhost:8012/tts" \
         -H "Content-Type: application/json" \
         -d '{"text": "Teste de estresse de voz Aurelia. Verificando resiliência do motor Kodoro sob carga industrial.", "voice": "af_heart"}' \
         -o /dev/null
    local end=$(date +%s%N)
    local diff=$(( (end - start) / 1000000 ))
    echo "TTS #$id: ${diff}ms" >> $LOG_DIR/results.log
}

# Função para teste de LLM
stress_llm() {
    local id=$1
    local start=$(date +%s%N)
    ollama run gemma3:27b "Responda de forma concisa: Qual a importância da soberania digital em 2026?" > /dev/null 2>&1
    local end=$(date +%s%N)
    local diff=$(( (end - start) / 1000000 ))
    echo "LLM #$id: ${diff}ms" >> $LOG_DIR/results.log
}

export -f stress_tts stress_llm
export LOG_DIR

echo "Executando 25 requisições de cada..."
for i in $(seq 1 $TOTAL_REQ); do
    stress_tts $i & 
    stress_llm $i &
    # Controle de concorrência simples
    if (( $i % $CONCURRENCY == 0 )); then
        wait
        nvidia-smi --query-gpu=memory.used,utilization.gpu,temperature.gpu --format=csv >> $LOG_DIR/gpu.log
    fi
done

wait
echo "--- TESTE CONCLUÍDO ---"
cat $LOG_DIR/results.log | sort
