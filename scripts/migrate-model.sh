#!/bin/bash
# Script para migrar referências de qwen3.5:9b e qwen2.5vl:7b para gemma3:27b-it-qat
# Executar na raiz do repositório

set -e

echo "🔄 Migrando referências de modelos..."

# Substituir no visao-qwen-vl skill
sed -i 's/qwen3\.5:9b/gemma3:27b-it-qat/g' .agent/skills/visao-qwen-vl/SKILL.md
sed -i 's/qwen2\.5vl:7b/gemma3:27b-it-qat/g' .agent/skills/visao-qwen-vl/SKILL.md
sed -i 's/qwen3\.5:9b/gemma3:27b-it-qat/g' .agent/skills/visao-qwen-vl/SKILL.md

# Substituir no litellm config
sed -i 's/qwen3\.5:9b/gemma3:27b-it-qat/g' configs/litellm/config.yaml

# Substituir no rate-limit-check
sed -i 's/qwen3\.5:9b/gemma3:27b-it-qat/g' scripts/rate-limit-check.sh

# Substituir na documentação
sed -i 's/qwen3\.5:9b/gemma3:27b-it-qat/g' docs/architecture-complete.md
sed -i 's/qwen2\.5vl:7b/gemma3:27b-it-qat/g' docs/architecture-complete.md
sed -i 's/qwen3\.5:9b/gemma3:27b-it-qat/g' docs/reports/20260401-infra-audit.md
sed -i 's/qwen2\.5vl:7b/gemma3:27b-it-qat/g' docs/reports/20260401-infra-audit.md

# Substituir no .env.example
sed -i 's/qwen3\.5:9b/gemma3:27b-it-qat/g' .env.example

# NÃO mexer em:
# - internal/telegram/input_guard.go (qwen2.5:0.5b é o guard small model)
# - docs/governance/PORTEIRO_SENTINEL_2026.md (guard model)
# - .agent/skills/porteiro-ops/SKILL.md (guard model)

echo "✅ Concluído. Verifique as mudanças com:"
echo "   git diff --stat"
