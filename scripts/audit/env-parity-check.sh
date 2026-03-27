#!/bin/bash
# ==============================================================================
# AURÉLIA — SOBERANA 2026 (ENV PARITY AUDITOR)
# ==============================================================================
# Verifica se o .env local é um espelho estrutural do .env.example.
# ==============================================================================

ENV_FILE=".env"
EXAMPLE_FILE=".env.example"

if [ ! -f "$ENV_FILE" ]; then
    echo "❌ Erro: Arquivo $ENV_FILE não encontrado!"
    exit 1
fi

if [ ! -f "$EXAMPLE_FILE" ]; then
    echo "❌ Erro: Arquivo $EXAMPLE_FILE não encontrado!"
    exit 1
fi

# Extrair chaves (ignorando comentários e linhas vazias)
keys_example=$(grep -E '^[A-Z0-9_]+=' "$EXAMPLE_FILE" | cut -d'=' -f1)
keys_env=$(grep -E '^[A-Z0-9_]+=' "$ENV_FILE" | cut -d'=' -f1)

missing_keys=""
for key in $keys_example; do
    if ! echo "$keys_env" | grep -q "^$key$"; then
        missing_keys="$missing_keys $key"
    fi
done

if [ -n "$missing_keys" ]; then
    echo "⚠️  AVISO de Desigualdade Detectado! (Soberania 2026)"
    echo "As seguintes chaves estão no .env.example mas NÃO no .env:"
    for key in $missing_keys; do
        echo "  - $key"
    done
    echo "----------------------------------------------------------------"
    echo "Ação necessária: Adicione estas chaves ao seu .env para manter o espelho."
    exit 1
else
    echo "✅ Paridade de Ambiente Confirmada. (Mirror-Perfect 100%)"
    exit 0
fi
