#!/bin/bash
# =============================================================================
# setup-kokoro-gpu.sh — Kokoro/Kodoro TTS GPU Setup
# =============================================================================
# Verifica, inicia e valida o motor Kokoro TTS com GPU NVIDIA.
# Requer: Kokoro ou Kodoro instalado (pip install kokoro-ng ou kodor)
# =============================================================================

set -euo pipefail

KOKORO_BASE_URL="${KOKORO_BASE_URL:-http://127.0.0.1:8012}"
VOICE_ID="${VOICE_ID:-aurelia-jarvis}"
TEST_TEXT="Olá, eu sou a Aurélia, assistente soberana do ecossistema."
TIMEOUT_HEALTH=5

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# =============================================================================
# 1. Verificar GPU NVIDIA
# =============================================================================
check_gpu() {
    log_info "Verificando GPU NVIDIA..."
    if command -v nvidia-smi &>/dev/null; then
        local gpu_name
        gpu_name=$(nvidia-smi --query-gpu=name --format=csv,noheader 2>/dev/null | head -1 || echo "GPU NVIDIA")
        local vram
        vram=$(nvidia-smi --query-gpu=memory.total --format=csv,noheader 2>/dev/null | head -1 | tr -d ' MiB' || echo "N/A")
        log_info "GPU: $gpu_name | VRAM: ${vram} MiB"
    else
        log_warn "nvidia-smi não encontrado. Kokoro pode funcionar em CPU (mais lento)."
    fi
}

# =============================================================================
# 2. Verificar se Kokoro/Kodoro está a correr
# =============================================================================
check_service() {
    log_info "Verificando Kokoro/Kodoro em $KOKORO_BASE_URL..."
    if curl -sf --max-time "$TIMEOUT_HEALTH" "$KOKORO_BASE_URL/health" &>/dev/null; then
        log_info "Kokoro/Kodoro está a correr em $KOKORO_BASE_URL"
        return 0
    fi

    log_warn "Kokoro/Kodoro não está a correr em $KOKORO_BASE_URL"
    log_info "Opções de instalação:"

    if command -v pip &>/dev/null || command -v pip3 &>/dev/null; then
        echo "  pip install kodor     # recomendada (fork ativa)"
        echo "  pip install kokoro-ng  # original"
    fi
    echo "  ou use o Docker:"
    echo "  docker run -d -p 8012:5001 ghcr.io/remsky/kokoro-onnx:latest"

    log_info "Após instalar, inicia com:"
    echo "  python -m kodor --port 5001 --device cuda  # GPU"
    echo "  python -m kodor --port 5001 --device cpu    # CPU"

    return 1
}

# =============================================================================
# 3. Verificar se a voz aurelia-jarvis existe
# =============================================================================
check_voice() {
    log_info "Verificando voz '$VOICE_ID'..."

    local voices
    voices=$(curl -sf --max-time 5 "$KOKORO_BASE_URL/v1/voices" 2>/dev/null | \
        grep -o "\"name\":\"[^\"]*\"" | sed 's/"name":"//g' | sed 's/"//g' || echo "")

    if echo "$voices" | grep -q "$VOICE_ID"; then
        log_info "Voz '$VOICE_ID' encontrada ✓"
        return 0
    fi

    log_warn "Voz '$VOICE_ID' não encontrada na lista de vozes disponíveis."
    if [ -n "$voices" ]; then
        log_info "Vozes disponíveis:"
        echo "$voices" | while read -r v; do echo "  - $v"; done
    fi
    log_info "Para usar voz PT-BR, grava e carrega em: voices/aurelia-jarvis.bin"
    return 1
}

# =============================================================================
# 4. Testar sintetização
# =============================================================================
test_synthesis() {
    log_info "Testando sintetização PT-BR..."

    local response
    local http_code
    http_code=$(curl -sf --max-time 30 \
        -X POST "$KOKORO_BASE_URL/v1/audio/speech" \
        -H "Content-Type: application/json" \
        -d "{\"model\":\"kokoro\",\"input\":\"$TEST_TEXT\",\"voice\":\"$VOICE_ID\",\"response_format\":\"mp3\"}" \
        --output /tmp/aurelia_tts_test.mp3 \
        -w "%{http_code}" 2>/dev/null || echo "000")

    if [ "$http_code" = "200" ] && [ -s /tmp/aurelia_tts_test.mp3 ]; then
        local size
        size=$(du -h /tmp/aurelia_tts_test.mp3 | cut -f1)
        log_info "Síntese OK — ${size} — guardado em /tmp/aurelia_tts_test.mp3"
        rm -f /tmp/aurelia_tts_test.mp3
        return 0
    fi

    log_error "Síntese falhou (HTTP $http_code)"
    log_info "Verifica se o servidor está em modo CUDA: python -m kodor --device cuda"
    return 1
}

# =============================================================================
# 5. Instruções de configuração da variável de ambiente
# =============================================================================
print_env() {
    log_info "Variáveis de ambiente para o .env:"
    echo ""
    echo "  # Kokoro TTS (PT-BR)"
    echo "  export TTS_BASE_URL=$KOKORO_BASE_URL"
    echo "  export TTS_MODEL=kokoro"
    echo "  export VOICE_ID=$VOICE_ID"
    echo "  export TTS_FORMAT=mp3"
    echo "  export TTS_SPEED=1.0"
    echo ""
    log_info "Para usar no aurelia, adiciona ao ~/.aurelia/config/app.json:"
    echo '  { "tts_base_url": "'"$KOKORO_BASE_URL"'", "tts_model": "kokoro", "tts_voice": "'"$VOICE_ID"'" }'
}

# =============================================================================
# MAIN
# =============================================================================
main() {
    echo ""
    echo "═══════════════════════════════════════════════"
    echo "  Kokoro/Kodoro TTS GPU Setup — Aurélia 2026"
    echo "═══════════════════════════════════════════════"
    echo ""

    check_gpu
    echo ""

    local service_ok=true
    check_service || service_ok=false
    echo ""

    if [ "$service_ok" = "true" ]; then
        check_voice || true
        echo ""
        test_synthesis || true
    else
        log_warn "Kokoro/Kodoro não está a correr. Instala primeiro."
    fi

    echo ""
    print_env
    echo ""
    log_info "Feito."
}

main "$@"
