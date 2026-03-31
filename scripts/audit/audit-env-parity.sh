#!/usr/bin/env bash
# scripts/audit/audit-env-parity.sh
# 🔍 Auditoria de Paridade de Ambiente — SOTA 2026

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DOTENV="$REPO_ROOT/.env"
EXAMPLE="$REPO_ROOT/.env.example"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔍 Verificando paridade de ambiente...${NC}"

if [ ! -f "$DOTENV" ]; then
    echo -e "${RED}❌ Arquivo .env não encontrado em $DOTENV${NC}"
    exit 1
fi

if [ ! -f "$EXAMPLE" ]; then
    echo -e "${YELLOW}⚠️  .env.example não encontrado. Criando base a partir do .env...${NC}"
    grep -v '^#' "$DOTENV" | grep '=' | cut -d'=' -f1 | sed 's/$/=/' > "$EXAMPLE"
    echo -e "${GREEN}✅ .env.example gerado.${NC}"
    exit 0
fi

# Extrair chaves de ambos os arquivos
KEYS_EXAMPLE=$(grep -v '^#' "$EXAMPLE" | grep '=' | cut -d'=' -f1 | sort)
KEYS_DOTENV=$(grep -v '^#' "$DOTENV" | grep '=' | cut -d'=' -f1 | sort)

MISSING_IN_ENV=()
for key in $KEYS_EXAMPLE; do
    if ! grep -q "^${key}=" "$DOTENV"; then
        MISSING_IN_ENV+=("$key")
    fi
done

if [ ${#MISSING_IN_ENV[@]} -eq 0 ]; then
    echo -e "${GREEN}✅ Paridade OK: Todas as chaves do .env.example estão no .env${NC}"
    exit 0
else
    echo -e "${RED}❌ Paridade FALHOU: Chaves faltando no .env:${NC}"
    for key in "${MISSING_IN_ENV[@]}"; do
        echo -e "  - $key"
    done
    exit 1
fi
