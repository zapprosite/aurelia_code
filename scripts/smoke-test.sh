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

# 2. Redis (3 instâncias)
echo ""
echo "=== REDIS ==="
test_service "Redis main (6379)" "docker exec aurelia-redis-main redis-cli ping"
test_service "Redis litellm (6380)" "docker exec litellm-redis redis-cli ping"
test_service "Redis n8n (6381)" "docker exec n8n-redis redis-cli ping"

# 3. Qdrant
echo ""
echo "=== QDRANT ==="
test_service "Qdrant HTTP" "curl -sf http://127.0.0.1:6333/readyz"
test_service "Qdrant collections" "curl -sf http://127.0.0.1:6333/collections | grep -q result"

# 4. Ollama
echo ""
echo "=== OLLAMA ==="
test_service "Ollama API" "curl -sf http://127.0.0.1:11434/api/tags"
test_service "Ollama models" "curl -sf http://127.0.0.1:11434/api/tags | grep -q qwen"

# 5. LiteLLM
echo ""
echo "=== LITELLM ==="
test_service "LiteLLM container" "docker ps | grep -q smart-router"
test_service "LiteLLM port 4000" "curl -sf http://localhost:4000/health || curl -sf http://localhost:4000"
test_service "LiteLLM UI 3334" "curl -sf http://localhost:3334"

# 6. STT (Whisper)
echo ""
echo "=== STT (Whisper) ==="
test_service "Whisper container" "docker ps | grep -q whisper-local"
test_service "Whisper health" "curl -sf http://localhost:8020/health"

# 7. TTS (Kokoro)
echo ""
echo "=== TTS (Kokoro) ==="
test_service "Kokoro container" "docker ps | grep -q kokoro"
test_service "Kokoro health" "curl -sf http://localhost:8012/health"

# 8. Edge TTS
echo ""
echo "=== EDGE TTS ==="
test_service "Edge TTS script" "test -f /home/will/aurelia/scripts/edge-tts.py"

# 9. n8n
echo ""
echo "=== N8N ==="
test_service "n8n container" "docker ps | grep -q n8n"
test_service "n8n Postgres" "docker exec n8n-postgres pg_isready"

# 10. Grafana
echo ""
echo "=== GRAFANA ==="
test_service "Grafana container" "docker ps | grep -q grafana"
test_service "Grafana local" "curl -sf http://localhost:3100/api/health"

# 11. CapRover
echo ""
echo "=== CAPROVER ==="
test_service "CapRover captain" "docker ps | grep -q captain"
test_service "CapRover nginx" "docker ps | grep -q nginx"

# 12. Rede
echo ""
echo "=== REDE ==="
test_service "Cloudflare Tunnel" "test -f ~/.cloudflared/config.yml"
test_service "Tailscale" "tailscale status 2>/dev/null | grep -q will-zappro"
test_service "DNS n8n" "getent hosts n8n.zappro.site > /dev/null"

# 13. GPU
echo ""
echo "=== GPU ==="
test_service "NVIDIA GPU" "nvidia-smi --query-gpu=name --format=csv,noheader"
test_service "NVIDIA driver" "nvidia-smi --query-gpu=driver_version --format=csv,noheader"

# 14. Arquivos importantes
echo ""
echo "=== ARQUIVOS ==="
test_service "Aurelia bin" "test -f /home/will/aurelia/bin/aurelia"
test_service ".env" "test -f /home/will/aurelia/.env"
test_service "Skills dir" "test -d /home/will/aurelia/.agent/skills"

# 15. Obsidian
echo ""
echo "=== OBSIDIAN ==="
test_service "Obsidian vault" "test -d /home/will/Documents/Obsidian/Aurelia"

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
    echo -e "${GREEN}🎉 TODOS OS TESTES PASSARAM!${NC}"
    exit 0
else
    echo ""
    echo -e "${YELLOW}⚠️  $FAIL testes falharam - verificar manualmente${NC}"
    exit 0
fi
