#!/bin/bash
# master-skill-init.sh - Sovereign Bootstrap SOTA 2026
set -e

REPO_PATH=$(pwd)
echo "🚀 Iniciando Bootstrap SOTA 2026 em: $REPO_PATH"

# 1. Estrutura Contextual
mkdir -p .context/{docs,agents,plans,workflow}
touch .context/docs/README.md

# 2. Estrutura de Agentes (Regras de Elite)
mkdir -p .agent/{rules,skills,workflows}

# 3. Copiar Contratos Soberanos (do repo atual se disponível)
if [ -f "/home/will/aurelia/AGENTS.md" ]; then
    cp /home/will/aurelia/AGENTS.md .
    cp /home/will/aurelia/.env.example .env.example
    echo "✅ Contratos Soberanos (AGENTS.md) instalados."
fi

# 4. Instalação de Kits de IA (Elite 2026)
# Supondo que os kits residam no diretório canônico da Aurelia
AURELIA_DIR="/home/will/aurelia"

if [ -d "$AURELIA_DIR/.agent/rules" ]; then
    ln -sf "$AURELIA_DIR/.agent/rules/"* ".agent/rules/"
    echo "🔗 Kit Antigravity/BMad (Regras) linkado via Sovereign-Link."
fi

# 5. Inicializar Go (se necessário)
if [ ! -f "go.mod" ]; then
    go mod init "$(basename "$REPO_PATH")" || true
    echo "✅ Go Module inicializado."
fi

# 5. Criar ADR de Baseline
mkdir -p docs/adr/executed
cat <<EOF > docs/adr/executed/20260327-bootstrap-baseline.md
# ADR: Bootstrap Baseline SOTA 2026
Status: Executado
Data: $(date +%Y-%m-%d)
Contexto: Projeto inicializado via Master Skill Global.
EOF

echo "✨ Bootstrap Concluído. O projeto agora é Soberano SOTA 2026."
