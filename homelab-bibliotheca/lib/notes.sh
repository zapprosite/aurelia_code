#!/usr/bin/env bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Sovereign-Bibliotheca v2 — Motor de Notas (Obsidian & Markdown)
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

source "$(dirname "${BASH_SOURCE[0]}")/config.sh"

function usage() {
    echo -e "Uso: $0 [comando]"
    echo -e "  create <título> [conteúdo]   Cria uma nova nota"
    echo -e "  append <título> <texto>      Adiciona texto a uma nota"
    echo -e "  list                         Lista notas no vault/pasta"
}

function get_note_path() {
    TITLE="$1"
    # Fallback: Se o vault não existir, usa a pasta de knowledge do monorepo
    if [ -d "$OBSIDIAN_VAULT_PATH" ]; then
        echo "$OBSIDIAN_VAULT_PATH/$TITLE.md"
    else
        echo "$PROJECT_ROOT/knowledge/$TITLE.md"
    fi
}

function create_note() {
    TITLE="$1"; CONTENT="$2"
    PATH_NOTE=$(get_note_path "$TITLE")
    mkdir -p "$(dirname "$PATH_NOTE")"

    # Uso do obsidian-cli se disponível, caso contrário fallback nativo
    if command -v obsidian-cli &> /dev/null; then
        obsidian-cli create -t "$TITLE" -c "$CONTENT"
    else
        echo -e "---\ncreated_at: $(date --iso-8601=seconds)\nsource: sovereign-lib\n---\n\n$CONTENT" > "$PATH_NOTE"
    fi
    echo -e "${CLR_SUCCESS}Nota '$TITLE' criada em $PATH_NOTE${CLR_RESET}"
}

function append_note() {
    TITLE="$1"; TEXT="$2"
    PATH_NOTE=$(get_note_path "$TITLE")
    
    if [ ! -f "$PATH_NOTE" ]; then
        create_note "$TITLE" "$TEXT"
    else
        echo -e "\n\n### Update: $(date --iso-8601=seconds)\n$TEXT" >> "$PATH_NOTE"
        echo -e "${CLR_SUCCESS}Conteúdo adicionado a '$TITLE'${CLR_RESET}"
    fi
}

case $1 in
    create) create_note "$2" "$3" ;;
    append) append_note "$2" "$3" ;;
    list) ls -1 "$(dirname "$(get_note_path "test")")"/*.md 2>/dev/null ;;
    *) usage ;;
esac
