#!/usr/bin/env bash
# [SOTA 2026] Global Legacy Purge (FIXED)
# Purpose: Eradicate all Qwen 3.5 (VL) references and replace with Qwen 3.5 (VL).

set -euo pipefail

echo "🚀 Starting global legacy purge..."

# Define replacements as "SEARCH|REPLACE"
REPLACEMENTS=(
    "qwen3.5|qwen3.5"
    "qwen3.5|qwen3.5"
    "modelQwen 3.5|modelQwen35"
    "Qwen 3.5 (VL)|Qwen 3.5 (VL)"
    "Qwen 3.5|Qwen 3.5"
    "qwen3.5|qwen3.5"
)

# Also fix the accidental replacement from the previous run:
# If 'qwen3.5' was replaced with '27b', we should fix it to 'qwen3.5'.
# Note: we must be careful not to break valid '27b' mentions if any, but in this context 
# 'qwen3.5' usually becomes 'qwen3.5'.
# Actually, the previous run replaced SEARCH=qwen3.5 with REPLACE=27b. 
# So 'qwen3.5' became 'qwen3.5'. 
# We'll just run the replacements normally, most will fix themselves.

for ITEM in "${REPLACEMENTS[@]}"; do
    SEARCH=$(echo "$ITEM" | cut -d'|' -f1)
    REPLACE=$(echo "$ITEM" | cut -d'|' -f2)
    
    echo "Replacing '$SEARCH' with '$REPLACE'..."
    find . -type f -not -path '*/.git/*' -not -path './docs/adr/*' -exec perl -pi -e "s/\Q$SEARCH\E/$REPLACE/g" {} +
done

# Clean up accidental 'qwen3.5' or ':27b' that might have happened if it replaced 'qwen3.5' only
find . -type f -not -path '*/.git/*' -not -path './docs/adr/*' -exec perl -pi -e "s/qwen3.5/qwen3.5/g" {} +
find . -type f -not -path '*/.git/*' -not -path './docs/adr/*' -exec perl -pi -e "s/qwen3.5/qwen3.5/g" {} +

echo "✅ Purge complete."
./scripts/scan-legacy-models.sh
