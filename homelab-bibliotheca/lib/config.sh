#!/usr/bin/env bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Sovereign-Bibliotheca v2 — Configuração Centralizada
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

# Caminhos Base
export LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export BIBLIOTHECA_ROOT="$(dirname "$LIB_DIR")"
export PROJECT_ROOT="$(dirname "$BIBLIOTHECA_ROOT")"
export SKILLS_DIR="$BIBLIOTHECA_ROOT/skills/open-claw"
export AURELIA_SKILLS_DIR="$BIBLIOTHECA_ROOT/skills/aurelia"

# Configurações de API (Local first)
export QDRANT_URL="http://localhost:6333"
export SUPABASE_REST_URL="http://localhost:3000"
export SUPABASE_STUDIO_URL="http://localhost:54323"
export AURELIA_API_URL="http://localhost:3334"
export TELEGRAM_IMPERSONATE_URL="http://localhost:8484/v1/telegram/impersonate"

# Caminhos de Dados
export AURELIA_DB_PATH="$HOME/.aurelia/data/aurelia.db"
export SECRETS_FILE="$HOME/.aurelia/config/secrets.env"

# Tentativa de carregar segredos
if [ -f "$SECRETS_FILE" ]; then
    source "$SECRETS_FILE"
fi

# Fallback para Obsidian Vault (Master, favor ajustar se necessário)
export OBSIDIAN_VAULT_PATH="${OBSIDIAN_VAULT_PATH:-$HOME/Documents/ObsidianVault}"

# Cores para Output
export CLR_INFO="\033[0;34m"
export CLR_SUCCESS="\033[0;32m"
export CLR_WARN="\033[1;33m"
export CLR_ERROR="\033[0;31m"
export CLR_RESET="\033[0m"
