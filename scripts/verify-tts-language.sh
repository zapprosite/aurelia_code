#!/usr/bin/env bash
set -euo pipefail

instance_root=${AURELIA_HOME:-$HOME/.aurelia}
config_path="$instance_root/config/app.json"

if [[ ! -f "$config_path" ]]; then
  echo "app config not found: $config_path" >&2
  exit 1
fi

read_language_python() {
  cat <<'PY' | python3 - 2>/dev/null || python - 2>/dev/null || true
import json
import pathlib
raw = pathlib.Path(r"${config_path}").read_text()
data = json.loads(raw)
print(data.get("tts_language", ""))
PY
}

language=$(read_language_python)

if [[ -z "$language" ]]; then
  echo "Failed to read tts_language from $config_path" >&2
  exit 3
fi

if [[ "$language" != "pt" ]]; then
  echo "TTS language must be pt, current value: $language" >&2
  exit 2
fi

echo "TTS language enforcement: PT-BR confirmed ($config_path)"

if [[ "${VERIFY_TTS_ENDPOINT:-1}" -ne 0 ]]; then
  echo "Validating synthesize endpoint..."
  curl -fsS --retry 3 --retry-delay 1 \
    --data '{"text":"verificando idioma PTBR"}' \
    http://127.0.0.1:8484/v1/voice/synthesize >/dev/null
  echo "Endpoint responded with HTTP 200"
else
  echo "Endpoint validation skipped (set VERIFY_TTS_ENDPOINT=1 to enable)"
fi
