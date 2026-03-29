#!/usr/bin/env bash
# [SOTA 2026] Legacy Model Scanner
# Purpose: Detect forbidden legacy model references (Gemma 3) to ensure sovereign purity.

set -euo pipefail

COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_NC='\033[0m'

echo -e "${COLOR_YELLOW}🔍 Scanning for legacy model references (Gemma 3)...${COLOR_NC}"

# Define patterns to search (DANGER: DO NOT RENAME THESE TO QWEN3.5)
LEGACY_PATTERNS=(
    "gemma3:27b"
    "gemma3:12b"
    "modelGemma3"
    "gemma3"
)

TOTAL_MATCHES=0

for PATTERN in "${LEGACY_PATTERNS[@]}"; do
    echo -e "\nChecking for: ${COLOR_YELLOW}$PATTERN${COLOR_NC}"
    # Search in code and scripts, excluding common ignore directories
    # Note: excluding docs/adr as they are historical records, but including active docs/governance
    MATCHES=$(grep -rEi "$PATTERN" . \
        --exclude-dir={.git,.history,node_modules,dist,tmp,vendor,docs/adr} \
        --exclude={aurelia.log,*.png,*.jpg,*.jpeg,*.gif,*.svg,*.pdf,scan-legacy-models.sh,purge-gemma.sh} || true)
    
    if [ -n "$MATCHES" ]; then
        echo -e "${COLOR_RED}Found matches:${COLOR_NC}"
        echo "$MATCHES"
        COUNT=$(echo "$MATCHES" | wc -l)
        TOTAL_MATCHES=$((TOTAL_MATCHES + COUNT))
    else
        echo -e "${COLOR_GREEN}Clear!${COLOR_GREEN}"
    fi
done

echo -e "\n--------------------------------"
if [ "$TOTAL_MATCHES" -eq 0 ]; then
    echo -e "${COLOR_GREEN}✅ No legacy model references found. Total Purity Achieved.${COLOR_NC}"
    exit 0
else
    echo -e "${COLOR_RED}❌ Found $TOTAL_MATCHES legacy references. Action required!${COLOR_NC}"
    exit 1
fi
