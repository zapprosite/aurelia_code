#!/usr/bin/env bash
set -euo pipefail

VENV_DIR="${AURELIA_VOICE_VENV:-$HOME/.aurelia/voice-capture/venv}"
PYTHON_BIN="${PYTHON_BIN:-python3}"

"$PYTHON_BIN" -m venv "$VENV_DIR"
source "$VENV_DIR/bin/activate"
python -m pip install --upgrade pip setuptools wheel
python -m pip install openwakeword==0.4.0

echo "voice capture env ready: $VENV_DIR"
