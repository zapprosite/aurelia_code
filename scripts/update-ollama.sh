#!/usr/bin/env bash
set -euo pipefail

echo "🤖 Atualizando kit local do Ollama..."

MAIN_MODEL="${MAIN_MODEL:-gemma3:12b}"
LIGHT_MODEL="${LIGHT_MODEL:-gemma3:4b}"
EMBED_MODEL="${EMBED_MODEL:-bge-m3:latest}"
OPTIONAL_MODEL="${OPTIONAL_MODEL:-gemma3:27b-it-q4_K_M}"

echo "📥 Puxando modelo principal: $MAIN_MODEL"
ollama pull "$MAIN_MODEL"

echo "📥 Puxando modelo leve: $LIGHT_MODEL"
ollama pull "$LIGHT_MODEL"

echo "📥 Puxando modelo de embeddings: $EMBED_MODEL"
ollama pull "$EMBED_MODEL"

if [[ -n "$OPTIONAL_MODEL" ]]; then
  echo "📥 Puxando modelo opcional: $OPTIONAL_MODEL"
  ollama pull "$OPTIONAL_MODEL"
fi

echo "--------------------------------"
echo "✅ Kit Ollama atualizado"
echo "Política recomendada:"
echo "  OLLAMA_NUM_PARALLEL=1"
echo "  OLLAMA_FLASH_ATTENTION=1"
echo "  OLLAMA_KV_CACHE_TYPE=q4_0"
echo "  OLLAMA_CONTEXT_LENGTH=8192"
echo "  Residente padrão: $MAIN_MODEL"
echo "  Sob demanda: $LIGHT_MODEL"
ollama list
