#!/usr/bin/env bash
# ADR Semparar Validation — Garante estabilidade do workflow de slices nonstop
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
ADR_DIR="$ROOT_DIR/docs/adr"
TASKMASTER_DIR="$ADR_DIR/taskmaster"

VALID_MD_STATUSES=("Proposto" "Em execução" "Aceito" "Bloqueado" "Cancelado")
VALID_JSON_STATUSES=("proposed" "in_progress" "accepted" "blocked" "cancelled")

echo "🔍 Validando workflow /adr-semparar..."

# Contar ADRs
ADR_COUNT=$(find "$ADR_DIR" -maxdepth 1 -name "ADR-202603*-*.md" -type f | wc -l)
JSON_COUNT=$(find "$TASKMASTER_DIR" -maxdepth 1 -name "ADR-202603*-*.json" -type f | wc -l)

echo "  ADRs: $ADR_COUNT | JSONs: $JSON_COUNT"

if [[ "$ADR_COUNT" -ne "$JSON_COUNT" ]]; then
  echo "❌ Mismatch: $ADR_COUNT MDs vs $JSON_COUNT JSONs"
  exit 1
fi

# Validar cada par MD + JSON
ERRORS=0
WARNINGS=0

for adr_md in $(find "$ADR_DIR" -maxdepth 1 -name "ADR-202603*-*.md" -type f | sort); do
  BASENAME=$(basename "$adr_md" .md)
  JSON_FILE="$TASKMASTER_DIR/${BASENAME}.json"

  if [[ ! -f "$JSON_FILE" ]]; then
    echo "❌ Falta JSON para: $BASENAME"
    ((ERRORS++))
    continue
  fi

  # Checar status em MD
  if ! grep -q "^- \(Proposto\|Em execução\|Aceito\|Bloqueado\|Cancelado\)$" "$adr_md"; then
    echo "⚠️  Status inválido em $BASENAME.md"
    ((WARNINGS++))
  fi

  # Checar JSON válido
  if ! jq empty "$JSON_FILE" 2>/dev/null; then
    echo "❌ JSON inválido: $BASENAME.json"
    ((ERRORS++))
    continue
  fi

  # Checar campos obrigatórios em JSON
  REQUIRED_FIELDS=("adr_id" "title" "status" "progress" "goal" "next_actions" "handoff")
  for field in "${REQUIRED_FIELDS[@]}"; do
    if ! jq -e ".$field" "$JSON_FILE" >/dev/null 2>&1; then
      echo "❌ Campo obrigatório faltando em $BASENAME.json: $field"
      ((ERRORS++))
    fi
  done

  # Checar handoff.resume_prompt preenchido
  RESUME_PROMPT=$(jq -r '.handoff.resume_prompt // empty' "$JSON_FILE")
  if [[ -z "$RESUME_PROMPT" ]]; then
    echo "⚠️  Resume prompt vazio em $BASENAME.json"
    ((WARNINGS++))
  fi

  # Validar frontmatter do MD
  if ! head -5 "$adr_md" | grep -q "^---$"; then
    echo "⚠️  Frontmatter faltando em $BASENAME.md"
    ((WARNINGS++))
  fi
done

echo ""
echo "📊 Resumo:"
echo "  ✓ $ADR_COUNT ADRs com JSON correspondente"
echo "  ❌ Erros críticos: $ERRORS"
echo "  ⚠️  Avisos: $WARNINGS"

if [[ $ERRORS -gt 0 ]]; then
  echo ""
  echo "❌ Validação FALHOU"
  exit 1
else
  echo ""
  echo "✅ Validação OK — Workflow /adr-semparar está ESTÁVEL"
  exit 0
fi
