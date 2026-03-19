#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Uso:
  scripts/minimax-voice-list.sh [--type all|system|voice_cloning|voice_generation] [--dry-run]

Objetivo:
  Listar vozes disponiveis na conta MiniMax para selecionar a voz da Aurelia.
EOF
}

VOICE_TYPE="all"
DRY_RUN=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --type)
      VOICE_TYPE="${2:-all}"
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
  echo "curl -sS -X POST https://api.minimax.io/v1/get_voice -H 'Authorization: Bearer <MINIMAX_API_KEY>' -H 'Content-Type: application/json' -d '{\"voice_type\":\"$VOICE_TYPE\"}'"
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

curl -sS -X POST "https://api.minimax.io/v1/get_voice" \
  -H "Authorization: Bearer $key" \
  -H "Content-Type: application/json" \
  -d "{\"voice_type\":\"$VOICE_TYPE\"}" | jq '.'
