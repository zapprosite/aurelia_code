#!/bin/bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Aurelia Telegram CLI Send (Impersonation) - Sovereign 2026
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

TEXT="$*"
if [ -z "$TEXT" ]; then
    echo "Uso: $0 <mensagem>"
    exit 1
fi

# Endpoint de impersonação local
URL="http://localhost:8484/v1/telegram/impersonate"

# Payload JSON (usa fallbacks automáticos se userID/chatID forem omitidos no backend)
PAYLOAD=$(jq -n --arg text "$TEXT" '{text: $text}')

echo "🛰️ Enviando comando via impersonação: \"$TEXT\""

curl -s -X POST "$URL" \
     -H "Content-Type: application/json" \
     -d "$PAYLOAD" | jq .

echo -e "\n✅ Comando injetado no pipeline do Telegram."
