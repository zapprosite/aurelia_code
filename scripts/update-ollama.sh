#!/bin/bash
echo "🤖 Atualizando Modelos Ollama..."

# 1. Definir Modelos
CHAT_MODEL="qwen2.5:32b" # Nota: Ajustei para o disponível mais próximo se o 3.5 não carregar
EMBED_MODEL="mxbai-embed-large"

# 2. Pull
echo "📥 Puxando modelo de chat: $CHAT_MODEL..."
ollama pull $CHAT_MODEL

echo "📥 Puxando modelo de embeddings: $EMBED_MODEL..."
ollama pull $EMBED_MODEL

echo "--------------------------------"
echo "✅ Todos os modelos estão atualizados!"
ollama list
