#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Uso:
  scripts/gemini-tts-smoke.sh [--voice Sulafat] [--text "texto"] [--output /tmp/out.wav] [--dry-run]

Objetivo:
  Validar Gemini TTS com uma voz pronta, sem mudar o runtime ativo.
EOF
}

VOICE="Sulafat"
TEXT="Olá. Eu sou Aurélia. Estou pronta para ajudar com clareza, serenidade e profissionalismo."
OUTPUT="/tmp/aurelia-gemini-tts.wav"
DRY_RUN=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --voice)
      VOICE="${2:-}"
      shift 2
      ;;
    --text)
      TEXT="${2:-}"
      shift 2
      ;;
    --output)
      OUTPUT="${2:-}"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "[ERROR] argumento desconhecido: $1" >&2
      exit 1
      ;;
  esac
done

if [[ "$DRY_RUN" -eq 1 ]]; then
  echo "POST https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-preview-tts:generateContent voice=$VOICE format=wav"
  exit 0
fi

key="${GEMINI_API_KEY:-}"
if [[ -z "$key" && -f "$HOME/.aurelia/config/app.json" ]]; then
  key="$(jq -r '.google_api_key // ""' "$HOME/.aurelia/config/app.json")"
fi
if [[ -z "$key" ]]; then
  echo "[ERROR] GEMINI_API_KEY ausente" >&2
  exit 1
fi

resp="$(mktemp -t aurelia-gemini-tts-XXXXXX.json)"
pcm="$(mktemp -t aurelia-gemini-tts-XXXXXX.pcm)"
trap 'rm -f "$resp" "$pcm"' EXIT

curl -sS "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-preview-tts:generateContent" \
  -H "x-goog-api-key: $key" \
  -X POST \
  -H "Content-Type: application/json" \
  -d "$(jq -nc --arg text "Fale em português do Brasil com voz feminina, tom doce, calmo e acolhedor, dicção clara e elegante, sem gírias, sem regionalismos informais e sem portunhol. Ritmo pausado, equilibrado e profissional. Texto: $TEXT" --arg voice "$VOICE" '{
    contents:[{parts:[{text:$text}]}],
    generationConfig:{
      responseModalities:["AUDIO"],
      speechConfig:{
        voiceConfig:{
          prebuiltVoiceConfig:{voiceName:$voice}
        }
      }
    },
    model:"gemini-2.5-flash-preview-tts"
  }')" > "$resp"

jq -r '.candidates[0].content.parts[0].inlineData.data // empty' "$resp" | base64 --decode > "$pcm"
if [[ ! -s "$pcm" ]]; then
  echo "[ERROR] resposta sem áudio PCM" >&2
  cat "$resp" >&2
  exit 1
fi

python3 - "$pcm" "$OUTPUT" <<'PY'
import sys, wave
pcm_path, wav_path = sys.argv[1], sys.argv[2]
with open(pcm_path, "rb") as f:
    pcm = f.read()
with wave.open(wav_path, "wb") as wf:
    wf.setnchannels(1)
    wf.setsampwidth(2)
    wf.setframerate(24000)
    wf.writeframes(pcm)
print(wav_path)
PY
