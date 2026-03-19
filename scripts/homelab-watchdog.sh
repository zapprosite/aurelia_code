#!/bin/bash
# 🛸 Home Lab Watchdog — Auto-Reparação Proativa
# Este script verifica a saúde dos containers e reinicia os que estiverem instáveis.

echo "--- Iniciando Watchdog do Home Lab ($(date)) ---"

# 1. Verificar Containers Docker
UNHEALTHY_CONTAINERS=$(docker ps -f "health=unhealthy" --format "{{.Names}}")

if [ -n "$UNHEALTHY_CONTAINERS" ]; then
    echo "⚠️ Containers UNHEALTHY detectados: $UNHEALTHY_CONTAINERS"
    for container in $UNHEALTHY_CONTAINERS; do
        echo "🔄 Reiniciando $container..."
        docker restart "$container"
    done
else
    echo "✅ Todos os containers monitorados estão saudáveis."
fi

# 2. Verificar Status do Ollama
if curl -s http://localhost:11434/api/tags > /dev/null; then
    echo "✅ Ollama Online."
else
    echo "❌ Ollama Offline! Tentando reiniciar serviço..."
    sudo systemctl restart ollama
fi

# 3. Verificar GPU NVIDIA
if nvidia-smi > /dev/null 2>&1; then
    echo "✅ GPU NVIDIA Operacional."
else
    echo "❌ Erro de comunicação com a GPU!"
fi

echo "--- Watchdog Concluído ---"
