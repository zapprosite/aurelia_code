#!/bin/bash

# 🛡️ Latency Check Script (SOTA 2026.1)
# Monitora Ollama, LiteLLM e Endpoints Cloud.

source /home/will/aurelia/.env

CHECK_URLS=(
    "Ollama (Local):http://127.0.0.1:11434/api/tags"
    "LiteLLM (Smart Router):http://127.0.0.1:4000/health"
    "Qdrant (Vector DB):http://127.0.0.1:6333/healthz"
)

echo "📊 --- Aurelia Latency Audit ---"
for entry in "${CHECK_URLS[@]}"; do
    NAME="${entry%%:*}"
    URL="${entry#*:}"
    
    start_time=$(date +%s%N)
    status=$(curl -o /dev/null -s -w "%{http_code}" "$URL")
    end_time=$(date +%s%N)
    
    latency=$(( (end_time - start_time) / 1000000 ))
    
    if [ "$status" -eq 200 ]; then
        echo "✅ $NAME: ${latency}ms"
    else
        echo "❌ $NAME: FAILED (Status $status)"
    fi
done

# Check GPU Load (Contention indicator)
GPU_LOAD=$(nvidia-smi --query-gpu=utilization.gpu --format=csv,noheader,nounits)
echo "🔥 GPU Load: ${GPU_LOAD}%"
