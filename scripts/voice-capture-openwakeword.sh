#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
VENV_DIR="${AURELIA_VOICE_VENV:-$HOME/.aurelia/voice-capture/venv}"

if [[ ! -x "$VENV_DIR/bin/python" ]]; then
  echo "voice capture venv missing: $VENV_DIR" >&2
  echo "run: $ROOT_DIR/scripts/setup-voice-capture-env.sh" >&2
  exit 1
fi

exec "$VENV_DIR/bin/python" "$ROOT_DIR/scripts/voice-capture-openwakeword.py" "$@"
