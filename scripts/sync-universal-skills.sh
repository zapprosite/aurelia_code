#!/usr/bin/env bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Universal Skill Sync SOTA 2026.1 — Sovereign Footprint Orchestrator
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

# Paths
export PROJECT_ROOT="/home/will/aurelia"
export SOURCE_DIR="$PROJECT_ROOT/.agent/skills"
export CLAUDE_DIR="$PROJECT_ROOT/.claude/skills"
export BIBLIOTHECA_DIR="$PROJECT_ROOT/homelab-bibliotheca/skills/aurelia"
export DB_PATH="/home/will/.aurelia/data/aurelia.db"

echo "🚀 Iniciando Sincronização Universal de Skills..."

# 1. Garantir diretórios
mkdir -p "$CLAUDE_DIR"
mkdir -p "$BIBLIOTHECA_DIR"

# 2. Espelhamento Soberano (Rsync/CP)
echo "📦 Espelhando habilidades para Claude e Obsidian..."
cp -r "$SOURCE_DIR/"* "$CLAUDE_DIR/"
cp -r "$SOURCE_DIR/"* "$BIBLIOTHECA_DIR/"

# 3. Persistência de Metadados (SQLite)
echo "💾 Atualizando metadados em $DB_PATH..."
sqlite3 "$DB_PATH" <<EOF
CREATE TABLE IF NOT EXISTS skills_meta (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE,
    path TEXT,
    last_sync TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    source TEXT DEFAULT '.agent'
);
EOF

# Inserir/Atualizar skills no DB
find "$SOURCE_DIR" -maxdepth 1 -type d | while read -r dir; do
    SKILL_NAME=$(basename "$dir")
    if [ "$SKILL_NAME" != "skills" ] && [ "$SKILL_NAME" != "." ]; then
        sqlite3 "$DB_PATH" "INSERT OR REPLACE INTO skills_meta (name, path, source) VALUES ('$SKILL_NAME', '$dir', '.agent');"
    fi
done

# 4. Atualizar Manifesto da Bibliotheca
if [ -f "$PROJECT_ROOT/homelab-bibliotheca/lib/skills.sh" ]; then
    echo "📜 Regenerando manifesto de habilidades..."
    bash "$PROJECT_ROOT/homelab-bibliotheca/lib/skills.sh" manifest
fi

echo "✅ Sincronização Universal concluída em $(date)."
