#!/usr/bin/env bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Sovereign-Bibliotheca v2 — Motor de Comunicação (Telegram Multi-Bot)
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

source "$(dirname "${BASH_SOURCE[0]}")/config.sh"

function usage() {
    echo -e "Uso: $0 [comando]"
    echo -e "  send <mensagem> [bot_id]     Envia mensagem via bot específico"
    echo -e "  impersonate <texto>          Fala como a Aurélia (interno)"
    echo -e "  list                         Lista bots configurados"
}

function send_msg() {
    MSG="$1"; BOT_ID="${2:-aurelia}"
    echo -e "${CLR_INFO}Enviando via bot '$BOT_ID'...${CLR_RESET}"
    
    # Integração com o bot-cli.sh existente
    "$PROJECT_ROOT/scripts/bot-cli.sh" ping "$MSG"
}

function list_bots() {
    "$PROJECT_ROOT/scripts/bot-cli.sh" list
}

case $1 in
    send) send_msg "$2" "$3" ;;
    impersonate) send_msg "Impersonate: $2" "aurelia" ;;
    list) list_bots ;;
    *) usage ;;
esac
