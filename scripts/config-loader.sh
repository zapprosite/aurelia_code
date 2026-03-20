#!/usr/bin/env bash
# scripts/config-loader.sh
# Exports environment variables from secrets.env for use by MCP/app at runtime
# Do NOT modify config files — app must read env vars at startup

set -euo pipefail

CONFIG_DIR="$HOME/.aurelia/config"
SECRETS_FILE="$CONFIG_DIR/secrets.env"

if [[ ! -f "$SECRETS_FILE" ]]; then
    echo "❌ Error: $SECRETS_FILE not found"
    exit 1
fi

# Load and export all secrets as environment variables
# Use 'set -a' to export all variables, then 'set +a' to disable
set -a
# shellcheck disable=SC1090
source "$SECRETS_FILE"
set +a

echo "✅ Secrets loaded from $SECRETS_FILE"
echo "Environment variables exported:"
echo "  - CLOUDFLARE_MCP_TOKEN"
echo "  - TELEGRAM_BOT_TOKEN"
echo "  - GOOGLE_API_KEY"
echo "  - OPENROUTER_API_KEY"
echo "  - GROQ_API_KEY"
echo "  - QDRANT_API_KEY"
echo "  - GITHUB_TOKEN"
echo "  - POSTGRES_PASSWORD"
echo ""
echo "Usage: source <(bash scripts/config-loader.sh)"
