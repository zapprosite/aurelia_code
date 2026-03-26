#!/usr/bin/env bash
set -euo pipefail

echo "🤖 Atualizando kit local do Ollama..."

MAIN_MODEL="${MAIN_MODEL:-gemma3:12b}"
EMBED_MODEL="${EMBED_MODEL:-nomic-embed-text}"

echo "📥 Puxando modelo principal: $MAIN_MODEL"
ollama pull "$MAIN_MODEL"

echo "📥 Puxando modelo de embeddings: $EMBED_MODEL"
ollama pull "$EMBED_MODEL"

echo "--------------------------------"
echo "✅ Kit Ollama atualizado"
echo "Política recomendada:"
echo "  OLLAMA_NUM_PARALLEL=1"
echo "  OLLAMA_FLASH_ATTENTION=1"
echo "  OLLAMA_KV_CACHE_TYPE=q4_0"
echo "  OLLAMA_CONTEXT_LENGTH=8192"
echo "  Residente padrão: $MAIN_MODEL"
echo "  Embedding padrão: $EMBED_MODEL"
ollama list
