#!/bin/bash
# Script de Verificação de Saúde do Homelab (SOTA 2026)
# Valida se os serviços essenciais estão rodando nas portas da ADR 20240401

echo "--- Homelab Health Check (Governança de Portas) ---"

check_port() {
    local port=$1
    local name=$2
    if ss -tunlp | grep -q ":$port " ; then
        echo -e "[OK] $name está rodando na porta $port"
    else
        echo -e "[ERRO] $name NÃO detectado na porta $port"
    fi
}

echo "Verificando serviços críticos..."
check_port 3000 "CapRover UI"
check_port 3001 "Grafana"
check_port 6333 "Qdrant REST"
check_port 4000 "LiteLLM"
check_port 11434 "Ollama"
check_port 5678 "n8n"
check_port 80 "CapRover Ingress (HTTP)"
check_port 443 "CapRover Ingress (HTTPS)"

echo "---------------------------------------------------"
echo "Verificando containers Docker ativos..."
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "qdrant|caprover|litellm|ollama|grafana|n8n" || echo "Nenhum container da stack homelab encontrado."

echo "--- Fim do Check ---"
