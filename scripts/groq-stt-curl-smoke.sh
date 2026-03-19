#!/usr/bin/env bash
set -euo pipefail

API_BASE="${GROQ_API_BASE:-https://api.groq.com/openai/v1}"
MODEL="${GROQ_STT_MODEL:-whisper-large-v3-turbo}"
AUDIO_FILE="${1:-}"

usage() {
  cat <<'EOF'
Usage:
  GROQ_API_KEY=... scripts/groq-stt-curl-smoke.sh /path/to/audio.wav

Behavior:
  - without GROQ_API_KEY: prints the exact curl command as a dry-run
  - with GROQ_API_KEY: executes the request against Groq STT

Optional env:
  GROQ_API_BASE   default: https://api.groq.com/openai/v1
  GROQ_STT_MODEL  default: whisper-large-v3-turbo
EOF
}

if [[ -z "${AUDIO_FILE}" ]]; then
  usage
  exit 1
fi

if [[ ! -f "${AUDIO_FILE}" ]]; then
  echo "[ERROR] audio file not found: ${AUDIO_FILE}" >&2
  exit 1
fi

CURL_CMD=(
  curl -sS
  -X POST "${API_BASE}/audio/transcriptions"
  -H "Authorization: Bearer ${GROQ_API_KEY:-<YOUR_GROQ_API_KEY>}"
  -F "file=@${AUDIO_FILE}"
  -F "model=${MODEL}"
  -F "response_format=json"
)

if [[ -z "${GROQ_API_KEY:-}" ]]; then
  echo "[DRY-RUN] GROQ_API_KEY is not set."
  printf '%q ' "${CURL_CMD[@]}"
  echo
  exit 0
fi

"${CURL_CMD[@]}"
