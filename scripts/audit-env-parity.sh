#!/bin/bash
# [Aurelia Industrial SOTA 2026]
# Script para auditar a paridade entre .env e .env.example.

ENV_FILE="/home/will/aurelia/.env"
EXAMPLE_FILE="/home/will/aurelia/.env.example"

echo "🔍 Auditando paridade .env <-> .env.example..."

KEYS_ENV=$(grep -oP '^[^#\s][^=]*' "$ENV_FILE" | sort)
KEYS_EXAMPLE=$(grep -oP '^[^#\s][^=]*' "$EXAMPLE_FILE" | sort)

DIFF=$(diff <(echo "$KEYS_ENV") <(echo "$KEYS_EXAMPLE"))

if [ -z "$DIFF" ]; then
    echo "✅ Paridade OK! Estruturas sincronizadas."
    exit 0
else
    echo "❌ DISPARIDADE DETECTADA!"
    echo "$DIFF"
    echo ""
    echo "💡 Sugestão: Execute o comando abaixo para sincronizar estruturalmente:"
    echo "perl -ne 'if (/^([^#\s][^=]*)=/) { print \"\$1=\\n\" } else { print }' $ENV_FILE > $EXAMPLE_FILE"
    exit 1
fi
