#!/bin/bash

set -euo pipefail

# Nome do binário (usa o nome da pasta se não for passado argumento)
BINARY_NAME=${1:-"aurelia-elite"}
ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

resolve_go() {
    if command -v go >/dev/null 2>&1; then
        command -v go
        return 0
    fi
    if [ -x /snap/bin/go ]; then
        echo /snap/bin/go
        return 0
    fi
    return 1
}

echo "🚀 Iniciando build operacional do Aurelia..."

GO_BIN=$(resolve_go || true)
if [ -z "${GO_BIN:-}" ]; then
    echo "❌ ERRO OPERACIONAL: Go (Golang) não encontrado no PATH nem em /snap/bin/go."
    echo "Para rodar este runtime, instale o Go 1.25 ou superior: https://go.dev/doc/install"
    echo "Dica: Tente 'sudo snap install go --classic' no Ubuntu."
    exit 1
fi

echo "✅ Go encontrado: $("$GO_BIN" version)"

# 2. Configurar variáveis para Linux 64-bit (Ubuntu Desktop)
export GOOS=linux
export GOARCH=amd64

# 3. Limpar cache de builds antigos para garantir um binário "puro"
echo "🧹 Limpando caches antigos..."
"$GO_BIN" clean -cache

# 4. Executar o build otimizado
# -ldflags="-s -w" remove tabelas de símbolos e debug, deixando o binário menor e mais rápido
echo "🔨 Compilando binário: $BINARY_NAME..."
# Aponta para o entrypoint correto do projeto
cd "$ROOT_DIR"
"$GO_BIN" build -ldflags="-s -w" -o "$BINARY_NAME" ./cmd/aurelia

# 5. Permissão de execução
chmod +x "$BINARY_NAME"

echo "------------------------------------------"
echo "✅ Sucesso! Binário gerado: ./$BINARY_NAME"
echo "💻 Plataforma alvo: $GOOS | Arch: $GOARCH"
echo "------------------------------------------"
