#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
SCRIPT="$ROOT_DIR/scripts/voice-capture-openwakeword.sh"
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

SILENCE_WAV="$TMP_DIR/silence.wav"
arecord -q -D default -f S16_LE -r 16000 -c 1 -d 1 -t wav "$SILENCE_WAV"

OUTPUT=$("$SCRIPT" --input-wav "$SILENCE_WAV" --output-dir "$TMP_DIR/out" || true)
if [[ -n "$OUTPUT" ]]; then
  echo "expected no detection on silence, got: $OUTPUT" >&2
  exit 1
fi

echo "voice capture smoke ok"
