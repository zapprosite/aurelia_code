#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Uso:
  scripts/media-transcript.sh --input <arquivo-ou-url> [--output <txt>] [--dry-run]

Objetivo:
  Extrair transcript de audio/video local ou de um link (ex.: YouTube) usando:
  - yt-dlp para download de audio quando a entrada for URL
  - ffmpeg para normalizar em WAV mono 16k
  - Groq STT para transcrever em PT-BR

Notas:
  - Requer: curl, jq, ffmpeg, ffprobe
  - Para URL, tambem requer: yt-dlp
  - Le GROQ_API_KEY do ambiente ou de ~/.aurelia/config/app.json
  - Nao usa o audio baixado como base de clonagem de voz; serve para transcript/estudo
EOF
}

INPUT=""
OUTPUT=""
DRY_RUN=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --input)
      INPUT="${2:-}"
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
      usage >&2
      exit 1
      ;;
  esac
done

if [[ -z "$INPUT" ]]; then
  echo "[ERROR] --input e obrigatorio" >&2
  usage >&2
  exit 1
fi

is_url=0
if [[ "$INPUT" =~ ^https?:// ]]; then
  is_url=1
fi

need_bin() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "[ERROR] dependência ausente: $1" >&2
    exit 1
  }
}

read_groq_key() {
  if [[ -n "${GROQ_API_KEY:-}" ]]; then
    printf '%s' "$GROQ_API_KEY"
    return
  fi
  if [[ -f "$HOME/.aurelia/config/app.json" ]]; then
    jq -r '.groq_api_key // ""' "$HOME/.aurelia/config/app.json"
    return
  fi
  printf ''
}

if [[ "$DRY_RUN" -eq 1 ]]; then
  echo "input=$INPUT"
  echo "mode=$([[ $is_url -eq 1 ]] && echo url || echo local)"
  if [[ "$is_url" -eq 1 ]]; then
    echo "step=yt-dlp -> bestaudio -> temp file"
  fi
  echo "step=ffmpeg -> mono 16k wav"
  echo "step=curl Groq /audio/transcriptions model=whisper-large-v3-turbo language=pt temperature=0"
  echo "step=stdout transcript or write to --output"
  exit 0
fi

need_bin curl
need_bin jq
need_bin ffmpeg
need_bin ffprobe
if [[ "$is_url" -eq 1 ]]; then
  need_bin yt-dlp
fi

GROQ_KEY="$(read_groq_key)"
if [[ -z "$GROQ_KEY" ]]; then
  echo "[ERROR] GROQ_API_KEY nao encontrado no ambiente nem em ~/.aurelia/config/app.json" >&2
  exit 1
fi

tmpdir="$(mktemp -d -t aurelia-media-transcript-XXXXXX)"
trap 'rm -rf "$tmpdir"' EXIT

source_file=""
if [[ "$is_url" -eq 1 ]]; then
  yt-dlp -f bestaudio/best -o "$tmpdir/source.%(ext)s" "$INPUT" >/dev/null
  source_file="$(find "$tmpdir" -maxdepth 1 -type f -name 'source.*' | head -n 1)"
else
  source_file="$INPUT"
fi

if [[ ! -f "$source_file" ]]; then
  echo "[ERROR] arquivo de entrada nao encontrado: $source_file" >&2
  exit 1
fi

normalized="$tmpdir/normalized.wav"
ffmpeg -y -i "$source_file" -ac 1 -ar 16000 -c:a pcm_s16le "$normalized" >/dev/null 2>&1

response="$tmpdir/transcript.json"
curl -sS -X POST "https://api.groq.com/openai/v1/audio/transcriptions" \
  -H "Authorization: Bearer $GROQ_KEY" \
  -F "file=@$normalized" \
  -F "model=whisper-large-v3-turbo" \
  -F "language=pt" \
  -F "temperature=0" \
  -F "response_format=json" >"$response"

transcript="$(jq -r '.text // ""' "$response")"
if [[ -z "$transcript" ]]; then
  echo "[ERROR] resposta sem campo text" >&2
  cat "$response" >&2
  exit 1
fi

if [[ -n "$OUTPUT" ]]; then
  printf '%s\n' "$transcript" >"$OUTPUT"
else
  printf '%s\n' "$transcript"
fi
