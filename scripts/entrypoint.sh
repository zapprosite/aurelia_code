#!/bin/bash

# ==============================================================================
# 🛰️ AURÉLIA — SETUP INICIAL AMIGÁVEL
# ==============================================================================

set -e

# Cores para o terminal
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}"
echo "  /$$$$$$                                /$$ /$$"
echo " /$$__  $$                              | $$|__/"
echo "| $$  \ $$ /$$   /$$  /$$$$$$   /$$$$$$ | $$ /$$  /$$$$$$"
echo "| $$$$$$$$| $$  | $$ /$$__  $$ /$$__  $$| $$| $$ /$$__  $$"
echo "| $$__  $$| $$  | $$| $$  \__/| $$$$$$$$| $$| $$| $$  \ $$"
echo "| $$  | $$| $$  | $$| $$      | $$_____/| $$| $$| $$  | $$"
echo "| $$  | $$|  $$$$$$/| $$      |  $$$$$$$| $$| $$|  $$$$$$/"
echo "|__/  |__/ \______/ |__/       \_______/|__/|__/ \______/ "
echo -e "${NC}"
echo -e "${GREEN}Bem-vindo ao Setup da Aurélia — Sua IA Soberana (2026)${NC}"
echo "----------------------------------------------------------"

# 1. Verificação de Dependências
echo -e "\n${BLUE}[1/4] Verificando ferramentas...${NC}"

check_cmd() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}Erro: $1 não está instalado.${NC}"
        return 1
    fi
    echo -e "${GREEN}OK: $1 encontrado.${NC}"
}

check_cmd docker || { echo "Por favor, instale o Docker: https://docs.docker.com/get-docker/"; exit 1; }
check_cmd go || { echo "Por favor, instale o Go: https://go.dev/doc/install"; exit 1; }

# 2. Escolha de Modo
echo -e "\n${BLUE}[2/4] Escolha seu modo de operação:${NC}"
echo "1) 🛰️  Soberano (Hardware Pesado - GPU NVIDIA necessária)"
echo "2) ☁️  Lite (Hardware Comum - Usa APIs na nuvem)"
read -p "Opção (1 ou 2): " MODE_OPT

if [ "$MODE_OPT" == "1" ]; then
    A_MODE="sovereign"
    echo -e "${YELLOW}Modo Soberano selecionado. Prepare sua GPU!${NC}"
else
    A_MODE="lite"
    echo -e "${YELLOW}Modo Lite selecionado. Economizando seus recursos locais.${NC}"
fi

# 3. Configuração de Ambiente (.env)
echo -e "\n${BLUE}[3/4] Configuração de chaves...${NC}"

if [ ! -f .env ]; then
    cp .env.example .env
    echo ".env criado a partir do template."
fi

# Função para atualizar .env
update_env() {
    local key=$1
    local value=$2
    if grep -q "^$key=" .env; then
        sed -i "s|^$key=.*|$key=$value|" .env
    else
        echo "$key=$value" >> .env
    fi
}

update_env "AURELIA_MODE" "$A_MODE"

if grep -q "^TELEGRAM_BOT_TOKEN=$" .env || [ -z "$(grep "^TELEGRAM_BOT_TOKEN=" .env | cut -d'=' -f2)" ]; then
    echo -e "${YELLOW}Você precisa de um Token do Telegram.${NC}"
    echo "Crie um bot aqui: https://t.me/BotFather"
    read -p "Cole seu Token aqui: " TG_TOKEN
    update_env "TELEGRAM_BOT_TOKEN" "$TG_TOKEN"
fi

if [ "$A_MODE" == "lite" ]; then
    echo -e "\nNo modo Lite, recomendamos usar OpenRouter ou Google Gemini."
    read -p "API Key do OpenRouter (opcional): " OR_KEY
    if [ ! -z "$OR_KEY" ]; then update_env "OPENROUTER_API_KEY" "$OR_KEY"; fi
fi

# 4. Inicialização
echo -e "\n${BLUE}[4/4] Iniciando o ecossistema...${NC}"

if [ "$A_MODE" == "sovereign" ]; then
    docker-compose up -d
else
    # No modo lite, subimos apenas o necessário (Postgres, Qdrant)
    docker-compose up -d postgres qdrant
fi

echo -e "\n${GREEN}🚀 TUDO PRONTO!${NC}"
echo "----------------------------------------------------------"
echo "O Dashboard está iniciando em: http://localhost:3334"
echo "Para ler o guia de iniciantes, acesse: docs/BEM-VINDO.md"
echo -e "\n${YELLOW}Dica Senior: Use 'go run ./cmd/aurelia' para rodar o core agora.${NC}"
