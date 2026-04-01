#!/bin/bash
# Script de Gestão de Estabilidade Homelab (Ubuntu Desktop)
# Habilita persistência e garante autorestart de serviços

echo "--- Homelab Stability Manager ---"

# 1. Habilitar Ollama no systemd
if command -v ollama &> /dev/null; then
    echo "[1/3] Habilitando Ollama no boot (systemd)..."
    sudo systemctl enable ollama 2>/dev/null
    sudo systemctl start ollama 2>/dev/null
fi

# 2. Configurar Auto-restart para Containers Docker
echo "[2/3] Configurando restart:always para containers críticos..."
CONTAINERS=$(docker ps -a --format "{{.Names}}" | grep -E "qdrant|caprover|litellm|grafana|n8n")
for c in $CONTAINERS; do
    echo "  - Aplicando em $c"
    docker update --restart unless-stopped "$c"
done

# 3. Limpeza de Logs do Docker (prevenir estouro de disco)
echo "[3/3] Limpando logs Docker antigos..."
sudo truncate -s 0 /var/lib/docker/containers/*/*-json.log 2>/dev/null || echo "Aviso: Sem permissão para truncar logs."

echo "--- Estabilização Concluída! ---"
