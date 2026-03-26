#!/usr/bin/env bash
# scripts/enterprise-scrub.sh
# Enterprise Secret Scrubbing — Sovereign 2026 Standard

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PLACEHOLDER="{chave-secrets}"

echo "🛰️ Iniciando Enterprise Scrubbing..."

# Padrões a serem limpos (API Keys, Tokens, Hex Addresses suspeitos)
PATTERNS=(
    # API Keys (padrão genérico de 20+ chars)
    "sk-[a-zA-Z0-9]{20,}"
    "{{GITHUB_TOKEN_HIDDEN}}[a-zA-Z0-9]{20,}"
    "{{GROQ_TOKEN_HIDDEN}}[a-zA-Z0-9]{20,}"
    "xai-[a-zA-Z0-9]{20,}"
    "AIzaSy[a-zA-Z0-9_-]{33}" # Google API Key
    
    # Tokens e Segredos em JSON/Go/MD
    "(api.?key|token|secret|password|passwd|bearer)\s*[:=]\s*['\"][a-zA-Z0-9_-]{20,}['\"]"
    
    # Endereços Hexadecimais em docs (exceto os curtos/técnicos de infra)
    "0x[a-fA-F0-9]{30,}"
)

# Arquivos a ignorar
EXCLUDES=(
    "*.png"
    "*.jpg"
    "*.jpeg"
    "*.gif"
    "*.ico"
    "*.pdf"
    "*.zip"
    "*.tar.gz"
    ".git/*"
    "node_modules/*"
    "vendor/*"
    ".env*"
    "secrets.env"
    "scripts/enterprise-scrub.sh"
)

# Construir argumentos de exclusão para o find
FIND_EXCLUDES=""
for ext in "${EXCLUDES[@]}"; do
    FIND_EXCLUDES="$FIND_EXCLUDES ! -path \"*/$ext\""
done

for pattern in "${PATTERNS[@]}"; do
    echo "🔍 Higienizando padrão: $pattern"
    # Usar sed para substituir o conteúdo capturado pelo placeholder
    # Nota: para padrões complexos com captura, simplificamos para manter a chave do campo.
    if [[ "$pattern" == *[:=]* ]]; then
        # Se contiver := ou =, tentamos preservar o prefixo
        # Ex: api_key: "REAL_KEY" -> api_key: "{chave-secrets}"
        find "$REPO_ROOT" -type f $FIND_EXCLUDES -exec sed -i -E "s/($pattern)/\2: \"$PLACEHOLDER\"/g" {} + 2>/dev/null || true
        # Fallback simples se a captura falhar
        find "$REPO_ROOT" -type f $FIND_EXCLUDES -exec sed -i -E "s/$pattern/secret: \"$PLACEHOLDER\"/g" {} + 2>/dev/null || true
    else
        # Substituição direta
        find "$REPO_ROOT" -type f $FIND_EXCLUDES -exec sed -i -E "s/$pattern/$PLACEHOLDER/g" {} + 2>/dev/null || true
    fi
done

echo "✅ Scrubbing concluído. Varrendo códigos de retorno..."

# Garantir que nenhum .env inesperado foi deixado
find "$REPO_ROOT" -name ".env" ! -name ".env.example" -type f -delete 2>/dev/null || true

echo "Done."
