#!/bin/bash
echo "📡 Aurelia Home Lab Health-Check"
echo "--------------------------------"

# 1. NVIDIA
if command -v nvidia-smi &> /dev/null; then
    echo "✅ GPU: $(nvidia-smi --query-gpu=gpu_name,memory.used,memory.total --format=csv,noheader,nounits)"
else
    echo "❌ GPU: NVIDIA SMI não encontrado."
fi

# 2. Ollama
OLLAMA_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:11434/api/tags)
if [ "$OLLAMA_STATUS" == "200" ]; then
    echo "✅ Ollama: Online ($(curl -s http://localhost:11434/api/tags | jq -r '.models | length') modelos)"
else
    echo "❌ Ollama: Offline ou inacessível."
fi

# 3. Docker
if command -v docker &> /dev/null; then
    RUNNING_CONTAINERS=$(docker ps -q | wc -l)
    echo "✅ Docker: Online ($RUNNING_CONTAINERS containers rodando)"
else
    echo "❌ Docker: não encontrado."
fi

# 4. ZFS (se disponível)
if command -v zpool &> /dev/null; then
    echo "✅ ZFS: $(zpool list -H -o name,health || echo 'Nenhum pool found')"
fi

echo "--------------------------------"
