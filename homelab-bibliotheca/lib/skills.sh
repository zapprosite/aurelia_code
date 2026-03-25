#!/usr/bin/env bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Sovereign-Bibliotheca v2 — Gestor de Skills (OpenClaw Indexer)
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

source "$(dirname "${BASH_SOURCE[0]}")/config.sh"

function usage() {
    echo -e "Uso: $0 [comando]"
    echo -e "  list              Lista todas as skills importadas"
    echo -e "  show <nome>       Mostra detalhes da skill (SKILL.md)"
    echo -e "  manifest          Gera manifest.json de todas as skills"
}

function list_skills() {
    echo -e "${CLR_INFO}=== Skills OpenClaw ===${CLR_RESET}"
    find "$SKILLS_DIR/skills" -maxdepth 3 -name "SKILL.md" -type f | while read -r line; do
        SKILL_NAME=$(basename "$(dirname "$line")")
        echo -e "  - $SKILL_NAME"
    done
}

function show_skill() {
    NAME="$1"
    FILE=$(find "$SKILLS_DIR/skills" -maxdepth 3 -name "$NAME" -type d)/SKILL.md
    if [ -f "$FILE" ]; then
        cat "$FILE"
    else
        echo -e "${CLR_ERROR}Skill '$NAME' não encontrada.${CLR_RESET}"
    fi
}

function generate_manifest() {
    echo -e "${CLR_INFO}Gerando manifesto de skills...${CLR_RESET}"
    MANIFEST_FILE="$BIBLIOTHECA_ROOT/skills-registry.json"
    
    echo "[" > "$MANIFEST_FILE"
    FIRST=true
    # Otimizado: evita grep em binários e usa head para velocidade
    find "$SKILLS_DIR/skills" -maxdepth 3 -name "SKILL.md" -type f | while read -r line; do
        if [ "$FIRST" = true ]; then FIRST=false; else echo "," >> "$MANIFEST_FILE"; fi
        
        NAME=$(basename "$(dirname "$line")")
        PATH_REL=${line#$PROJECT_ROOT/}
        
        # Extrair descrição de forma segura via head (evita grep pesado em binários)
        DESC=$(head -n 20 "$line" | grep -v "^#" | grep -v "^$" | head -n 1 | tr -d '"' | tr -d '\n' | tr -d '\r' | cut -c 1-100)
        
        printf '  {"name": "%s", "path": "%s", "description": "%s"}' "$NAME" "$PATH_REL" "$DESC" >> "$MANIFEST_FILE"
    done
    echo -e "\n]" >> "$MANIFEST_FILE"
    
    echo -e "${CLR_SUCCESS}Manifesto gerado em $MANIFEST_FILE${CLR_RESET}"
}

case $1 in
    list)     list_skills ;;
    show)     show_skill "$2" ;;
    manifest) generate_manifest ;;
    *)        usage ;;
esac
