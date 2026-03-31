#!/bin/bash
# Registra voz clonada "aurelia-jarvis" no Voxtral TTS
# Ref: assets/voice/aurelia.mp3 — feminina, doce, PT-BR
set -euo pipefail

VOXTRAL_URL="${VOXTRAL_URL:-http://localhost:8012}"
VOICE_FILE="${1:-/home/will/aurelia/assets/voice/aurelia.mp3}"
VOICE_NAME="${2:-aurelia-jarvis}"

echo "Registrando voz '${VOICE_NAME}' de ${VOICE_FILE}..."

curl -sf -X POST "${VOXTRAL_URL}/v1/audio/voices" \
  -F "audio_sample=@${VOICE_FILE}" \
  -F "name=${VOICE_NAME}" \
  -F "ref_text=Olá, eu sou a Aurélia, sua assistente pessoal. Estou aqui para ajudar você com tudo que precisar." \
  -F "consent=aurelia-homelab-2026"

echo ""
echo "Voz '${VOICE_NAME}' registrada com sucesso."
echo "Teste: curl -X POST ${VOXTRAL_URL}/v1/audio/speech -H 'Content-Type: application/json' -d '{\"input\":\"Olá! O valor é R\$ 1.250,00.\",\"model\":\"mistralai/Voxtral-4B-TTS-2603\",\"voice\":\"${VOICE_NAME}\",\"response_format\":\"opus\"}' --output /tmp/test-voxtral.opus"
