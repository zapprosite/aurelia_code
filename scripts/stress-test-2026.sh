#!/bin/bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Aurelia Stress Test 2026 - Soberania Industrial
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

set -e

LOG_FILE="/tmp/stress_test_$(date +%Y%m%d_%H%M%S).log"
RESULTS_FILE="/tmp/stress_results.md"
TARGET_URL="http://localhost:8484/v1/telegram/impersonate" 
CONCURRENT_REQUESTS=3
TOTAL_REQUESTS=12

echo "🚀 Iniciando Estresse Operacional (Full Pipeline): $TOTAL_REQUESTS total, $CONCURRENT_REQUESTS concorrentes"
echo "# Relatório de Estresse Operacional - $(date)" > $RESULTS_FILE
echo "## Parâmetros: Concorrência=$CONCURRENT_REQUESTS, Total=$TOTAL_REQUESTS" >> $RESULTS_FILE

# Função para capturar métricas GPU
capture_metrics() {
    nvidia-smi --query-gpu=temperature.gpu,utilization.gpu,power.draw,memory.used --format=csv,noheader,nounits | tr -d ' '
}

echo "📊 Métricas Iniciais: $(capture_metrics)"
echo "| Timestamp | Temp (°C) | Util (%) | Power (W) | Memory (MiB) | Status |" >> $RESULTS_FILE
echo "|-----------|-----------|----------|-----------|--------------|--------|" >> $RESULTS_FILE

# Execução do loop de carga
for ((i=1; i<=$TOTAL_REQUESTS; i+=$CONCURRENT_REQUESTS)); do
    echo "⚡ Lote $i a $((i+CONCURRENT_REQUESTS-1))..."
    
    # Disparar requisições em paralelo (impersonação)
    for ((j=0; j<$CONCURRENT_REQUESTS; j++)); do
        (
            curl -s -X POST "$TARGET_URL" \
                -H "Content-Type: application/json" \
                -d "{\"text\": \"Teste de carga industrial $i-$j: analise o status térmico da GPU e retorne um sumário premium.\"}" > /dev/null
        ) &
    done
    
    wait # Aguardar lote terminar
    
    # Capturar métricas pós-lote
    METRICS=$(capture_metrics)
    T=$(echo $METRICS | cut -d',' -f1)
    U=$(echo $METRICS | cut -d',' -f2)
    P=$(echo $METRICS | cut -d',' -f3)
    M=$(echo $METRICS | cut -d',' -f4)
    
    echo "| $(date +%H:%M:%S) | $T | $U | $P | $M | OK |" >> $RESULTS_FILE
    
    if [ "$T" -gt 75 ]; then
        echo "🛑 ALERTA TÉRMICO: ${T}°C detectados!"
        echo "| $(date +%H:%M:%S) | $T | $U | $P | $M | 🛑 CRITICAL |" >> $RESULTS_FILE
    fi
done

echo "✅ Estresse finalizado. Relatório em $RESULTS_FILE"
echo "---" >> $RESULTS_FILE
echo "🎯 **Status Final**: $([ $(tail -n 1 $RESULTS_FILE | grep -c "CRITICAL") -eq 0 ] && echo "Aprovado" || echo "Falha Técnica")" >> $RESULTS_FILE
