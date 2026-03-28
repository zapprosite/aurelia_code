#!/bin/bash

# 💎 Aurelia Serene Audio: Perfection Script (SOTA 2026.1)
# Aplica os ajustes de fonetização, padding e chunking.

APP_JSON="/home/will/.aurelia/config/app.json"
FACTORY_GO="/home/will/aurelia/pkg/tts/factory.go"
OPENAI_GO="/home/will/aurelia/pkg/tts/openai_compatible.go"

echo "🎯 Aplicando Protocolo de Perfeição Portuguesa..."

# 1. Ajuste de Chunking no factory.go (Max 1200 chars)
sed -i "s/maxChars\s*=\s*[0-9]*/maxChars = 1200/g" "$FACTORY_GO"

# 2. Ajuste de Padding no openai_compatible.go
sed -i 's/" \. \. \. \. \."/" . . . . . "/g' "$OPENAI_GO"

# 3. Verificação de lang_code no app.json
sed -i 's/"tts_language":\s*"[^"]*"/"tts_language": "pt-br"/g' "$APP_JSON"

echo "✅ Protocolo aplicado com sucesso."
echo "🚀 Recompilando o sistema..."
go build -v -o bin/aurelia ./cmd/aurelia
sudo systemctl restart aurelia.service
echo "✨ Aurélia agora fala com perfeição SOTA 2026."
