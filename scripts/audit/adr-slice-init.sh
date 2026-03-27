#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
ADR_DIR="$ROOT_DIR/docs/adr"
TASKMASTER_DIR="$ADR_DIR/taskmaster"
TEMPLATE_MD="$ADR_DIR/TEMPLATE-NONSTOP-SLICE.md"
TEMPLATE_JSON="$TASKMASTER_DIR/TEMPLATE-NONSTOP-SLICE.json"

usage() {
  cat <<'EOF'
usage: ./scripts/adr-slice-init.sh <slug> [--title "Title"] [--date YYYYMMDD] [--dry-run]
EOF
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

SLUG="$1"
shift
TITLE=""
DATE_OVERRIDE=""
DRY_RUN=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --title)
      shift
      TITLE="${1:-}"
      ;;
    --date)
      shift
      DATE_OVERRIDE="${1:-}"
      ;;
    --dry-run)
      DRY_RUN=1
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
  shift || true
done

DATE_VALUE="${DATE_OVERRIDE:-$(date +%Y%m%d)}"
ADR_ID="ADR-${DATE_VALUE}-${SLUG}"
ADR_PATH="$ADR_DIR/${ADR_ID}.md"
JSON_PATH="$TASKMASTER_DIR/${ADR_ID}.json"

render_md() {
  sed \
    -e "s/ADR-YYYYMMDD-slug/${ADR_ID}/g" \
    -e "s/- slug:/- slug: ${SLUG}/" \
    -e "s|- json de continuidade:|- json de continuidade: ${JSON_PATH#"$ROOT_DIR/"}|" \
    "$TEMPLATE_MD"
}

render_json() {
  sed \
    -e "s/ADR-YYYYMMDD-slug/${ADR_ID}/g" \
    -e "s/Slice title/${TITLE:-$SLUG}/g" \
    "$TEMPLATE_JSON"
}

if [[ "$DRY_RUN" -eq 1 ]]; then
  echo "ADR_PATH=$ADR_PATH"
  echo "JSON_PATH=$JSON_PATH"
  exit 0
fi

mkdir -p "$ADR_DIR" "$TASKMASTER_DIR"

if [[ -e "$ADR_PATH" || -e "$JSON_PATH" ]]; then
  echo "ADR or JSON already exists for $ADR_ID" >&2
  exit 1
fi

render_md > "$ADR_PATH"
render_json > "$JSON_PATH"

echo "created:"
echo "  $ADR_PATH"
echo "  $JSON_PATH"
