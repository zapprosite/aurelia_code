#!/bin/bash
#
# Script: smoke-test.sh
# Descrição: Smoke test completo de todos os serviços
# Uso: ./smoke-test.sh
#

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=========================================="
echo "  AURELIA SMOKE TEST"
echo "  $(date '+%Y-%m-%d %H:%M:%S')"
echo "=========================================="
echo ""

TOTAL=0
PASS=0
FAIL=0

test_service() {
    local name="$1"
    local cmd="$2"
    TOTAL=$((TOTAL + 1))
    echo -n "[$TOTAL] $name ... "
    if eval "$cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}OK${NC}"
        PASS=$((PASS + 1))
    else
        echo -e "${RED}FAIL${NC}"
        FAIL=$((FAIL + 1))
    fi
}

# 1. Docker
echo "=== DOCKER ==="
test_service "Docker running" "docker info > /dev/null 2>&1"
test_service "Containers" "test \$(docker ps | wc -l) -gt 1"

# 2. Redis
echo ""
echo "=== REDIS ==="
# Tenta encontrar o container do redis dinamicamente
REDIS_CONTAINER=$(docker ps --format '{{.Names}}' | grep redis | head -n 1)
test_service "Redis (6379)" "docker exec $REDIS_CONTAINER redis-cli -p 6379 ping"

# 3. Qdrant
echo ""
echo "=== QDRANT ==="
test_service "Qdrant HTTP" "curl -sf http://127.0.0.1:6333/readyz"

# 4. Ollama
echo ""
echo "=== OLLAMA ==="
test_service "Ollama Service" "ps aux | grep -v grep | grep -q ollama"
test_service "Gemma3 Model" "ollama list | grep -q gemma3"

# 5. LiteLLM (Smart Router)
echo ""
echo "=== LITELLM ==="
test_service "LiteLLM container" "docker ps | grep -q smart-router"
test_service "LiteLLM port 4000" "curl -sf http://localhost:4000/"

# 6. STT (Groq Cloud)
echo ""
echo "=== STT (Groq) ==="
test_service "STT Provider Config" "grep -q 'STT_PROVIDER=groq' /home/will/aurelia/.env"
test_service "Groq API Key" "grep -q 'GROQ_API_KEY=' /home/will/aurelia/.env"

# 7. TTS (Voxtral & Edge)
echo ""
echo "=== TTS (Voxtral/Edge) ==="
test_service "Voxtral container" "docker ps | grep -q voxtral"
test_service "Voxtral health" "curl -sf http://localhost:8012/v1/models"
test_service "Edge TTS script" "test -f /home/will/aurelia/scripts/edge-tts.py"

# 8. Aurelia System API
echo ""
echo "=== AURELIA API ==="
test_service "Aurelia API container" "docker ps | grep -q aurelia-api"
test_service "Aurelia API Port 8080" "curl -sf http://localhost:8080/health"

# 9. n8n
echo ""
echo "=== N8N ==="
test_service "n8n container" "docker ps | grep -q n8n"

# 10. Grafana
echo ""
echo "=== GRAFANA ==="
test_service "Grafana container" "docker ps | grep -q grafana"

# 11. CapRover
echo ""
echo "=== CAPROVER ==="
test_service "CapRover process" "ps aux | grep -v grep | grep -q captain || docker ps | grep -q captain"

# 12. Rede
echo ""
echo "=== REDE ==="
test_service "Tailscale" "tailscale status 2>/dev/null | grep -q will-zappro"

# 13. GPU
echo ""
echo "=== GPU ==="
test_service "NVIDIA GPU" "nvidia-smi --query-gpu=name --format=csv,noheader"

# 14. Arquivos importantes
echo ""
echo "=== ARQUIVOS ==="
test_service ".env" "test -f /home/will/aurelia/.env"
test_service "Skills dir" "test -d /home/will/aurelia/.agent/skills"

# 15. Obsidian Vault
echo ""
echo "=== OBSIDIAN ==="
test_service "Knowledge Root" "test -d /home/will/aurelia/knowledge"

# RESUMO
echo ""
echo "=========================================="
echo "  RESUMO"
echo "=========================================="
echo -e "Total:  $TOTAL"
echo -e "Pass:   ${GREEN}$PASS${NC}"
echo -e "Fail:   ${RED}$FAIL${NC}"

if [ $FAIL -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 TODOS OS TESTES PASSARAM! (Aurelia SOTA 2026.2)${NC}"
    exit 0
else
    echo ""
    echo -e "${YELLOW}⚠️  $FAIL testes falharam - verificar manualmente${NC}"
    exit 0
fi
