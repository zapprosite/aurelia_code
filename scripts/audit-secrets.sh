#!/bin/bash

# ==============================================================================
# 🛰️  SOVEREIGN SECRET AUDITOR (INDUSTRIAL v2026)
# Standard: Zero Hardcode & Homelab Governance
# Developer: Antigravity / Aurelia_Code
# ==============================================================================

# --- Configuration ---
BOLD="\033[1m"
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
CYAN="\033[0;36m"
NC="\033[0m"

EXIT_CODE=0
FINDINGS_COUNT=0

# Exclude directories that are safe or known to contain metadata matching patterns
EXCLUDE_DIRS=(
    ".git"
    "node_modules"
    "vendor"
    ".aurelia"
    ".context"
    "docs"
    "build"
    "homelab-bibliotheca"
    "testdata"
    "assets"
)

# Exclude files that are safe or intended to have long strings (configs/registries)
EXCLUDE_FILES=(
    ".env"
    ".env.example"
    "audit-secrets.sh"
    "*.md"
    "*.json"
    "cost.go"      # Registros de custos de modelos
    "resolver.go"  # Definições de caminhos estruturais
    "*_test.go"    # Arquivos de teste
)

# --- Functions ---

function log_step() {
    echo -e "\n${BOLD}${CYAN}» $1${NC}"
}

function check_parity() {
    log_step "Auditando Paridade de Variáveis (.env vs .env.example)"
    
    if [ ! -f ".env" ] || [ ! -f ".env.example" ]; then
        echo -e "  [${YELLOW}SKIP${NC}] Arquivos de ambiente incompletos para auditoria de paridade."
        return
    fi

    # Extrair chaves (sanitizado)
    keys_env=$(grep -vE '^#|^$' .env | cut -d'=' -f1 | tr -d ' ' | sort | uniq)
    keys_example=$(grep -vE '^#|^$' .env.example | cut -d'=' -f1 | tr -d ' ' | sort | uniq)
    
    missing=$(comm -23 <(echo "$keys_example") <(echo "$keys_env"))
    
    if [ -n "$missing" ]; then
        echo -e "  [${RED}FAIL${NC}] Variáveis documentadas, mas ausentes no .env local:"
        echo "$missing" | sed 's/^/    - /'
        FINDINGS_COUNT=$((FINDINGS_COUNT + 1))
        EXIT_CODE=1
    else
        echo -e "  [${GREEN}OK${NC}] Paridade de ambiente em conformidade absoluta."
    fi
}

function check_hardcode() {
    log_step "Escaneando Segredos Hardcoded (Zero Hardcode Policy)"
    
    # Padrão: Atribuições agressivas que parecem strings constantes
    # Foca em sk-... (OpenAI), Bearer ..., ou chaves alfanuméricas longas
    grep_patterns="api[_-]?key|token|secret|password|bearer|sk-[a-zA-Z0-9]{15,}"
    
    # Construir exclusões para o grep/rg
    exclude_cmd=""
    for dir in "${EXCLUDE_DIRS[@]}"; do exclude_cmd="$exclude_cmd --exclude-dir=$dir"; done
    for file in "${EXCLUDE_FILES[@]}"; do exclude_cmd="$exclude_cmd --exclude=$file"; done

    # Lógica de detecção: 
    # 1. Busca por patterns em fontes.
    # 2. Filtra por atribuições ( : ou = ).
    # 3. Filtra por strings entre aspas de min 15 chars (evita nomes curtos).
    # 4. Exclui metadados de sistema (json tags, IDs, caminhos).
    leaks=$(grep -rEi "$grep_patterns" . $exclude_cmd | \
            grep -E "[:=]\s*['\"][a-zA-Z0-9/_-]{15,}['\"]" | \
            grep -v "{chave-para-env}" | \
            grep -vE "json:\"|ID:|\"path\":|test|mock|example|placeholder|config\.|cfg\." | \
            head -n 10)

    if [ -n "$leaks" ]; then
        echo -e "  [${RED}FAIL${NC}] Possíveis vazamentos detectados (Strings constantes):"
        echo "$leaks" | sed 's/^/    | /'
        FINDINGS_COUNT=$((FINDINGS_COUNT + 1))
        EXIT_CODE=1
    else
        echo -e "  [${GREEN}OK${NC}] Nenhum segredo óbvio detectado nos arquivos do repositório."
    fi
}

function check_logs() {
    log_step "Auditando Logs de Sistema (Homelab Control)"
    
    AURELIA_LOGS="$HOME/.aurelia/logs"
    if [ ! -d "$AURELIA_LOGS" ]; then
        echo -e "  [${YELLOW}SKIP${NC}] Diretório de logs globais ($AURELIA_LOGS) inacessível."
        return
    fi

    # Scan por tokens reais em logs (agora é silencioso se OK)
    log_leaks=$(grep -rhEi "Bearer sk-|sk-[a-zA-Z0-9]{20,}" "$AURELIA_LOGS" 2>/dev/null | head -n 3)
    
    if [ -n "$log_leaks" ]; then
        echo -e "  [${RED}WARN${NC}] Tokens reais vazados nos logs do daemon!"
        echo "$log_leaks" | sed 's/^/    | /'
        echo -e "  ${YELLOW}Sugestão: Execute '> $AURELIA_LOGS/daemon.log' para limpar.${NC}"
        # Não bloqueia o push por logs de máquina local, apenas alerta.
    else
        echo -e "  [${GREEN}OK${NC}] Logs locais sanitizados."
    fi
}

function check_armor() {
    log_step "Validando Blindagem de Repositório (.gitignore)"
    
    if grep -q "^\.env$" .gitignore; then
        echo -e "  [${GREEN}OK${NC}] .env está devidamente blindado."
    else
        echo -e "  [${RED}CRITICAL${NC}] .env EXPOSTO! Blindando automaticamente..."
        echo ".env" >> .gitignore
    fi
}

# --- Execution ---

header="🛰️  SOVEREIGN SCAN — $(date +'%Y-%m-%d %H:%M:%S')"
echo -e "${BOLD}${CYAN}$header${NC}"
echo -e "${CYAN}$(printf '%.s=' $(seq 1 ${#header}))${NC}"

check_parity
check_hardcode
check_logs
check_armor

echo -e "\n${BOLD}${CYAN}$(printf '%.s=' $(seq 1 ${#header}))${NC}"
if [ "$EXIT_CODE" -eq 0 ]; then
    echo -e "${GREEN}${BOLD}✅ AUDIT PASSED: Ecossistema em Conformidade Soberana.${NC}"
    exit 0
else
    echo -e "${RED}${BOLD}❌ AUDIT FAILED: $FINDINGS_COUNT problema(s) detectado(s).${NC}"
    echo -e "${RED}Corrija as falhas antes de prosseguir com o push.${NC}"
    exit 1
fi
