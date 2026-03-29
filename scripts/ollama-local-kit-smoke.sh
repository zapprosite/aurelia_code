#!/usr/bin/env bash
set -euo pipefail

MODEL="${MODEL:-qwen3.5}"
PROMPT="${PROMPT:-Responda apenas OK.}"

export OLLAMA_NUM_PARALLEL="${OLLAMA_NUM_PARALLEL:-1}"
export OLLAMA_FLASH_ATTENTION="${OLLAMA_FLASH_ATTENTION:-1}"
export OLLAMA_KV_CACHE_TYPE="${OLLAMA_KV_CACHE_TYPE:-q4_0}"
export OLLAMA_CONTEXT_LENGTH="${OLLAMA_CONTEXT_LENGTH:-8192}"

echo "🔎 Validando modelo local: $MODEL"
echo "⚙️  Política ativa: parallel=$OLLAMA_NUM_PARALLEL flash=$OLLAMA_FLASH_ATTENTION kv=$OLLAMA_KV_CACHE_TYPE ctx=$OLLAMA_CONTEXT_LENGTH"

payload=$(jq -nc \
  --arg model "$MODEL" \
  --arg prompt "$PROMPT" \
  --argjson num_ctx "$OLLAMA_CONTEXT_LENGTH" \
  '{model:$model,prompt:$prompt,stream:false,options:{num_ctx:$num_ctx}}')

response=$(curl -fsS http://127.0.0.1:11434/api/generate \
  -H "Content-Type: application/json" \
  -d "$payload")

answer=$(printf '%s' "$response" | jq -r '.response // empty')

if [[ -z "$answer" ]]; then
  echo "❌ Resposta vazia do Ollama"
  exit 1
fi

echo "✅ Smoke do kit local passou"
printf '%s\n' "$answer"
