#!/bin/bash
# Aurelia Auto-Changelog Generator SOTA 2026

echo "# CHANGELOG - Ecossistema Aurélia 🏮" > CHANGELOG.md
echo "" >> CHANGELOG.md
echo "## [$(date +'%Y-%m-%d')] v$(cat services/aurelia-api/main.go | grep '"version"' | cut -d'"' -f4 || echo "0.12.1-nebula")" >> CHANGELOG.md
echo "" >> CHANGELOG.md

echo "### 🚀 Features" >> CHANGELOG.md
git log --pretty=format:"- %s (%h)" | grep -E "^- (feat|FEAT)" >> CHANGELOG.md

echo "" >> CHANGELOG.md
echo "### 🐛 Bug Fixes" >> CHANGELOG.md
git log --pretty=format:"- %s (%h)" | grep -E "^- (fix|FIX)" >> CHANGELOG.md

echo "" >> CHANGELOG.md
echo "### ⚙️ Maintenance & Chore" >> CHANGELOG.md
git log --pretty=format:"- %s (%h)" | grep -E "^- (chore|CHORE|refactor|docs|test)" >> CHANGELOG.md

echo "" >> CHANGELOG.md
echo "---" >> CHANGELOG.md
echo "Gerado automaticamente pela Aurélia em $(date)" >> CHANGELOG.md
