#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
CONFIG_PATH="${AURELIA_APP_CONFIG:-$HOME/.aurelia/config/app.json}"
DATA_DIR="${AURELIA_DATA_DIR:-$HOME/.aurelia/data}"
STATUS_PATH="${GEMINI_STATUS_PATH:-$DATA_DIR/gemini_smoke.json}"
MODEL="${GEMINI_MODEL:-gemini-2.5-flash}"
PROMPT="${GEMINI_SMOKE_PROMPT:-reply with ok only}"

read_key_from_config() {
  if [[ -f "$CONFIG_PATH" ]]; then
    jq -r '.google_api_key // empty' "$CONFIG_PATH"
  fi
}

API_KEY="${GOOGLE_API_KEY:-$(read_key_from_config)}"

write_status() {
  local status="$1"
  local error_message="${2:-}"
  mkdir -p "$DATA_DIR"
  jq -nc \
    --arg status "$status" \
    --arg model "$MODEL" \
    --arg checked_at "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
    --arg error "$error_message" \
    '{
      status: $status,
      model: $model,
      checked_at: $checked_at,
      error: (if $error == "" then null else $error end)
    }' > "${STATUS_PATH}.tmp"
  install -m 600 "${STATUS_PATH}.tmp" "$STATUS_PATH"
  rm -f "${STATUS_PATH}.tmp"
}

usage() {
  cat <<'EOF'
Usage:
  GOOGLE_API_KEY=... scripts/gemini-smoke.sh

Optional env:
  GEMINI_MODEL         default: gemini-2.5-flash
  GEMINI_SMOKE_PROMPT  default: reply with ok only
  AURELIA_APP_CONFIG   default: ~/.aurelia/config/app.json
EOF
}

if [[ -z "$API_KEY" ]]; then
  usage
  echo "[ERROR] missing Google API key in env or runtime config." >&2
  exit 1
fi

echo "==> list models"
MODELS_OUTPUT="$(curl -fsS "https://generativelanguage.googleapis.com/v1beta/models?key=${API_KEY}" \
  | jq -r '.models[]?.name' \
  | rg '^models/gemini-' \
  | sed -n '1,12p')" || {
    write_status "error" "list_models_failed"
    exit 1
  }
printf '%s\n' "$MODELS_OUTPUT"

echo
echo "==> generate content (${MODEL})"
RESPONSE_TEXT="$(curl -fsS \
  -H 'Content-Type: application/json' \
  -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}" \
  -d "$(jq -nc --arg prompt "$PROMPT" '{contents:[{parts:[{text:$prompt}]}]}')" \
  | jq -r '.candidates[0].content.parts[0].text // empty')" || {
    write_status "error" "generate_content_failed"
    exit 1
  }
printf '%s\n' "$RESPONSE_TEXT"
write_status "ok"
