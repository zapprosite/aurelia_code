#!/bin/bash

# Nome do binário (usa o nome da pasta se não for passado argumento)
BINARY_NAME=${1:-"aurelia-elite"}

echo "🚀 Iniciando auto-instalação e build para Ubuntu 24.04..."

# 1. Verificar se o Go está instalado
if ! command -v go &> /dev/null; then
    echo "📦 Go não encontrado. Instalando via snap..."
    sudo snap install go --classic
else
    echo "✅ Go já está instalado: $(go version)"
fi

# 2. Configurar variáveis para Linux 64-bit (Ubuntu Desktop)
export GOOS=linux
export GOARCH=amd64

# 3. Limpar cache de builds antigos para garantir um binário "puro"
echo "🧹 Limpando caches antigos..."
go clean -cache

# 4. Executar o build otimizado
# -ldflags="-s -w" remove tabelas de símbolos e debug, deixando o binário menor e mais rápido
echo "🔨 Compilando binário: $BINARY_NAME..."
# Aponta para o entrypoint correto do projeto
go build -ldflags="-s -w" -o "$BINARY_NAME" ./cmd/aurelia

# 5. Permissão de execução
chmod +x "$BINARY_NAME"

echo "------------------------------------------"
echo "✅ Sucesso! Binário gerado: ./$BINARY_NAME"
echo "💻 Plataforma alvo: $GOOS | Arch: $GOARCH"
echo "------------------------------------------"
