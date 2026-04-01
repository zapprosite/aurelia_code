#!/bin/bash
# =============================================================================
# aurelia-tts.sh — Edge TTS + Kokoro GPU fallback
# =============================================================================
# Uses Microsoft Edge TTS (free, natural PT-BR) as primary
# Falls back to Kokoro GPU if Edge TTS fails
# =============================================================================

set -euo pipefail

VENV="${AURELIA_VENV:-/home/will/aurelia/.venv}"
EDGE_TTS="/home/will/aurelia/scripts/edge-tts.py"
KOKORO_URL="${KOKORO_URL:-http://127.0.0.1:8012}"

VOICE="${1:-}"
TEXT="${2:-}"
OUTPUT="${3:-}"

usage() {
    echo "Usage: $0 <voice> <text> <output.mp3>"
    echo ""
    echo "Voices (Edge TTS - free):"
    echo "  pt-BR-FranciscaNeural  (Female - recommended)"
    echo "  pt-BR-ThalitaMultilingualNeural (Female)"
    echo "  pt-BR-AntonioNeural (Male)"
    echo ""
    echo "Kokoro voices:"
    echo "  af_heart, af_sarah, af_jessica"
    exit 1
}

[[ -z "$VOICE" || -z "$TEXT" || -z "$OUTPUT" ]] && usage

# Try Edge TTS first (free, natural)
edge_tts() {
    source "$VENV/bin/activate" 2>/dev/null || true
    python3 "$EDGE_TTS" \
        --text "$TEXT" \
        --voice "$VOICE" \
        --output "$OUTPUT" 2>/dev/null
}

# Fallback to Kokoro GPU
kokoro_tts() {
    curl -sf --max-time 30 -X POST "${KOKORO_URL}/v1/audio/speech" \
        -H "Content-Type: application/json" \
        -d "{\"model\":\"kokoro\",\"input\":\"$TEXT\",\"voice\":\"$VOICE\",\"response_format\":\"mp3\"}" \
        -o "$OUTPUT"
}

# Main
if [[ "$VOICE" == pt-BR-* ]]; then
    # Edge TTS for PT-BR voices
    if edge_tts; then
        echo "Edge TTS OK: $OUTPUT"
        exit 0
    fi
    echo "Edge TTS failed, trying Kokoro..."
fi

# Kokoro fallback
if kokoro_tts; then
    echo "Kokoro OK: $OUTPUT"
    exit 0
fi

echo "ERROR: Both TTS engines failed"
exit 1
