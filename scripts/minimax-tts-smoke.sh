#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Uso:
  scripts/minimax-tts-smoke.sh --voice-id <voice_id> [--text <texto>] [--output <arquivo.mp3>] [--dry-run]

Objetivo:
  Validar a lane TTS da MiniMax para a voz oficial da Aurelia em PT-BR.
EOF
}

VOICE_ID=""
TEXT="Olá. Eu sou Aurélia. Estou pronta para ajudar com clareza, serenidade e profissionalismo."
OUTPUT=""
DRY_RUN=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --voice-id)
      VOICE_ID="${2:-}"
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

if [[ -z "$VOICE_ID" ]]; then
  echo "[ERROR] --voice-id e obrigatorio" >&2
  exit 1
fi

if [[ "$DRY_RUN" -eq 1 ]]; then
  echo "POST https://api.minimax.io/v1/t2a_v2 model=speech-2.8-hd voice_id=$VOICE_ID language_boost=Portuguese format=mp3"
  exit 0
fi

key="${MINIMAX_API_KEY:-}"
if [[ -z "$key" && -f "$HOME/.aurelia/config/app.json" ]]; then
  key="$(jq -r '.minimax_api_key // ""' "$HOME/.aurelia/config/app.json")"
fi
if [[ -z "$key" ]]; then
  echo "[ERROR] MINIMAX_API_KEY ausente" >&2
  exit 1
fi

tmpjson="$(mktemp -t aurelia-minimax-tts-XXXXXX.json)"
trap 'rm -f "$tmpjson"' EXIT

curl -sS -X POST "https://api.minimax.io/v1/t2a_v2" \
  -H "Authorization: Bearer $key" \
  -H "Content-Type: application/json" \
  -d "$(jq -nc \
      --arg model "speech-2.8-hd" \
      --arg text "$TEXT" \
      --arg voice "$VOICE_ID" \
      '{
        model:$model,
        text:$text,
        stream:false,
        language_boost:"Portuguese",
        output_format:"hex",
        voice_setting:{voice_id:$voice,speed:1,vol:1,pitch:0},
        audio_setting:{sample_rate:32000,bitrate:128000,format:"mp3",channel:1}
      }')" >"$tmpjson"

if [[ "$(jq -r '.base_resp.status_code // -1' "$tmpjson")" != "0" ]]; then
  cat "$tmpjson" >&2
  exit 1
fi

if [[ -n "$OUTPUT" ]]; then
  jq -r '.data.audio' "$tmpjson" | xxd -r -p >"$OUTPUT"
  echo "$OUTPUT"
  exit 0
fi

jq -r '.data.audio' "$tmpjson" | xxd -r -p >/dev/stdout
