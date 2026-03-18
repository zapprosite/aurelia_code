#!/bin/bash
# Senior Health-Check Script v1.0
# Optimized for Ubuntu 24.04 + RTX 4090

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== GPU (NVIDIA RTX 4090) ===${NC}"
nvidia-smi --query-gpu=name,memory.total,memory.used,temperature.gpu,utilization.gpu --format=csv,noheader

echo -e "\n${GREEN}=== DOCKER SERVICES ===${NC}"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Size}}"

echo -e "\n${GREEN}=== NATIVE SERVICES ===${NC}"
pgrep qdrant > /dev/null && echo -e "Qdrant: ${GREEN}ONLINE${NC} (Native)" || echo -e "Qdrant: ${RED}OFFLINE${NC}"
pgrep ollama > /dev/null && echo -e "Ollama: ${GREEN}ONLINE${NC} (Native)" || echo -e "Ollama: ${RED}OFFLINE${NC}"

echo -e "\n${GREEN}=== OLLAMA MODELS ===${NC}"
curl -s http://localhost:11434/api/tags | jq -r '.models[].name'

echo -e "\n${GREEN}=== INFRA HEALTH ===${NC}"
echo -n "Qdrant API: " && curl -s http://localhost:6333/readyz || echo "ERROR"
echo -n "LiteLLM API: " && curl -s http://localhost:4000/health/liveliness || echo "ERROR"
echo -n "ZFS Pool (tank): " && zpool list -H -o health tank

echo -e "\n${GREEN}=== DB INTEGRITY ===${NC}"
python3 -c "import sqlite3; conn = sqlite3.connect('/home/will/.aurelia/data/aurelia.db'); print('Aurelia DB Journal Mode:', conn.execute('PRAGMA journal_mode').fetchone()[0])"
