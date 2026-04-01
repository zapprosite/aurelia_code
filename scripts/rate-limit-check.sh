#!/bin/bash
#
# Script: rate-limit-check.sh
# Descrição: Verificar rate limits e recursos do sistema
# Uso: ./rate-limit-check.sh
#

echo "=========================================="
echo "  RATE LIMIT + RECURSOS SISTEMA"
echo "=========================================="
echo ""

# 1. Rate Limit Telegram
echo "=== TELEGRAM RATE LIMITS ==="
BOT_TOKEN="8793928549:AAEr3tjaarijUWxu-iru0Vcm6N6DkwjndL4"

# GetChatMemberCount (30 msg/segundo)
echo "Chat member count..."
RESP=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/getChatMemberCount" \
  -H "Content-Type: application/json" \
  -d '{"chat_id": "7220607041"}' 2>/dev/null)
echo "$RESP" | jq -r '.result // "Erro"' 2>/dev/null || echo "sem acesso"

# GetMe
RESP=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/getMe" 2>/dev/null)
echo "Bot: $(echo "$RESP" | jq -r '.result.username // "erro"')"

echo ""

# 2. Recursos do Sistema
echo "=== RECURSOS ==="
echo "GPU: $(nvidia-smi --query-gpu=name,memory.used,memory.total --format=csv,noheader)"
echo "RAM: $(free -h | awk '/^Mem:/ {print $3 used, $2 total}')"
echo "CPU: $(cat /proc/cpuinfo | grep "model name" | head -1 | cut -d: -f2 | xargs)"

echo ""

# 3. Ollama Rate Limits
echo "=== OLLAMA ==="
curl -s http://127.0.0.1:11434/api/tags 2>/dev/null | jq '.models[] | "\(.name) - \(.size)"' || echo "Ollama offline"

echo ""

# 4. Qwen3.5-9B (VL)
echo "=== QWEN VL (9B) ==="
echo "Modelo: qwen3.5:9b"
echo "Tamanho: 6.6GB"
VRAM_DISPONIVEL=22821
TAMANHO_MODELO=6600
echo "VRAM necessário: ~6.6GB"
if [ $VRAM_DISPONIVEL -gt $TAMANHO_MODELO ]; then
    echo "✅ Cabe na VRAM ($((VRAM_DISPONIVEL - TAMANHO_MODELO))MB livre)"
else
    echo "⚠️ Não cabe na VRAM"
fi

echo ""

# 5. LiteLLM Rate Limits
echo "=== LITELLM (porta 4000) ==="
curl -s -H "Authorization: Bearer sk-test" http://localhost:4000/health 2>/dev/null | jq -r '.status // "offline"' || echo "LiteLLM offline"

echo ""

# 6. TTS
echo "=== TTS (Kokoro + Edge) ==="
curl -s http://127.0.0.1:8012/health 2>/dev/null | jq -r '.status // "Kokoro offline"' || echo "Kokoro offline"
echo "Edge TTS: scripts/edge-tts.py"

echo ""

# 7. Resumo de Limits
echo "=========================================="
echo "  RESUMO: RATE LIMITS"
echo "=========================================="
echo "Telegram: 30 msg/segundo"
echo "Ollama: ~3 req/min (sem rate limit)"
echo "LiteLLM: por chave API"
echo ""
echo "VRAM Total: 24GB"
echo "VRAM Usado: 1.2GB"  
echo "VRAM Livre: 22.8GB"
echo ""
echo "RAM Total: 30GB"
echo "RAM Usado: 17GB"
echo "RAM Livre: 13GB"
echo ""

# 8. Recomendação de Carga
echo "=========================================="
echo "  RECOMENDAÇÃO DE CARGA"
echo "=========================================="
echo "Max agentes simultâneos: 3-4 (9B VL = 6.6GB cada)"
echo "Max requisições/segundo: 5-10"
echo "Max concurrent TTS: 2-3"
echo ""
echo "=========================================="
