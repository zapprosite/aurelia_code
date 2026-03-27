#!/usr/bin/env bash
# scripts/launcher.sh
# Launcher wrapper that sources secrets.env before starting aurelia app
# Usage: bash scripts/launcher.sh [app-args...]

set -euo pipefail

CONFIG_DIR="$HOME/.aurelia/config"
SECRETS_FILE="$CONFIG_DIR/secrets.env"

# Load secrets into environment
if [[ ! -f "$SECRETS_FILE" ]]; then
    echo "❌ Error: secrets.env not found at $SECRETS_FILE"
    exit 1
fi

# Source secrets as environment variables
set -a
# shellcheck disable=SC1090
source "$SECRETS_FILE"
set +a

echo "✅ Secrets loaded from environment"
echo "🚀 Starting Aurelia app with config: $CONFIG_DIR/app.json"

# Pass all remaining args to main app
# This is a template — replace 'aurelia-app' with actual app command
# exec aurelia-app "$@"
