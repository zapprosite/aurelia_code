#!/bin/bash
# obsidian-sync.sh — Sincroniza Regras e Skills com a Vault do Obsidian.
# SOTA 2026.1. Sênior Direto ao Ponto.

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "$REPO_ROOT/.env"

VAULT_PATH="${OBSIDIAN_VAULT_PATH:-$HOME/Documents/Obsidian/Aurelia}"
RULES_DIR="$REPO_ROOT/.agent/rules"
SKILLS_DIR="$REPO_ROOT/.agent/skills"
TARGET_RULES="$VAULT_PATH/Governance/Rules"
TARGET_SKILLS="$VAULT_PATH/Skills/Core"
TARGET_DOCS="$VAULT_PATH/Training"

echo "🔄 Sincronizando com Obsidian Vault em $VAULT_PATH..."

mkdir -p "$TARGET_RULES" "$TARGET_SKILLS" "$TARGET_DOCS"

# Sincronizar Regras e Docs de Treinamento
cp -v "$RULES_DIR"/*.md "$TARGET_RULES/"
cp -v "$REPO_ROOT/.agent/docs"/*.md "$TARGET_DOCS/"

# Sincronizar Skills (Discovery Index)
# Cria um arquivo MD para cada skill no Obsidian
for skill in $(ls -1 "$SKILLS_DIR" | grep -v "README"); do
    echo "Skill: $skill"
    cat <<EOF > "$TARGET_SKILLS/$skill.md"
---
name: $skill
type: sovereign-skill
location: $SKILLS_DIR/$skill
synced: $(date)
---
# $skill

Ver o diretório local para detalhes: [Link Local](file://$SKILLS_DIR/$skill/SKILL.md)
EOF
done

echo "✅ Sincronização Obsidian concluída."
