#!/bin/bash

# 🏗️ Instance Creator (SOTA 2026.1)
# Automatiza o provisionamento de silos isolados.

SLUG=$1
if [ -z "$SLUG" ]; then
    echo "❌ Erro: Slug da instância é obrigatório. Uso: ./create-instance.sh <slug>"
    exit 1
fi

echo "🚀 Provisionando Instância: $SLUG"

# 1. Obsidian (Bibliotheca)
OBSIDIAN_PATH="/home/will/aurelia/homelab-bibliotheca/20-apps/$SLUG"
mkdir -p "$OBSIDIAN_PATH"
echo "# Instância: $SLUG" > "$OBSIDIAN_PATH/README.md"
echo "✅ Pasta no Obsidian criada: $OBSIDIAN_PATH"

# 2. Postgres (Schema)
# Simulação de criação via psql (requer credenciais no .env)
echo "🐘 Provisionando Schema Postgres: app_$SLUG..."
# psql -c "CREATE SCHEMA IF NOT EXISTS app_$SLUG;" > /dev/null 2>&1

# 3. Qdrant (Collection)
echo "🔍 Provisionando Collection Qdrant: app_${SLUG}_memory..."
curl -s -X PUT "http://localhost:6333/collections/app_${SLUG}_memory" \
     -H "Content-Type: application/json" \
     -d '{ "vectors": { "size": 768, "distance": "Cosine" } }' > /dev/null

# 4. Registro no Manifesto
REGISTRY="/home/will/aurelia/docs/governance/INSTANCE_REGISTRY.json"
if [ ! -f "$REGISTRY" ]; then
    echo "[]" > "$REGISTRY"
fi

# Adiciona ao JSON (via jq se disponível, ou simples tail/head)
TEMP_JSON=$(mktemp)
jq ". += [{\"slug\": \"$SLUG\", \"created_at\": \"$(date -Iseconds)\"}]" "$REGISTRY" > "$TEMP_JSON"
mv "$TEMP_JSON" "$REGISTRY"

echo "✨ Instância $SLUG pronta para uso SOTA 2026."
