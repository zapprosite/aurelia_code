#!/bin/bash
# =============================================================================
# setup-aurelia-voice.sh — Voice Cloning for Aurélia using KokoClone
# =============================================================================
# Clona a voz da Aurélia a partir de um sample de referência.
# Usa KokoClone para voice cloning zero-shot.
# =============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

KOKOCLONE_DIR="${KOKOCLONE_DIR:-/tmp/kokoclone}"
SAMPLE_WAV="${SAMPLE_WAV:-$PROJECT_ROOT/assets/voice/aurelia_sample.wav}"
REF_AUDIO="${REF_AUDIO:-$PROJECT_ROOT/assets/voice/aurelia.mp3}"
OUTPUT_DIR="${OUTPUT_DIR:-$PROJECT_ROOT/assets/voice/cloned}"
LANG="${LANG:-pt}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

check_kokoclone() {
    log_info "Verificando KokoClone..."
    if [ -d "$KOKOCLONE_DIR" ]; then
        log_info "KokoClone encontrado em $KOKOCLONE_DIR"
        return 0
    fi
    log_info "Clonando KokoClone..."
    git clone https://github.com/Ashish-Patnaik/kokoclone.git "$KOKOCLONE_DIR"
    cd "$KOKOCLONE_DIR"
    uv sync
    log_info "KokoClone instalado!"
    return 0
}

prepare_sample() {
    log_info "Preparando sample de áudio..."
    
    mkdir -p "$OUTPUT_DIR"
    
    if [ ! -f "$REF_AUDIO" ]; then
        log_error "Áudio de referência não encontrado: $REF_AUDIO"
        return 1
    fi
    
    if [ -f "$SAMPLE_WAV" ] && [ "${FORCE_RECREATE:-false}" = "false" ]; then
        log_info "Sample já existe: $SAMPLE_WAV"
    else
        log_info "Extraindo 15s do áudio original..."
        ffmpeg -i "$REF_AUDIO" -ss 0 -t 15 \
            -acodec pcm_s16le -ar 16000 -ac 1 \
            "$SAMPLE_WAV" -y 2>/dev/null
        log_info "Sample guardado em: $SAMPLE_WAV"
    fi
    
    local size
    size=$(du -h "$SAMPLE_WAV" | cut -f1)
    log_info "Sample: $size"
}

clone_voice() {
    local text="$1"
    local output_file="$2"
    
    log_info "Clonando voz: $text"
    
    if [ ! -f "$SAMPLE_WAV" ]; then
        log_error "Sample não encontrado: $SAMPLE_WAV"
        return 1
    fi
    
    cd "$KOKOCLONE_DIR"
    source .venv/bin/activate
    
    python cli.py \
        --text "$text" \
        --lang "$LANG" \
        --ref "$SAMPLE_WAV" \
        --out "$output_file" 2>&1 | grep -E "(Success|Downloading|Synthesizing|Applying)"
    
    if [ -f "$output_file" ]; then
        log_info "Voz clonada guardada em: $output_file"
        return 0
    else
        log_error "Falha ao clonar voz"
        return 1
    fi
}

batch_clone() {
    log_info "Clonando batch de frases..."
    
    local phrases=(
        "Olá, eu sou a Aurélia, assistente soberana do ecossistema."
        "Como posso ajudá-lo hoje?"
        "Processando a sua solicitação."
        "Operação concluída com sucesso."
        "Desculpe, não entendi. Pode repetir?"
    )
    
    mkdir -p "$OUTPUT_DIR/phrases"
    
    for i in "${!phrases[@]}"; do
        local phrase="${phrases[$i]}"
        local output="$OUTPUT_DIR/phrases/phrase_$((i+1)).wav"
        clone_voice "$phrase" "$output" || true
    done
    
    log_info "Batch completo! Ficheiros em: $OUTPUT_DIR/phrases/"
}

main() {
    echo ""
    echo "═══════════════════════════════════════════════"
    echo "  Aurélia Voice Clone — KokoClone Setup"
    echo "═══════════════════════════════════════════════"
    echo ""
    
    check_kokoclone
    prepare_sample
    
    case "${1:-test}" in
        test)
            local test_text="Olá, eu sou a Aurélia, assistente soberana."
            local test_output="$OUTPUT_DIR/aurelia_test.wav"
            clone_voice "$test_text" "$test_output"
            echo ""
            echo "Teste concluído!"
            echo "Ouvir: ffplay $test_output"
            ;;
        batch)
            batch_clone
            ;;
        cli)
            shift
            local text="${1:-}"
            local output="${2:-$OUTPUT_DIR/custom.wav}"
            if [ -z "$text" ]; then
                log_error "Usage: $0 cli \"texto\" [output.wav]"
                exit 1
            fi
            clone_voice "$text" "$output"
            ;;
        *)
            echo "Usage: $0 {test|batch|cli}"
            echo ""
            echo "  test   - Teste com frase padrão"
            echo "  batch  - Clonar batch de frases"
            echo "  cli    - CLI interativo (use: cli \"texto\" [output])"
            ;;
    esac
}

main "$@"
