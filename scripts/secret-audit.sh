#!/usr/bin/env bash
# scripts/secret-audit.sh
# Audit script to detect plaintext credentials, API keys, and tokens
# that should NOT be in logs, source code, or config files
# Usage: bash scripts/secret-audit.sh
# Recommended: crontab -e → 0 6 * * 1 bash ~/aurelia/scripts/secret-audit.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_FILE="$HOME/.aurelia/logs/secret-audit.log"
FINDINGS=0

mkdir -p "$(dirname "$LOG_FILE")"

echo "🔍 Secret Audit — $(date)" | tee -a "$LOG_FILE"
echo "Repository: $REPO_ROOT" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Pattern definitions
PATTERNS=(
    # Passwords
    "password['\"]?\s*[:=]\s*['\"]?[a-zA-Z0-9!@#$%^&*._-]{8,}"
    "passwd['\"]?\s*[:=]\s*['\"]?[a-zA-Z0-9!@#$%^&*._-]{8,}"

    # API Keys
    "api[_-]?key['\"]?\s*[:=]\s*['\"]?[a-zA-Z0-9_-]{20,}"
    "apikey['\"]?\s*[:=]\s*['\"]?[a-zA-Z0-9_-]{20,}"

    # Bearer tokens
    "bearer\s+[A-Za-z0-9_.-]{20,}"
    "authorization['\"]?\s*[:=]\s*['\"]?Bearer\s+[A-Za-z0-9_.-]{20,}"

    # AWS keys
    "AKIA[0-9A-Z]{16}"
    "aws_secret_access_key['\"]?\s*[:=]"

    # GitHub tokens
    "ghp_[A-Za-z0-9_]{36,}"
    "github[_-]?token['\"]?\s*[:=]"

    # Telegram tokens
    "[0-9]{10}:AA[A-Za-z0-9_-]{25,}"
    "telegram[_-]?token['\"]?\s*[:=]"

    # Private keys
    "-----BEGIN (RSA|DSA|EC|PGP|OPENSSH) PRIVATE KEY"

    # Database passwords
    "postgresql://.*:.*@"
    "mysql://.*:.*@"
    "mongodb+srv://.*:.*@"
)

# Files to exclude
EXCLUDE_PATTERNS=(
    ".git/"
    ".gitignore"
    "node_modules/"
    "vendor/"
    ".venv/"
    "__pycache__/"
    "*.min.js"
    "*.min.css"
    "secrets.env"  # This is expected to exist
    "secret-audit.sh"  # This script itself
    ".bak"
    ".backup"
)

# Build exclude regex
EXCLUDE_REGEX=""
for pattern in "${EXCLUDE_PATTERNS[@]}"; do
    if [ -z "$EXCLUDE_REGEX" ]; then
        EXCLUDE_REGEX="$pattern"
    else
        EXCLUDE_REGEX="$EXCLUDE_REGEX|$pattern"
    fi
done

echo "Scanning for patterns:" | tee -a "$LOG_FILE"
for pattern in "${PATTERNS[@]:0:3}"; do
    echo "  • $(echo "$pattern" | cut -c 1-50)..." | tee -a "$LOG_FILE"
done
echo "  (and $(( ${#PATTERNS[@]} - 3 )) more patterns)" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"

# Scan git history for secrets
echo "Scanning git history for secrets..." | tee -a "$LOG_FILE"
if git -C "$REPO_ROOT" rev-parse --git-dir > /dev/null 2>&1; then
    # Look for common patterns in commit history
    if git -C "$REPO_ROOT" log -p -S "password" -S "token" -S "secret" -S "key" --all 2>/dev/null | grep -iE "(password|token|secret|api.?key|bearer)" | head -5 | grep -vE "$EXCLUDE_REGEX" > /dev/null; then
        echo -e "${RED}⚠️  Found potential secrets in git history${NC}" | tee -a "$LOG_FILE"
        ((FINDINGS++))
    else
        echo -e "${GREEN}✅ No obvious secrets in git history${NC}" | tee -a "$LOG_FILE"
    fi
fi

echo "" | tee -a "$LOG_FILE"

# Scan logs directory
echo "Scanning logs for plaintext credentials..." | tee -a "$LOG_FILE"
if [ -d "$HOME/.aurelia/logs" ]; then
    for logfile in "$HOME/.aurelia/logs"/*.log; do
        if [ -f "$logfile" ]; then
            for pattern in "${PATTERNS[@]}"; do
                if grep -iE "$pattern" "$logfile" 2>/dev/null | grep -v "pattern\|regex" | head -1 > /dev/null; then
                    echo -e "${RED}⚠️  Found in $logfile: $(basename "$logfile")${NC}" | tee -a "$LOG_FILE"
                    ((FINDINGS++))
                    break
                fi
            done
        fi
    done
fi

if [ $FINDINGS -eq 0 ]; then
    echo -e "${GREEN}✅ No secrets found in logs${NC}" | tee -a "$LOG_FILE"
fi

echo "" | tee -a "$LOG_FILE"

# Scan source code (excluding common safe patterns)
echo "Scanning source code for hardcoded secrets..." | tee -a "$LOG_FILE"
FOUND_CODE_SECRETS=0
for pattern in "password.*=" "token.*=" "secret.*=" "api.?key.*="; do
    # Search in source files but exclude safe patterns
    if find "$REPO_ROOT" -type f \( -name "*.go" -o -name "*.py" -o -name "*.js" -o -name "*.ts" \) \
        ! -path "$REPO_ROOT/.git/*" \
        ! -path "$REPO_ROOT/node_modules/*" \
        ! -path "$REPO_ROOT/vendor/*" \
        -exec grep -l "$pattern" {} \; 2>/dev/null | head -3 > /dev/null; then

        # Check if it's actually a hardcoded value (not env var placeholder)
        FILES=$(find "$REPO_ROOT" -type f \( -name "*.go" -o -name "*.py" -o -name "*.js" -o -name "*.ts" \) \
            ! -path "$REPO_ROOT/.git/*" \
            ! -path "$REPO_ROOT/node_modules/*" \
            ! -path "$REPO_ROOT/vendor/*" \
            -exec grep -l "$pattern" {} \; 2>/dev/null | head -3)

        for file in $FILES; do
            # Look for actual values, not ${PLACEHOLDERS}
            if grep "$pattern" "$file" | grep -v '\${' | grep -vE '(test|example|mock|placeholder)' > /dev/null 2>&1; then
                echo -e "${YELLOW}⚠️  Check $file (may contain hardcoded secrets)${NC}" | tee -a "$LOG_FILE"
                ((FOUND_CODE_SECRETS++))
            fi
        done
    fi
done

if [ $FOUND_CODE_SECRETS -eq 0 ]; then
    echo -e "${GREEN}✅ No hardcoded secrets in source code${NC}" | tee -a "$LOG_FILE"
else
    ((FINDINGS += FOUND_CODE_SECRETS))
fi

echo "" | tee -a "$LOG_FILE"

# Check for .env files with credentials
echo "Scanning for .env files..." | tee -a "$LOG_FILE"
ENV_FILES=$(find "$REPO_ROOT" -name ".env*" -o -name "*.env" 2>/dev/null | grep -v "secrets.env\|\.env\.example")
if [ -n "$ENV_FILES" ]; then
    echo -e "${YELLOW}⚠️  Found .env files (verify they are git-ignored):${NC}" | tee -a "$LOG_FILE"
    echo "$ENV_FILES" | tee -a "$LOG_FILE"
    ((FINDINGS++))
else
    echo -e "${GREEN}✅ No unexpected .env files${NC}" | tee -a "$LOG_FILE"
fi

echo "" | tee -a "$LOG_FILE"

# Final verdict
if [ $FINDINGS -eq 0 ]; then
    echo -e "${GREEN}✅ AUDIT PASSED — No secrets found${NC}" | tee -a "$LOG_FILE"
    exit 0
else
    echo -e "${RED}⚠️  AUDIT FAILED — $FINDINGS potential issue(s) found${NC}" | tee -a "$LOG_FILE"
    echo "Review findings above and remediate if needed." | tee -a "$LOG_FILE"
    exit 1
fi
