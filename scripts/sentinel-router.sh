#!/bin/bash
# 🛠️ sentinel-router.sh (SOTA 2026) — Especialista LiteLLM
# Autoridade: Aurélia Sentinel Swarm

LOG_FILE="/home/will/aurelia/logs/sentinel.log"

echo "[$(date)] [SENTINEL-ROUTER] Iniciando Auditoria de Roteamento..." >> "$LOG_FILE"

# 1. Verificar Status da API LiteLLM
if curl -s -X GET "http://localhost:4000/health" > /dev/null; then
    echo "✅ LiteLLM Router Online." >> "$LOG_FILE"
else
    echo "❌ LiteLLM Router Offline! Reiniciando container..." >> "$LOG_FILE"
    docker restart litellm
fi

# 2. Teste de Latência de Endpoints (Ping)
# (Aqui poderíamos iterar sobre as chaves e testar cada uma)
echo "🔍 Testando latência de fallback..." >> "$LOG_FILE"
LATENCY=$(curl -o /dev/null -s -w "%{time_total}\n" http://localhost:4000/v1/models)
echo "📊 Latência Global de Resposta: ${LATENCY}s" >> "$LOG_FILE"

echo "[$(date)] [SENTINEL-ROUTER] Auditoria Concluída." >> "$LOG_FILE"
